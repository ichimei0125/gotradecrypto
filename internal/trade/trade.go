package trade

import (
	"log"

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

	if losscut(e, symbol) {
		log.Printf(",%s, %f, %f, %f, %f, %f, %f, %f, %f, %f, %d, %d", d[1].OpenTime, d[1].Open, d[1].SMA, d[1].EMA, d[1].BBands_Plus_3K, d[1].BBands_Plus_2K, d[1].BBands_Minus_2K, d[1].BBands_Minus_3K, d[1].SlowK, d[1].SlowD, -1, -1)
		return
	}

	status := strategy1(data)

	switch status {
	case BUY:
		log.Printf(",%s, %f, %f, %f, %f, %f, %f, %f, %f, %f, %d, %d", d[1].OpenTime, d[1].Open, d[1].SMA, d[1].EMA, d[1].BBands_Plus_3K, d[1].BBands_Plus_2K, d[1].BBands_Minus_2K, d[1].BBands_Minus_3K, d[1].SlowK, d[1].SlowD, 1, 0)
		buy(e, symbol, data)
	case SELL:
		log.Printf(",%s, %f, %f, %f, %f, %f, %f, %f, %f, %f, %d, %d", d[1].OpenTime, d[1].Open, d[1].SMA, d[1].EMA, d[1].BBands_Plus_3K, d[1].BBands_Plus_2K, d[1].BBands_Minus_2K, d[1].BBands_Minus_3K, d[1].SlowK, d[1].SlowD, 0, 1)
		sell(e, symbol, data)
	case DoNothing:
		log.Printf(",%s, %f, %f, %f, %f, %f, %f, %f, %f, %f, %d, %d", d[1].OpenTime, d[1].Open, d[1].SMA, d[1].EMA, d[1].BBands_Plus_3K, d[1].BBands_Plus_2K, d[1].BBands_Minus_2K, d[1].BBands_Minus_3K, d[1].SlowK, d[1].SlowD, 0, 0)
	}

}
