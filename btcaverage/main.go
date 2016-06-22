package btcaverage

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
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
