package tradeengine

import (
	"github.com/ichimei0125/gotradecrypto/internal/common"
	"github.com/ichimei0125/gotradecrypto/internal/config"
	"github.com/ichimei0125/gotradecrypto/internal/exchange"
)

func buy(e exchange.Exchange, symbol string, data []exchange.CandleStick) {
	price := data[0].Close

	c := config.GetConfig()
	_invest := c.Trade.InvestMoney
	size := float64(_invest) / price

	const sameKLineIntervalBuyTimes int = 2
	if e.GetOrderNum(symbol, exchange.ACTIVE, common.KLINE_INTERVAL, exchange.BUY) < sameKLineIntervalBuyTimes && e.GetOrderNum(symbol, exchange.COMPLETED, common.KLINE_INTERVAL, exchange.BUY) < sameKLineIntervalBuyTimes {
		e.BuyCypto(symbol, size, price)
	}
}

func sell(e exchange.Exchange, symbol string, data []exchange.CandleStick) {
	price := data[0].Close

	c := config.GetConfig()
	_invest := c.Trade.InvestMoney
	size := float64(_invest) / price

	const sameKLineIntervalSellTimes int = 1
	if e.GetOrderNum(symbol, exchange.ACTIVE, common.KLINE_INTERVAL, exchange.SELL) < sameKLineIntervalSellTimes && e.GetOrderNum(symbol, exchange.COMPLETED, common.KLINE_INTERVAL, exchange.SELL) < sameKLineIntervalSellTimes {
		e.SellCypto(symbol, size, price)
	}
}

func losscut(e exchange.Exchange, symbol string) bool {
	// TODO
	return false
	// coin, money := symbol.GetTradePair()

	// // TODO 多币种支持

	// money_amount, _ := e.GetBalance(money)
	// coin_amount, _ := e.GetBalance(coin)
	// size_limit := e.GetTradeSizeLimit(symbol)

	// c := config.GetConfig()
	// if money_amount <= float64(c.Trade.CutLoss) && coin_amount < size_limit*2 {
	// 	e.SellAllCypto()
	// 	return true
	// }
	// return false
}
