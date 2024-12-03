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

	// bitflyer
	var _bitflyer = new(bitflyer.Bitflyer)
	// bitflyer xrpjpy
	var klineBitflyerXRPJPY *[]exchange.KLine = &[]exchange.KLine{}
	// bitflyer fx_btc_jpy
	// var klineBitflyerFXBTCJPY *[]exchange.KLine = &[]exchange.KLine{}

	trades := []struct {
		exchange exchange.Exchange
		kine     *[]exchange.KLine
		symbol   exchange.Symbol
	}{
		{_bitflyer, klineBitflyerXRPJPY, exchange.XRPJPY},
		// {_bitflyer, klineBitflyerFXBTCJPY, exchange.FX_BTCJPY},
	}

	// init log
	logger.InitLogger(nil, "", 10, 5, 30, true)
	for _, t := range trades {
		logger.InitLogger(t.exchange, t.symbol, 10, 5, 30, true)
		logger.Print(t.exchange, t.symbol, "Time, PriceNow, Open, Close, High, Low, SMA, EMA, BBands+3, BBands+2, BBands-2, BBands-3, K, D, SMASlope, RSI, BUY, SELL")
	}

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
