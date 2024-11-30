package main

import (
	"sync"
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

	// bitflyer
	var _bitflyer = new(bitflyer.Bitflyer)
	// bitflyer xrpjpy
	var klineBitflyerXRPJPY *[]exchange.KLine = &[]exchange.KLine{}

	trades := []struct {
		exchange exchange.Exchange
		kine     *[]exchange.KLine
		symbol   exchange.Symbol
	}{
		{_bitflyer, klineBitflyerXRPJPY, exchange.XRPJPY},
	}

	logger.Info("Time, PriceNow, Open, Close, High, Low, SMA, EMA, BBands+3, BBands+2, BBands-2, BBands-3, K, D, SMASlope, RSI, BUY, SELL")
	wg := new(sync.WaitGroup)
	for {
		wg.Add(len(trades))

		for _, t := range trades {
			localT := t // 闭包变量捕获问题
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				localT.exchange.FetchKLine(localT.symbol, localT.kine)
				indicator.GetIndicators(localT.kine)

				go trade.Trade(localT.exchange, localT.symbol, localT.kine)
			}(wg)
		}

		wg.Wait()
		// sleep
		time.Sleep(time.Duration(common.REFRESH_INTERVAL) * time.Minute)
	}
}
