package exchange

import (
	"time"
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
	Side          string    `gorm:"not null"`
	Price         float64   `gorm:"not null"`
	Size          float64   `gorm:"not null"`
	ExecutionTime time.Time `gorm:"not null;index"`
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
