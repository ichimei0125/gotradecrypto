package indicator

import (
	"github.com/ichimei0125/gotradecrypto/internal/exchange"
	"github.com/markcheno/go-talib"
)

func GetIndicators(data *[]exchange.CandleStick) {
	d := *data

	closes := exchange.GetPrice(d, "close")
	highs := exchange.GetPrice(d, "high")
	lows := exchange.GetPrice(d, "low")
	// opens := exchange.GetPrice(d, "open")

	sma := talib.Sma(closes, 20)
	ema := talib.Ema(closes, 20)
	bbands_2_upper, _, bbands_2_lower := talib.BBands(closes, 20, 2, 2, talib.SMA)
	bbands_3_upper, _, bbands_3_lower := talib.BBands(closes, 20, 3, 3, talib.SMA)
	stochs_k, stochs_d := talib.Stoch(highs, lows, closes, 14, 3, talib.SMA, 3, talib.SMA)
	rsis := talib.Rsi(closes, 14)

	for i := 0; i < len(d); i++ {
		d[i].SMA = sma[i]
		d[i].EMA = ema[i]
		d[i].BBands_Plus_2K = bbands_2_upper[i]
		d[i].BBands_Minus_2K = bbands_2_lower[i]
		d[i].BBands_Plus_3K = bbands_3_upper[i]
		d[i].BBands_Minus_3K = bbands_3_lower[i]
		d[i].SlowK = stochs_k[i]
		d[i].SlowD = stochs_d[i]
		d[i].RSI = rsis[i]
	}
}
