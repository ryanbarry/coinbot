package btcaverage

import (
	"fmt"
	"testing"
	"time"
)

func TestGetSymbols(t *testing.T) {
	s, err := GetSymbols()
	if err != nil {
		t.Errorf("Got err! (%v)", err)
	}
	fmt.Printf("Symbols obj has %d Crypto, %d Global, and %d Local symbols\n", len(s.Crypto.Symbols), len(s.Global.Symbols), len(s.Local.Symbols))
}

func TestGetGlobalTicker(t *testing.T) {
	tracker, err := NewGlobalTracker()
	tckr, err := tracker.GetCurrentGlobalTicker("BTCUSD")

	if err != nil {
		t.Errorf("got err! (%v)", err)
	}
	fmt.Printf("Epoch parsed: %s\nDisplay parsed: %s\n", tckr.Timestamp.Format(time.RFC1123Z), tckr.DisplayTimestamp.Format(time.RFC1123Z))
}
