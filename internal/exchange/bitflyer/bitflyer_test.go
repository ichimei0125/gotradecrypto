package bitflyer_test

import (
	"testing"
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/db"
	"github.com/ichimei0125/gotradecrypto/internal/exchange/bitflyer"
)

func TestFetchTrades(t *testing.T) {
	db.InitDB()
	defer db.CloseDB()
	since := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	bitflyer.FetchTrades(since, "XLM_JPY")
}
