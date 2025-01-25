package simulator

import (
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/exchange"
	"github.com/ichimei0125/gotradecrypto/internal/logger"
	"github.com/ichimei0125/gotradecrypto/internal/tradeengine"
)

const (
	init_money   = 50000
	invest_money = 10000 // TODO: more invest strategy
	loss_cut     = 40000
)

var (
	money = 0.0
	coin  = 0.0
)

func Simulator(e exchange.Exchange, symbol string, startDate time.Time) {

	trades := e.FetchTrades(startDate, symbol)
	// TODO: realtime
	waitMin := e.GetInfo().Waittime

	endTime := startDate.Add(time.Duration((tradeengine.CANDLESTICK_LENGTH + 2) * tradeengine.CANDLESTICK_INTERVAL * time.Minute))
	_trades := []exchange.Trade{}
	for _, trade := range trades {
		if trade.ExecutionTime.Before(endTime) {
			_trades = append(_trades, trade)
			continue
		}
		candlesticks := exchange.TradesToCandleStickByMinute(_trades, tradeengine.CANDLESTICK_INTERVAL)
		tradestatus := tradeengine.Tradestrategy(candlesticks)
		if tradestatus == tradeengine.BUY {
			buy(candlesticks[len(candlesticks)-1], e, symbol)
		}
		if tradestatus == tradeengine.SELL {
			sell(candlesticks[len(candlesticks)-1], e, symbol)
		}

		endTime = endTime.Add(waitMin * time.Minute)
		startDate = startDate.Add(waitMin * time.Minute)

		// rm too oldtrade
		_newTrades := make([]exchange.Trade, 0, len(_trades))
		for _, t := range _trades {
			if t.ExecutionTime.After(startDate) {
				_newTrades = append(_newTrades, t)
			}
		}
		_trades = _newTrades
	}

}

func sell(c exchange.CandleStick, e exchange.Exchange, s string) {
	money += coin * c.Close
	if money < loss_cut {
		logger.Print(e, s, "loss cut", money)

	}

	logger.Print(e, s, money)

}

func buy(c exchange.CandleStick, e exchange.Exchange, s string) {
	if money > 0 {
		size := invest_money / c.Close
		coin += size
	}
	logger.Print(e, s, money)
}
