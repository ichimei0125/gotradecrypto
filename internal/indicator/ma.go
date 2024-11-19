package indicator

import (
	"math"

	"github.com/ichimei0125/gotradecrypto/internal/exchange"
)

func SMA(data *[]exchange.KLine, period int) {
	start_index := len(*data) - period
	d := *data

	if d[start_index].SMA == 0.0 {
		d[start_index].SMA = math.Round(exchange.SumClose(d[start_index:])/float64(period)*1000) / 1000 // 小数点後3位
	}

	for i := start_index - 1; i >= 0; i = i - 1 {
		d[i].SMA = d[i+1].SMA + (d[i].Close-d[i+period].Close)/float64(period)
	}
}

func EMA(data *[]exchange.KLine, period int) {
	d := *data
	alpha := math.Round(2.0/float64(period+1)*1000) / 1000 // 小数点後３位

	if d[len(d)-1].EMA == 0.0 {
		d[len(d)-1].EMA = d[len(d)-1].Close
	}

	for i := len(d) - 2; i >= 0; i = i - 1 {
		d[i].EMA = d[i+1].EMA + alpha*(d[i].Close-d[i+1].EMA)
	}
}
