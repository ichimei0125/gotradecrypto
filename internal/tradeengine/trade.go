package tradeengine

import (
	"github.com/ichimei0125/gotradecrypto/internal/exchange"
)

type TradeStatus int

const (
	BUY TradeStatus = iota
	SELL
	HOLD
)

const (
	CANDLESTICK_LENGTH   = 1000
	CANDLESTICK_INTERVAL = 3
)

func Trade(e exchange.Exchange, symbol string, data []exchange.CandleStick) {
	e.CheckUnfinishedOrder(symbol)

	// var output_buy, output_sell int = 0, 0

	if !e.GetInfo().IsDryRun && losscut(e, symbol) {
		// output_buy, output_sell = -1, -1
		return
	}

	status := Tradestrategy(data)

	switch status {
	case BUY:
		// output_buy, output_sell = 1, 0
		if !e.GetInfo().IsDryRun {
			buy(e, symbol, data)
		}
	case SELL:
		// output_buy, output_sell = 0, 1
		if !e.GetInfo().IsDryRun {
			sell(e, symbol, data)
		}
	case HOLD:
		// output_buy, output_sell = 0, 0
	}

	// logger.Print(e, symbol, fmt.Sprintf("%s, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %d, %d", time.Now(), d[0].Close, d[1].Open, d[1].Close, d[1].High, d[1].Low, d[1].SMA, d[1].EMA, d[1].BBands_Plus_3K, d[1].BBands_Plus_2K, d[1].BBands_Minus_2K, d[1].BBands_Minus_3K, d[1].SlowK, d[1].SlowD, d[1].SMASlope, d[1].RSI, output_buy, output_sell))
}
