package exchange

import (
	"time"
)

type KLine struct {
	Symbol    Symbol
	Open      float64
	Close     float64
	High      float64
	Low       float64
	Volume    float64
	CloseTime time.Time

	// indicators
	SMA             float64
	SMA_new         float64
	EMA             float64
	BBands_Plus_3K  float64
	BBands_Plus_2K  float64
	BBands_Minus_3K float64
	BBands_Minus_2K float64
	SlowK           float64
	SlowD           float64
}

func SumClose(klines []KLine) float64 {
	var sum float64
	for _, k := range klines {
		sum += k.Close
	}
	return sum
}

func GetClose(KLine []KLine) []float64 {
	res := make([]float64, len(KLine))
	for i, k := range KLine {
		res[i] = k.Close
	}
	return res
}
