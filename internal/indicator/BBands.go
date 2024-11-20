package indicator

import "github.com/ichimei0125/gotradecrypto/internal/exchange"

func bbands(data *[]exchange.KLine, period int) {
	d := *data
	for i := len(d) - period; i >= 0; i -= 1 {
		if d[i].SMA == 0.0 {
			continue
		}
		closes := exchange.GetClose(d[i : i+period])
		d[i].BBands_Plus_2K = d[i].SMA + 2*standardDeviation(closes)
		d[i].BBands_Plus_3K = d[i].SMA + 3*standardDeviation(closes)

		d[i].BBands_Minus_2K = d[i].SMA - 2*standardDeviation(closes)
		d[i].BBands_Minus_3K = d[i].SMA - 3*standardDeviation(closes)
	}
}
