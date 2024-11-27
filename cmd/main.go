package main

import (
	"fmt"
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/common"
	"github.com/ichimei0125/gotradecrypto/internal/exchange"
	"github.com/ichimei0125/gotradecrypto/internal/exchange/bitflyer"
	"github.com/ichimei0125/gotradecrypto/internal/indicator"
	"github.com/ichimei0125/gotradecrypto/internal/logger"
	"github.com/ichimei0125/gotradecrypto/internal/trade"
)

func main() {
	logger.InitLogger("log/app.log", 10, 5, 30, true)

	// bitflyer xprjpy
	var klineBitflyerXRPJPY *[]exchange.KLine = &[]exchange.KLine{}
	bitflyerXRPJPY := &bitflyer.Bitflyer{}
	bitflyerxrpjpy := exchange.XRPJPY

	trades := []struct {
		exchange exchange.Exchange
		kine     *[]exchange.KLine
		symbol   exchange.Symbol
	}{
		{bitflyerXRPJPY, klineBitflyerXRPJPY, bitflyerxrpjpy},
	}

	fmt.Println("Time, CloseTime, kline, SMA, EMA, BBands+3, BBands+2, BBands-2, BBands-3, K, D, SMASlope, BUY, SELL")
	for {

		for _, t := range trades {
			t.exchange.FetchKLine(t.symbol, t.kine)
			indicator.GetIndicators(t.kine)

			go trade.Trade(t.exchange, t.symbol, t.kine)
		}

		// sleep
		time.Sleep(time.Duration(common.REFRESH_INTERVAL) * time.Minute)
	}
}
