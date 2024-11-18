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
}
