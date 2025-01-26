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
	"github.com/ichimei0125/gotradecrypto/internal/indicator"
	"github.com/ichimei0125/gotradecrypto/internal/logger"
	"github.com/ichimei0125/gotradecrypto/internal/simulator"
	"github.com/ichimei0125/gotradecrypto/internal/tradeengine"
	"github.com/spf13/cobra"
)

// exchange to trading
var exchanges = []exchange.Exchange{new(bitflyer.Bitflyer)}

var (
	param_simulate_startdata time.Time
)

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

func tradebot(mode string) {
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
	logger.InitLogger("", "", 10, 5, 30, true)
	for _, t := range tradings {
		logger.InitLogger(t.exchange.GetInfo().Name, t.symbol, 10, 5, 30, true)
		// logger.Print(t.exchange, t.symbol, "Time, PriceNow, Open, Close, High, Low, SMA, EMA, BBands+3, BBands+2, BBands-2, BBands-3, K, D, SMASlope, RSI, BUY, SELL")
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

				if mode == "trade" {
					since := common.GetUTCNow().Add(time.Duration(-common.KLINE_INTERVAL * (common.KLINE_LENGTH + 2) * int(time.Minute)))
					candlesticks := localT.exchange.FetchCandleSticks(since, string(localT.symbol), time.Duration(3*time.Minute))
					indicator.GetIndicators(&candlesticks)

					go func() {
						// defer handlePanic()
						tradeengine.Trade(localT.exchange, localT.symbol, candlesticks)
					}()

				} else if mode == "simulate" {
					simulator.Simulator(localT.exchange, localT.symbol, param_simulate_startdata)
				} else if mode == "updatedata" {
					localT.exchange.FetchTrades(common.NULLDATE, localT.symbol)

				} else {
					panic(fmt.Sprintf("wrong mode: %s", mode))
				}

			}(wg)
		}

		wg.Wait()
		if mode == "trade" {
			time.Sleep(time.Duration(common.REFRESH_INTERVAL) * time.Minute)
		} else {
			break
		}
	}
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "tradebot",
		Short: "Trading Bot CLI",
		Long:  "A CLI tool to run a trading bot or backtest a strategy",
		Run: func(cmd *cobra.Command, args []string) {
			tradebot("trade")
		},
	}

	var simulateCmd = &cobra.Command{
		Use:   "simulate",
		Short: "Run backtesting simulator",
		Long:  "Run the backtesting simulator",
		Run: func(cmd *cobra.Command, args []string) {
			days, _ := cmd.Flags().GetInt("days")
			param_simulate_startdata = time.Now().AddDate(0, 0, -days).UTC()

			tradebot("simulate")
		},
	}

	var updatedataCmd = &cobra.Command{
		Use:   "updatedata",
		Short: "update trades data",
		Long:  "update trades data for exchange",
		Run: func(cmd *cobra.Command, args []string) {
			tradebot("updatedata")
		},
	}

	// init flag
	simulateCmd.Flags().Int("days", 90, "Number of days in the past to start the backtest")

	rootCmd.AddCommand(simulateCmd)
	rootCmd.AddCommand(updatedataCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
