package btcaverage

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"sync"
	"time"
)

type btatime struct {
	time.Time
}
type Avg struct {
	Rolling24h     float64 `json:"24h_avg"`
	Ask            float64
	Bid            float64
	Last           float64
	Timestamp      btatime
	Volume_btc     float64
	Volume_percent float64
}

type GlobalTracker struct {
	avg    *map[string]Avg
	mu     sync.RWMutex
	period time.Duration
}

func NewGlobalTracker() (*GlobalTracker, error) {
	initAvg, err := GetCurrentGlobalAvg()
	if err != nil {
		return nil, err
	}

	gt := &GlobalTracker{avg: initAvg, period: time.Minute}
	go gt.Poll()

	return gt, nil
}

func (gt *GlobalTracker) GetAvg(currency string) Avg {
	gt.mu.RLock()
	res := (*gt.avg)[currency]
	gt.mu.RUnlock()
	return res
}

func (gt *GlobalTracker) Poll() {
	for {
		time.Sleep(gt.period)
		log.Printf("Fetching new global average...")
		avg, err := GetCurrentGlobalAvg()
		if err != nil {
			log.Printf("Error fetching global avg: ", err.Error())
			continue
		}
		gt.mu.Lock()
		gt.avg = avg
		gt.mu.Unlock()
	}
}

func GetCurrentGlobalAvg() (*map[string]Avg, error) {
	res, err := http.Get("https://api.bitcoinaverage.com/ticker/global/all")
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	var avgs map[string]Avg
	if err = json.Unmarshal(body, &avgs); err != nil {
		//FIXME: this could probably be improved upon! the JSON returned by the
		// "all" call is a uniform object of currency codes to avg objects
		// EXCEPT for one additional key "timestamp" whose value is a timestamp
		// string so this expects that particular parse error and ignores it...
		if e, ok := err.(*json.UnmarshalTypeError); ok {
			if e.Type == reflect.TypeOf(Avg{}) && e.Value == "string" {
				return &avgs, nil
			}
		}
		return nil, err
	}

	return &avgs, nil
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
