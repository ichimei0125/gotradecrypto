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

// mode: "open", "close", "high", "low"
func GetPrice(candlesticks []CandleStick, mode string) []float64 {
	res := make([]float64, len(candlesticks))

	switch mode {
	case "close":
		for i, k := range candlesticks {
			res[i] = k.Close
		}
	case "open":
		for i, k := range candlesticks {
			res[i] = k.Open
		}
	case "high":
		for i, k := range candlesticks {
			res[i] = k.High
		}
	case "low":
		for i, k := range candlesticks {
			res[i] = k.Low
		}
	default:
		panic("mode must be \"open\", \"close\", \"high\", \"low\"")
	}

	return res
}

// trades: sorted old -> new
func TradesToCandleStickByMinute(trades []Trade, minutes int) []CandleStick {
	if !(minutes > 0 && minutes < 60) {
		panic(fmt.Sprintf("minutes must in 1~59, minutes: %d", minutes))
	}

	openTime := trades[0].ExecutionTime
	norm_minute := int(openTime.Minute()/minutes) * minutes
	openTime = time.Date(openTime.Year(), openTime.Month(), openTime.Day(), openTime.Hour(), norm_minute, 0, 0, openTime.Location())

	return convertTradesToCandleStick(trades, openTime, time.Duration(minutes*int(time.Minute)))

}

// trades: sorted old -> new
func convertTradesToCandleStick(trades []Trade, openTime time.Time, duration time.Duration) []CandleStick {
	candlesticks := make([]CandleStick, 0, 1200)

	nextOpenTime := openTime.Add(duration)
	prices := []float64{}
	vol := 0.0
	for i, trade := range trades {
		if trade.ExecutionTime.Before(nextOpenTime) {
			prices = append(prices, trade.Price)
			vol += trade.Size
			// the last piece will not continue
			if i < len(trades)-1 {
				continue
			}
		}
		// it is possible no trading
		if len(prices) > 0 {
			// append to candlestick
			candlesticks = append(candlesticks, CandleStick{
				Open:     prices[len(prices)-1],
				Close:    prices[0],
				High:     slices.Max(prices),
				Low:      slices.Min(prices),
				Volume:   common.Round(vol, 3),
				OpenTime: openTime,
			})
		}
		// reset tmp vals
		prices = []float64{}
		vol = 0.0
		openTime = nextOpenTime
		nextOpenTime = openTime.Add(duration)
	}

	return candlesticks
}
