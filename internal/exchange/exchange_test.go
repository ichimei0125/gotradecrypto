package exchange_test

import (
	"fmt"
	"testing"

	"github.com/ichimei0125/gotradecrypto/internal/exchange"
	"github.com/ichimei0125/gotradecrypto/internal/exchange/bitflyer"
)

func TestIsDryRun(t *testing.T) {
	symbol := exchange.FX_BTCJPY
	e := &bitflyer.Bitflyer{}
	is_dry_run := symbol.IsDryRun(e.Name())
	if is_dry_run != true {
		t.Error(is_dry_run)
	}

	// symbol.IsDryRun("test")
	symbol = exchange.ETHJPY
	symbol.IsDryRun(e.Name())

}

func TestGetSecrets(t *testing.T) {
	e := &bitflyer.Bitflyer{}
	key, secret := exchange.GetSecret(e.Name())
	fmt.Println(key, secret)

	exchange.GetSecret("test")
}
