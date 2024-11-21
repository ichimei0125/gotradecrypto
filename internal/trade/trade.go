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

func Trade(data *[]exchange.KLine) {
	status := strategy1(data)

	d := *data
	switch status {
	case BUY:
		log.Printf(",%s, %f, %f, %f, %f, %f, %f, %f, %f, %f, %d, %d", d[1].CloseTime, d[1].Open, d[1].SMA, d[1].EMA, d[1].BBands_Plus_3K, d[1].BBands_Plus_2K, d[1].BBands_Minus_2K, d[1].BBands_Minus_3K, d[1].SlowK, d[1].SlowD, 1, 0)
	case SELL:
		log.Printf(",%s, %f, %f, %f, %f, %f, %f, %f, %f, %f, %d, %d", d[1].CloseTime, d[1].Open, d[1].SMA, d[1].EMA, d[1].BBands_Plus_3K, d[1].BBands_Plus_2K, d[1].BBands_Minus_2K, d[1].BBands_Minus_3K, d[1].SlowK, d[1].SlowD, 0, 1)
	case DoNothing:
		log.Printf(",%s, %f, %f, %f, %f, %f, %f, %f, %f, %f, %d, %d", d[1].CloseTime, d[1].Open, d[1].SMA, d[1].EMA, d[1].BBands_Plus_3K, d[1].BBands_Plus_2K, d[1].BBands_Minus_2K, d[1].BBands_Minus_3K, d[1].SlowK, d[1].SlowD, 0, 0)
	}

}
