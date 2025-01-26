package simulator

import (
	"fmt"
	"sync"
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/common"
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

type balance struct {
	money float64
	coin  float64
}

var (
	balanceStore sync.Map
)

func Simulator(e exchange.Exchange, symbol string, startTime time.Time) {

	trades := e.FetchTrades(startTime, symbol)
	// TODO: realtime
	waitDuration := e.GetInfo().Waittime

	endTime := trades[0].ExecutionTime.Add(time.Duration(common.KLINE_INTERVAL * (common.KLINE_LENGTH + 2) * int(time.Minute)))
	_trades := []exchange.Trade{}
	for _, trade := range trades {
		if trade.ExecutionTime.Before(endTime) {
			_trades = append(_trades, trade)
			continue
		}
		candlesticks := exchange.TradesToCandleStickByMinute(_trades, common.KLINE_INTERVAL)
		indicator.GetIndicators(&candlesticks)
		tradestatus := tradeengine.Tradestrategy(candlesticks)
		key := common.GetUniqueName(e.GetInfo().Name, symbol)
		if tradestatus == tradeengine.BUY {
			buy(candlesticks[len(candlesticks)-1], e.GetInfo().Name, symbol, key)
		}
		if tradestatus == tradeengine.SELL {
			sell(candlesticks[len(candlesticks)-1], e.GetInfo().Name, symbol, key)
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

func sell(c exchange.CandleStick, e, s string, key string) {
	var b balance
	_b, ok := balanceStore.Load(key)
	if ok {
		b = _b.(balance)
	} else {
		return
	}

	if b.coin < 0.00000001 {
		return
	}

	b.money += b.coin * c.Close
	b.coin = 0.0
	if b.money < loss_cut {
		msg := fmt.Sprintf("!!!loss cut!!!%s, %.2f, %.5f", c.OpenTime.Format(time.DateTime), b.money, b.coin)
		logger.Print(e, s, msg)
	}

	msg := fmt.Sprintf("%s, %.2f, %.5f", c.OpenTime.Format(time.DateTime), b.money, b.coin)
	logger.Print(e, s, msg)

	balanceStore.Store(key, b)

}

func buy(c exchange.CandleStick, e, s string, key string) {
	var b balance
	_b, ok := balanceStore.Load(key)
	if ok {
		b = _b.(balance)
	} else {
		b = balance{
			money: float64(init_money),
			coin:  0,
		}
	}

	if b.money >= invest_money {
		b.coin += invest_money / c.Close
		b.money -= invest_money

		msg := fmt.Sprintf("%s, %.2f, %.5f", c.OpenTime.Format(time.DateTime), b.money, b.coin)
		logger.Print(e, s, msg)
	}

	balanceStore.Store(key, b)
}
