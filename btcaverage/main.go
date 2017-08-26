package btcaverage

import (
	"encoding/json"
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
	mu      sync.RWMutex
	period  time.Duration
}

func NewGlobalTracker() (*GlobalTracker, error) {
	btcusd, err := GetCurrentGlobalTicker("BTCUSD")
	if err != nil {
		return nil, err
	}

	gt := &GlobalTracker{tickers: &map[string]Ticker{"BTCUSD": *btcusd}, period: 10 * time.Minute}
	go gt.Poll()

	return gt, nil
}

func (gt *GlobalTracker) GetAvg(symbol string) Ticker {
	gt.mu.RLock()
	res := (*gt.tickers)[symbol]
	gt.mu.RUnlock()
	return res
}

func (gt *GlobalTracker) Poll() {
	for {
		time.Sleep(gt.period)
		log.Printf("Fetching new global average...")
		ticker, err := GetCurrentGlobalTicker("BTCUSD")
		if err != nil {
			log.Printf("Error fetching global avg: ", err.Error())
			continue
		}
		gt.mu.Lock()
		(*gt.tickers)["BTCUSD"] = *ticker
		gt.mu.Unlock()
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

	var ticker Ticker
	if err = json.Unmarshal(body, &ticker); err != nil {
		return nil, err
	}

	return &ticker, nil
}
