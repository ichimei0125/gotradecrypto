package trade

import (
	"fmt"
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/config"
	"github.com/ichimei0125/gotradecrypto/internal/exchange"
	"github.com/ichimei0125/gotradecrypto/internal/logger"
)

type TradeStatus int

const (
	BUY TradeStatus = iota
	SELL
	DoNothing
)

func Trade(e exchange.Exchange, symbol exchange.Symbol, data *[]exchange.KLine) {
	e.CheckUnfinishedOrder(symbol)

	d := *data
	c := config.GetConfig()

	var output_buy, output_sell int = 0, 0

	if !c.DryRun && losscut(e, symbol) {
		output_buy, output_sell = -1, -1
		return
	}

	status := tradestrategy(data)

	switch status {
	case BUY:
		output_buy, output_sell = 1, 0
		if !c.DryRun {
			buy(e, symbol, data)
		}
	case SELL:
		output_buy, output_sell = 0, 1
		if !c.DryRun {
			sell(e, symbol, data)
		}
	case DoNothing:
		output_buy, output_sell = 0, 0
	}

	// check unfinished order

	logger.Info(e, symbol, fmt.Sprintf("%s, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %d, %d", time.Now(), d[0].Close, d[1].Open, d[1].Close, d[1].High, d[1].Low, d[1].SMA, d[1].EMA, d[1].BBands_Plus_3K, d[1].BBands_Plus_2K, d[1].BBands_Minus_2K, d[1].BBands_Minus_3K, d[1].SlowK, d[1].SlowD, d[1].SMASlope, d[1].RSI, output_buy, output_sell))
}
