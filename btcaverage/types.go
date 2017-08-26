package btcaverage

import (
	"strconv"
	"time"
)

type btatime struct {
	time.Time
}

func (btat *btatime) UnmarshalJSON(b []byte) error {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}

	t, err := time.Parse("2006-01-02 15:04:05", string(b))
	if err != nil {
		return err
	}

	*btat = btatime{t}
	return nil
}

type epochsec struct {
	time.Time
}

func (eps *epochsec) UnmarshalJSON(b []byte) error {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}

	i, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}

	*eps = epochsec{time.Unix(i, 0)}
	return nil
}

type Hist struct {
	Hour    float64
	Day     float64
	Week    float64
	Month   float64
	Quarter float64 `json:"month_3"`
	Half    float64 `json:"month_6"`
	Year    float64
}

type Ticker struct {
	Ask      float64
	Bid      float64
	Last     float64
	High     float64
	Low      float64
	Open     Hist
	Averages struct {
		Day   float64
		Week  float64
		Month float64
	}
	Volume  float64
	Changes struct {
		Price   Hist
		Percent Hist
	}
	VolumePercent    float64 `json:"volume_percent"`
	Timestamp        epochsec
	DisplayTimestamp btatime `json:"display_timestamp"`
}
