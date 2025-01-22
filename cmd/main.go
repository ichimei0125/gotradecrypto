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

// exchange to trading
var exchanges = []exchange.Exchange{new(bitflyer.Bitflyer)}

// catch app error
func handlePanic() {
	if r := recover(); r != nil {
		db.InsertErr(fmt.Sprintf("panic: %v", r))
		os.Exit(1)
	}
}

type trading struct {
	exchange exchange.Exchange
	symbol   string
}

func main() {
	// defer handlePanic()

	tradings := []trading{}
	for _, _exchagne := range exchanges {
		for _, _symbol := range _exchagne.GetInfo().Symbols {
			_trading := trading{
				exchange: _exchagne,
				symbol:   _symbol,
			}
			tradings = append(tradings, _trading)
		}
	}

	// init log
	logger.InitLogger(nil, "", 10, 5, 30, true)
	for _, t := range tradings {
		logger.InitLogger(t.exchange, t.symbol, 10, 5, 30, true)
		logger.Print(t.exchange, t.symbol, "Time, PriceNow, Open, Close, High, Low, SMA, EMA, BBands+3, BBands+2, BBands-2, BBands-3, K, D, SMASlope, RSI, BUY, SELL")
	}

	// init db
	db.InitDB()
	defer db.CloseDB()

	// TOOD: recovery from AppErr

	wg := new(sync.WaitGroup)
	for {
		wg.Add(len(tradings))

		for _, t := range tradings {
			localT := t // 闭包变量捕获问题
			go func(wg *sync.WaitGroup) {
				// defer handlePanic()
				defer wg.Done()

				since := time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC)
				test := localT.exchange.FetchCandleSticks(since, string(localT.symbol), time.Duration(3*time.Minute))
				fmt.Println(localT.symbol, test[1].OpenTime, test[1].Open, test[1].High, test[1].Low, test[1].Close, test[1].Volume)
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
