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
	log.Println("Time, CloseTime, kline, SMA, SMA_NEW, EMA, BBands+3, BBands+2, BBands-2, BBands-3, K, D")
	for {
		bitflyer := &bitflyer.Bitflyer{}

		symbol := exchange.BTCJPY

		kline := exchange.GetKLine(bitflyer, symbol)
		indicator.GetIndicators(&kline)

		if len(kline) < common.KLINE_LENGTH {
			log.Printf("not enough data")
		} else {
			log.Printf(" ,%s, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f", kline[1].CloseTime, kline[1].Open, kline[1].SMA, kline[1].SMA_new, kline[1].EMA, kline[1].BBands_Plus_3K, kline[1].BBands_Plus_2K, kline[1].BBands_Minus_2K, kline[1].BBands_Minus_3K, kline[1].SlowK, kline[1].SlowD)
		}

		// sleep
		time.Sleep(time.Duration(common.REFRESH_INTERVAL) * time.Minute)

	}
}
