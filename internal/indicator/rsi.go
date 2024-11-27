package indicator

import (
	"fmt"

	"github.com/ichimei0125/gotradecrypto/internal/exchange"
)

func rsi(data *[]exchange.KLine, period int) {
	d := *data

	// U := make([]float64, len(d)-1)
	// D := make([]float64, len(d)-1)

	// for i := len(d) - 2; i >= 0; i-- {
	// 	sub := d[i].Close - d[i+1].Close
	// 	if sub > 0 {
	// 		U[i] = sub
	// 	} else if sub < 0 {
	// 		D[i] = sub
	// 	}
	// }

	// emaU := calculateEMA(U, period)
	// emaD := calculateEMA(D, period)

	// for i := range d {
	// 	if i > len(d)-period-2 {
	// 		break
	// 	}
	// 	d[i].RSI = emaU[i] / (emaU[i] + emaD[i]) * 100
	// }

	closes := exchange.GetClose(d)
	rsi := calculateRSI(closes, period)
	for i := 0; i < len(d)-period; i++ {
		d[i].RSI = rsi[i]
	}

}

// calculateRSI calculates the RSI for a given slice of prices and period.
func calculateRSI(prices []float64, period int) []float64 {
	if len(prices) < period+1 {
		panic(fmt.Errorf("not enough data points to calculate RSI: need at least %d, got %d", period+1, len(prices)))
	}

	// Reverse the array since index 0 represents today.
	reversedPrices := reverseArray(prices)

	var rsi []float64
	var gains, losses float64

	// Calculate the first average gains and losses
	for i := 1; i <= period; i++ {
		change := reversedPrices[i] - reversedPrices[i-1]
		if change > 0 {
			gains += change
		} else {
			losses -= change
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	// Append the first RSI value
	if avgLoss == 0 {
		rsi = append(rsi, 100)
	} else {
		rs := avgGain / avgLoss
		rsi = append(rsi, 100-(100/(1+rs)))
	}

	// Calculate subsequent RSI values
	for i := period + 1; i < len(reversedPrices); i++ {
		change := reversedPrices[i] - reversedPrices[i-1]
		if change > 0 {
			gain := change
			loss := 0.0
			avgGain = ((avgGain * float64(period-1)) + gain) / float64(period)
			avgLoss = ((avgLoss * float64(period-1)) + loss) / float64(period)
		} else {
			gain := 0.0
			loss := -change
			avgGain = ((avgGain * float64(period-1)) + gain) / float64(period)
			avgLoss = ((avgLoss * float64(period-1)) + loss) / float64(period)
		}

		if avgLoss == 0 {
			rsi = append(rsi, 100)
		} else {
			rs := avgGain / avgLoss
			rsi = append(rsi, 100-(100/(1+rs)))
		}
	}

	// Reverse the RSI back to match the original price array format
	return reverseArray(rsi)
}
