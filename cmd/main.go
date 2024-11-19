package main

import (
	"log"
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/common"
	"github.com/ichimei0125/gotradecrypto/internal/exchange"
	"github.com/ichimei0125/gotradecrypto/internal/exchange/bitflyer"
	"github.com/ichimei0125/gotradecrypto/internal/indicator"
)

func main() {
	for {
		bitflyer := &bitflyer.Bitflyer{}

		symbol := exchange.BTCJPY

		kline := exchange.GetKLine(bitflyer, symbol)
		indicator.GetIndicators(&kline)

		if len(kline) < common.KLINE_LENGTH {
			log.Printf("not enough data")
		} else {
			log.Printf("CloseTime %s, kline %f, SMA %f, EMA %f, BBands+3 %f, BBands+2 %f, BBands-2 %f, BBands-3 %f", kline[1].CloseTime, kline[1].Open, kline[1].SMA, kline[1].EMA, kline[1].BBands_Plus_3K, kline[1].BBands_Plus_2K, kline[1].BBands_Minus_2K, kline[1].BBands_Minus_3K)
		}

		// sleep
		time.Sleep(time.Duration(common.REFRESH_INTERVAL) * time.Minute)

	}
}
