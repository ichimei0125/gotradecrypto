package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/common"
	"github.com/ichimei0125/gotradecrypto/internal/exchange"
	"github.com/ichimei0125/gotradecrypto/internal/exchange/bitflyer"
	"github.com/ichimei0125/gotradecrypto/internal/indicator"
	"github.com/ichimei0125/gotradecrypto/internal/trade"
)

func main() {
	log.Println("GoTradeCypto started")

	// 捕获信号
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		log.Println("Time, CloseTime, kline, SMA, EMA, BBands+3, BBands+2, BBands-2, BBands-3, K, D, SMASlope, BUY, SELL")
		for {
			var kline *[]exchange.KLine = &[]exchange.KLine{}
			bitflyer := &bitflyer.Bitflyer{}
			symbol := exchange.XRPJPY

			exchange.GetKLine(bitflyer, symbol, kline)
			indicator.GetIndicators(kline)

			go trade.Trade(bitflyer, symbol, kline)

			// sleep
			time.Sleep(time.Duration(common.REFRESH_INTERVAL) * time.Minute)
		}

	}()

	<-stop
	log.Println("GoTradeCypto stopped")
}
