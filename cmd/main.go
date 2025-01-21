package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/common"
	"github.com/ichimei0125/gotradecrypto/internal/db"
	"github.com/ichimei0125/gotradecrypto/internal/exchange"
	"github.com/ichimei0125/gotradecrypto/internal/exchange/bitflyer"
	"github.com/ichimei0125/gotradecrypto/internal/logger"
)

// catch app error
func handlePanic() {
	if r := recover(); r != nil {
		db.InsertErr(fmt.Sprintf("panic: %v", r))
		os.Exit(1)
	}
}

func main() {
	// defer handlePanic()

	// bitflyer
	var _bitflyer = new(bitflyer.Bitflyer)
	// bitflyer xrpjpy
	// var klineBitflyerXRPJPY *[]exchange.CandleStick = &[]exchange.CandleStick{}
	// bitflyer fx_btc_jpy
	var klineBitflyerFXBTCJPY *[]exchange.CandleStick = &[]exchange.CandleStick{}

	trades := []struct {
		exchange exchange.Exchange
		kine     *[]exchange.CandleStick
		symbol   string
	}{
		{_bitflyer, klineBitflyerFXBTCJPY, exchange.XRPJPY},
		{_bitflyer, klineBitflyerFXBTCJPY, exchange.FX_BTCJPY},
		{_bitflyer, klineBitflyerFXBTCJPY, exchange.BTCJPY},
		{_bitflyer, klineBitflyerFXBTCJPY, exchange.ETHJPY},
		{_bitflyer, klineBitflyerFXBTCJPY, exchange.XLMJPY},
		{_bitflyer, klineBitflyerFXBTCJPY, exchange.MONAJPY},
		{_bitflyer, klineBitflyerFXBTCJPY, exchange.ETHBTC},
		{_bitflyer, klineBitflyerFXBTCJPY, exchange.BCHBTC},
	}

	// init log
	logger.InitLogger(nil, "", 10, 5, 30, true)
	for _, t := range trades {
		logger.InitLogger(t.exchange, t.symbol, 10, 5, 30, true)
		logger.Print(t.exchange, t.symbol, "Time, PriceNow, Open, Close, High, Low, SMA, EMA, BBands+3, BBands+2, BBands-2, BBands-3, K, D, SMASlope, RSI, BUY, SELL")
	}

	// init db
	db.InitDB()
	defer db.CloseDB()

	// TOOD: recovery from AppErr

	wg := new(sync.WaitGroup)
	for {
		wg.Add(len(trades))

		for _, t := range trades {
			localT := t // 闭包变量捕获问题
			go func(wg *sync.WaitGroup) {
				// defer handlePanic()
				defer wg.Done()

				since := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
				test := bitflyer.FetchTrades(since, string(localT.symbol))
				fmt.Println(test[0].ExecutionTime)
				// localT.exchange.FetchCandleSticks(localT.symbol, localT.kine)
				// indicator.GetIndicators(localT.kine)

				// go func() {
				// 	// defer handlePanic()
				// 	trade.Trade(localT.exchange, localT.symbol, localT.kine)
				// }()
			}(wg)
		}

		wg.Wait()
		// sleep
		time.Sleep(time.Duration(common.REFRESH_INTERVAL) * time.Minute)
	}
}
