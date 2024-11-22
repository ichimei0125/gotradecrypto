package main

import (
	"log"
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/common"
	"github.com/ichimei0125/gotradecrypto/internal/exchange"
	"github.com/ichimei0125/gotradecrypto/internal/exchange/bitflyer"
	"github.com/ichimei0125/gotradecrypto/internal/indicator"
	"github.com/ichimei0125/gotradecrypto/internal/trade"
)

func main() {
	var kline *[]exchange.KLine = &[]exchange.KLine{}

	log.Println("Time, CloseTime, kline, SMA, EMA, BBands+3, BBands+2, BBands-2, BBands-3, K, D, BUY, SELL")
	for {
		bitflyer := &bitflyer.Bitflyer{}
		symbol := exchange.XRPJPY

		exchange.GetKLine(bitflyer, symbol, kline)
		indicator.GetIndicators(kline)

		go trade.Trade(bitflyer, symbol, kline)

		// sleep
		time.Sleep(time.Duration(common.REFRESH_INTERVAL) * time.Minute)
	}
}
