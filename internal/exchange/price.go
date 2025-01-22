package exchange

import (
	"fmt"
	"slices"
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/common"
)

type CandleStick struct {
	Open     float64
	Close    float64
	High     float64
	Low      float64
	Volume   float64
	OpenTime time.Time

	// indicators
	SMA             float64
	EMA             float64
	BBands_Plus_3K  float64
	BBands_Plus_2K  float64
	BBands_Minus_3K float64
	BBands_Minus_2K float64
	SlowK           float64
	SlowD           float64
	SMASlope        float64
	RSI             float64
}

type Trade struct {
	ID            string    `gorm:"primaryKey"`
	Side          string    `gorm:"not null;type:varchar(4)"`
	Price         float64   `gorm:"not null;type:float"`
	Size          float64   `gorm:"not null;type:float"`
	ExecutionTime time.Time `gorm:"not null;type:datetime(6);index"`
}

func SumClose(klines []CandleStick) float64 {
	var sum float64
	for _, k := range klines {
		sum += k.Close
	}
	return sum
}

func GetClose(KLine []CandleStick) []float64 {
	res := make([]float64, len(KLine))
	for i, k := range KLine {
		res[i] = k.Close
	}
	return res
}

// trades: sorted new -> old
func TradesToCandleStickByMinute(trades []Trade, minutes int) []CandleStick {
	if !(minutes > 0 && minutes < 60) {
		panic(fmt.Sprintf("minutes must in 1~59, minutes: %d", minutes))
	}

	lastest_time := trades[0].ExecutionTime
	norm_minute := int(lastest_time.Minute()/minutes) * minutes
	start_time := time.Date(lastest_time.Year(), lastest_time.Month(), lastest_time.Day(), lastest_time.Hour(), norm_minute, 0, 0, lastest_time.Location())

	return convertTradesToCandleStick(trades, start_time, time.Duration(minutes*int(time.Minute)))

}

// trades: sorted new -> old
func convertTradesToCandleStick(trades []Trade, start_time time.Time, duration time.Duration) []CandleStick {
	candlesticks := []CandleStick{}

	prices := []float64{}
	vol := 0.0
	for _, trade := range trades {
		if trade.ExecutionTime.After(start_time) {
			prices = append(prices, trade.Price)
			vol += trade.Size
		} else {
			// it is possible no trading
			if len(prices) > 0 {
				// append to candlestick
				candlesticks = append(candlesticks, CandleStick{
					Open:     prices[len(prices)-1],
					Close:    prices[0],
					High:     slices.Max(prices),
					Low:      slices.Min(prices),
					Volume:   common.Round(vol, 3),
					OpenTime: start_time,
				})
			}
			// reset tmp vals
			prices = []float64{}
			vol = 0.0
			start_time = start_time.Add(-duration)
		}
	}

	return candlesticks
}
