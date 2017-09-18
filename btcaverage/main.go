package btcaverage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	baseURL = "https://apiv2.bitcoinaverage.com"
)

type GlobalTracker struct {
	tickers *map[string]Ticker
	symbols *Symbols
	mu      sync.RWMutex
	period  time.Duration
}

func NewGlobalTracker() (*GlobalTracker, error) {
	log.Printf("Initializing GlobalTracker...")
	sym, err := GetSymbols()
	if err != nil {
		return nil, err
	}

	// 10-minute period because free API level allows 5k calls/month, which is ~1 every 8.64 minutes
	gt := &GlobalTracker{tickers: &map[string]Ticker{}, symbols: sym, period: 10 * time.Minute}
	go gt.poll()

	return gt, nil
}

func (gt *GlobalTracker) GetAvg(symbol string) (Ticker, error) {
	if _, ok := (*gt.tickers)[symbol]; ok {
		gt.mu.RLock()
		t := (*gt.tickers)[symbol]
		gt.mu.RUnlock()
		return t, nil
	} else {
		log.Printf("Not tracking %s yet, fetching...", symbol)
		t, err := GetCurrentGlobalTicker(symbol)
		if err != nil {
			return Ticker{}, err
		}

		gt.mu.Lock()
		(*gt.tickers)[symbol] = *t
		gt.mu.Unlock()
		return *t, nil
	}
}

func (gt *GlobalTracker) poll() {
	for {
		time.Sleep(gt.period)

		for sym, _ := range *gt.tickers {
			log.Printf("Fetching avg for %s...", sym)
			ticker, err := GetCurrentGlobalTicker(sym)
			if err != nil {
				log.Printf("Error fetching global avg: %s", err.Error())
				continue
			}
			gt.mu.Lock()
			(*gt.tickers)[sym] = *ticker
			gt.mu.Unlock()
			time.Sleep(gt.period)
		}
	}
}

func GetCurrentGlobalTicker(symbol string) (*Ticker, error) {
	if symbol != "BTCUSD" {
		panic("Symbols other than BTCUSD are not supported!")
	}

	res, err := http.Get(baseURL + "/indices/global/ticker/" + symbol)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Got %d response: \"%s\"", res.StatusCode, body)
	}

	var ticker Ticker
	if err = json.Unmarshal(body, &ticker); err != nil {
		return nil, err
	}

	return &ticker, nil
}

func GetSymbols() (*Symbols, error) {
	res, err := http.Get(baseURL + "/constants/symbols")
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	var sym Symbols
	if err = json.Unmarshal(body, &sym); err != nil {
		return nil, err
	}

	return &sym, nil
}
