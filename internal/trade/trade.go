package trade

import (
	"fmt"

	"github.com/ichimei0125/gotradecrypto/internal/config"
	"github.com/ichimei0125/gotradecrypto/internal/exchange"
)

type TradeStatus int

const (
	BUY TradeStatus = iota
	SELL
	DoNothing
)

func Trade(e exchange.Exchange, symbol exchange.Symbol, data *[]exchange.KLine) {
	d := *data

	var output_buy, output_sell int = 0, 0

	if losscut(e, symbol) {
		output_buy, output_sell = -1, -1
		return
	}

	status := tradestrategy(data)
	c := config.GetConfig()

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

	fmt.Printf(",%s, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %d, %d", d[1].OpenTime, d[1].Open, d[1].SMA, d[1].EMA, d[1].BBands_Plus_3K, d[1].BBands_Plus_2K, d[1].BBands_Minus_2K, d[1].BBands_Minus_3K, d[1].SlowK, d[1].SlowD, d[1].SMASlope, output_buy, output_sell)

}
