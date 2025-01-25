package indicator

import "github.com/ichimei0125/gotradecrypto/internal/exchange"

func bbands(data *[]exchange.CandleStick, period int) {
	d := *data
	for i := len(d) - period; i >= 0; i -= 1 {
		if d[i].SMA == 0.0 {
			continue
		}
		closes := exchange.GetPrice(d[i:i+period], "close")
		d[i].BBands_Plus_2K = d[i].SMA + 2*standardDeviation(closes)
		d[i].BBands_Plus_3K = d[i].SMA + 3*standardDeviation(closes)

		d[i].BBands_Minus_2K = d[i].SMA - 2*standardDeviation(closes)
		d[i].BBands_Minus_3K = d[i].SMA - 3*standardDeviation(closes)
	}
}

func maSlope(data *[]exchange.CandleStick, period int) {
	d := *data

	for i := 0; i < len(d)-period; i++ {
		tmp_data := make([]float64, period)
		is_nodata := false
		for index, _d := range d[i : i+period] {
			if _d.SMA == 0.0 {
				is_nodata = true
				break
			}
			tmp_data[index] = _d.SMA
		}

		if is_nodata {
			continue
		}
		d[i].SMASlope = slope(tmp_data)
	}

}
