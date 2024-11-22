package trade

import (
	"github.com/ichimei0125/gotradecrypto/internal/config"
	"github.com/ichimei0125/gotradecrypto/internal/exchange"
)

func buy(e exchange.Exchange, symbol exchange.Symbol, data *[]exchange.KLine) {
	d := *data
	price := d[0].Close

	c := config.GetConfig()
	_invest := c.Trade.InvestMoney
	size := float64(_invest) / price

	if e.GetOrderNum(symbol, exchange.ACTIVE, 10, exchange.BUY) == 0 && e.GetOrderNum(symbol, exchange.COMPLETED, 10, exchange.BUY) == 0 {
		e.BuyCypto(symbol, size, price)
	}
}

func sell(e exchange.Exchange, symbol exchange.Symbol, data *[]exchange.KLine) {
	d := *data
	price := d[0].Close

	c := config.GetConfig()
	_invest := c.Trade.InvestMoney
	size := float64(_invest) / price

	if e.GetOrderNum(symbol, exchange.ACTIVE, 10, exchange.SELL) == 0 && e.GetOrderNum(symbol, exchange.COMPLETED, 10, exchange.SELL) == 0 {
		e.SellCypto(symbol, size, price)
	}
}

func losscut(e exchange.Exchange) bool {
	amount, _ := e.GetBalance(exchange.JPY)
	c := config.GetConfig()
	if amount < float64(c.Trade.CutLoss) {
		e.SellAllCypto()
		return true
	}
	return false
}
