package main

import (
	"log"
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/common"
	"github.com/ichimei0125/gotradecrypto/internal/exchange"
	"github.com/ichimei0125/gotradecrypto/internal/exchange/bitflyer"
)

func main() {
	for {
		bitflyer := &bitflyer.Bitflyer{}

		symbol := exchange.BTCJPY

		kline := exchange.GetKLineAndIndicators(bitflyer, symbol)

		if len(kline) < common.KLINE_LENGTH {
			log.Printf("not enough data")
		} else {
			log.Printf("kline %f, len %d, cap %d", kline[0].Open, len(kline), cap(kline))
		}

		// sleep
		time.Sleep(time.Duration(common.REFRESH_INTERVAL) * time.Minute)

	}
}
