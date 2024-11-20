package indicator

import (
	"math"

	"github.com/ichimei0125/gotradecrypto/internal/exchange"
)

func stochastic(data *[]exchange.KLine) {
	d := *data
	period := 14
	kSmoothingPeriod := 1
	DSmoothingPeriod := 3

	RSV := make([]float64, len(d))
	L := make([]float64, len(d))
	H := make([]float64, len(d))
	start_index := len(d) - period
	for i := start_index; i >= 0; i-- {
		lowest := d[i].Low
		highest := d[i].High
		for j := i + 1; j <= i+period; j++ {
			if d[j].Low < lowest {
				lowest = d[j].Low
			}
			if d[j].High > highest {
				highest = d[j].High
			}
		}
		L[i] = lowest
		H[i] = highest

		RSV[i] = (d[i].Close - L[i]) / (H[i] - L[i]) * 100
	}

	alpha_K := math.Round(2.0/float64(kSmoothingPeriod+1)*1000) / 1000 // 小数点後３位
	alpha_D := math.Round(2.0/float64(DSmoothingPeriod+1)*1000) / 1000 // 小数点後３位
	for i := start_index; i >= 0; i-- {
		if d[i+1].SlowK == 0.0 {
			d[i+1].SlowK = 0.5
		}
		if d[i+1].SlowD == 0.0 {
			d[i+1].SlowD = 0.5
		}

		d[i].SlowK = alpha_K*RSV[i] + (1-alpha_K)*d[i+1].SlowK
		d[i].SlowD = alpha_D*d[i].SlowK + (1-alpha_D)*d[i+1].SlowD
	}
}
