package indicator

import "github.com/ichimei0125/gotradecrypto/internal/exchange"

func GetIndicators(data *[]exchange.KLine) {
	sma(data, 20)
	ema(data, 20)
	bbands(data, 20)
	stochastic(data, 14, 1, 3)
	maSlope(data, 14)
	rsi(data, 14)
}
