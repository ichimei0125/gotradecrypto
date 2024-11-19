package indicator

import "github.com/ichimei0125/gotradecrypto/internal/exchange"

func GetIndicators(data *[]exchange.KLine) {
	SMA(data, 20)
	EMA(data, 20)
	BBands(data, 20)
}
