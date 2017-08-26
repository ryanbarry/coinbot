package btcaverage

import (
	"fmt"
	"testing"
	"time"
)

func TestGetGlobalTicker(t *testing.T) {
	tckr, err := GetCurrentGlobalTicker("BTCUSD")

	if err != nil {
		t.Errorf("got err! (%v)", err)
	}
	fmt.Printf("Epoch parsed: %s\nDisplay parsed: %s\n", tckr.Timestamp.Format(time.RFC1123Z), tckr.DisplayTimestamp.Format(time.RFC1123Z))
}
