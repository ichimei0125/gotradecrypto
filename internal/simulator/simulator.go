package simulator

import (
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/exchange"
	"github.com/ichimei0125/gotradecrypto/internal/indicator"
	"github.com/ichimei0125/gotradecrypto/internal/logger"
	"github.com/ichimei0125/gotradecrypto/internal/tradeengine"
)

const (
	init_money   = 50000
	invest_money = 10000 // TODO: more invest strategy
	loss_cut     = 40000
)

var (
	money = float64(init_money)
	coin  = 0.0
)

func Simulator(e exchange.Exchange, symbol string, startTime time.Time) {

	trades := e.FetchTrades(startTime, symbol)
	// TODO: realtime
	waitDuration := e.GetInfo().Waittime

	endTime := startTime.Add(time.Duration((tradeengine.CANDLESTICK_LENGTH + 2) * tradeengine.CANDLESTICK_INTERVAL * time.Minute))
	_trades := []exchange.Trade{}
	for _, trade := range trades {
		if trade.ExecutionTime.Before(endTime) {
			_trades = append(_trades, trade)
			continue
		}
		candlesticks := exchange.TradesToCandleStickByMinute(_trades, tradeengine.CANDLESTICK_INTERVAL)
		indicator.GetIndicators(&candlesticks)
		tradestatus := tradeengine.Tradestrategy(candlesticks)
		if tradestatus == tradeengine.BUY {
			buy(candlesticks[len(candlesticks)-1], e.GetInfo().Name, symbol)
		}
		if tradestatus == tradeengine.SELL {
			sell(candlesticks[len(candlesticks)-1], e.GetInfo().Name, symbol)
		}
		// logger.Print(e, symbol, candlesticks[len(candlesticks)-1].OpenTime)

		endTime = endTime.Add(waitDuration)
		startTime = startTime.Add(waitDuration)

		// rm too oldtrade
		_newTrades := make([]exchange.Trade, 0, len(_trades))
		for _, t := range _trades {
			if t.ExecutionTime.After(startTime) {
				_newTrades = append(_newTrades, t)
			}
		}
		_trades = _newTrades
	}

}

func sell(c exchange.CandleStick, e, s string) {
	money += coin * c.Close
	coin = 0.0
	if money < loss_cut {
		logger.Print(e, s, "loss cut", money)
	}

	logger.Print(e, s, money)

}

func buy(c exchange.CandleStick, e, s string) {
	if money > 0 {
		coin = invest_money / c.Close
		money -= init_money
	}
	logger.Print(e, s, money)
}
