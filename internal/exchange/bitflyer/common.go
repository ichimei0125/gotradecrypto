package bitflyer

import (
	"bytes"
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/config"
	"github.com/ichimei0125/gotradecrypto/internal/exchange"
)

const baseURL = "https://api.bitflyer.com"

var _symbols []string = []string{}

type Bitflyer struct{}

func (b *Bitflyer) GetInfo() exchange.ExchangeInfo {
	if len(_symbols) <= 0 {
		GetSymbols()
	}

	return exchange.ExchangeInfo{
		Name:       "bitflyer",
		IsDryRun:   true,
		IsReattime: false,
		Waittime:   time.Duration(1 * time.Minute),
		Symbols:    _symbols,
	}
}

func GetSymbols() []string {
	_symbols = config.GetConfig().Symbols[new(Bitflyer).GetInfo().Name]
	return _symbols
}

type Execution struct {
	ID                         int64      `json:"id"`
	Side                       string     `json:"side"`
	Price                      float64    `json:"price"`
	Size                       float64    `json:"size"`
	ExecDate                   CustomTime `json:"exec_date"`
	BuyChildOrderAcceptanceId  string     `json:"buy_child_order_acceptance_id"`
	SellChildOrderAcceptanceId string     `json:"sell_child_order_acceptance_id"`
}

type CustomTime struct {
	time.Time
}

var TimeLayout = "2006-01-02T15:04:05.000"

var timeFormats = []string{
	"2006-01-02T15:04:05.000",
	"2006-01-02T15:04:05.00",
	"2006-01-02T15:04:05.0",
	"2006-01-02T15:04:05",
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := bytes.Trim(b, "\"")
	if len(s) == 0 {
		return nil
	}
	var parseErr error
	for _, layout := range timeFormats {
		t, err := time.Parse(layout, string(s))
		if err == nil {
			ct.Time = t
			return nil
		}
		parseErr = err
	}
	return parseErr
}
