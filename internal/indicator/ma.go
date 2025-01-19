package indicator

import (
	"fmt"
	"math"

	"github.com/ichimei0125/gotradecrypto/internal/exchange"
)

func sma(data *[]exchange.CandleStick, period int) {
	start_index := len(*data) - period
	d := *data

	if d[start_index].SMA == 0.0 {
		d[start_index].SMA = math.Round(exchange.SumClose(d[start_index:])/float64(period)*1000) / 1000 // 小数点後3位
	}

	for i := start_index - 1; i >= 0; i = i - 1 {
		d[i].SMA = d[i+1].SMA + (d[i].Close-d[i+period].Close)/float64(period)
	}
}

func ema(data *[]exchange.CandleStick, period int) {
	d := *data
	alpha := math.Round(2.0/float64(period+1)*1000) / 1000 // 小数点後３位

	if d[len(d)-1].EMA == 0.0 {
		d[len(d)-1].EMA = d[len(d)-1].Close
	}

	for i := len(d) - 2; i >= 0; i = i - 1 {
		d[i].EMA = d[i+1].EMA + alpha*(d[i].Close-d[i+1].EMA)
	}
}

// func calculateSMA(values []float64) float64 {
// 	sum := 0.0
// 	for _, v := range values {
// 		sum += v
// 	}
// 	return sum / float64(len(values))
// }

func calculateEMA(values []float64, period int) []float64 {
	// ensure cover 99%
	if len(values) < int(2.3*float64(period)+1) {
		panic(fmt.Sprintf("not enough data for calculate ema %d", len(values)))
	}

	alpha := roundAt(2.0/float64(period+1), 3)

	res := make([]float64, len(values))
	res[len(values)-1] = values[len(values)-1]

	for i := len(values) - 2; i >= 0; i = i - 1 {
		res[i] = res[i+1] + alpha*(values[i]-res[i+1])
	}
	return res
}
