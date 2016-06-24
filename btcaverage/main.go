package btcaverage

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type btatime struct {
	time.Time
}
type GlobalAvg struct {
	Rolling24h     float64 `json:"24h_avg"`
	Ask            float64
	Bid            float64
	Last           float64
	Timestamp      btatime
	Volume_btc     float64
	Volume_percent float64
}

type GlobalTracker struct {
	avg      *GlobalAvg
	mu       sync.RWMutex
	currency string
	period   time.Duration
}

func NewGlobalTracker(currency string) (*GlobalTracker, error) {
	initAvg, err := GetCurrentGlobalAvg(currency)
	if err != nil {
		return nil, err
	}

	gt := &GlobalTracker{avg: initAvg, currency: currency, period: time.Minute}
	go gt.Poll()

	return gt, nil
}

func (gt *GlobalTracker) GetAvg() GlobalAvg {
	gt.mu.RLock()
	res := *gt.avg
	gt.mu.RUnlock()
	return res
}

func (gt *GlobalTracker) Poll() {
	for {
		time.Sleep(gt.period)
		log.Printf("Fetching new global average...")
		avg, err := GetCurrentGlobalAvg(gt.currency)
		if err != nil {
			log.Printf("Error fetching global avg: ", err.Error())
			continue
		}
		gt.mu.Lock()
		gt.avg = avg
		gt.mu.Unlock()
	}
}

func GetCurrentGlobalAvg(currency string) (*GlobalAvg, error) {
	if !isValidCurrency(currency) {
		return nil, errors.New("Invalid currency code! (" + currency + ")")
	}

	res, err := http.Get("https://api.bitcoinaverage.com/ticker/global/" + currency)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	var avgObj GlobalAvg
	if err = json.Unmarshal(body, &avgObj); err != nil {
		return nil, err
	}

	return &avgObj, nil
}

func isValidCurrency(currency string) bool {
	valid := []string{"USD", "EUR", "GBP"}
	for _, cc := range valid {
		if cc == currency {
			return true
		}
	}

	return false
}

func (btat *btatime) UnmarshalJSON(b []byte) error {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}

	t, err := time.Parse(time.RFC1123Z, string(b))
	if err != nil {
		return err
	}

	*btat = btatime{t}
	return nil
}
