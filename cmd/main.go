package main

import (
	"github.com/ichimei0125/gotradecrypto/internal/exchange"
	"github.com/ichimei0125/gotradecrypto/internal/exchange/bitflyer"
)

func main() {
	bitflyer := &bitflyer.Bitflyer{}

	symbol := exchange.BTCJPY

	exchange.FetchPriceFromExchange(bitflyer, symbol)
}
