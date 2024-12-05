package trade

import "github.com/ichimei0125/gotradecrypto/internal/exchange"

// 策略1
//
//	data[2].bbands in, data[1].bbands out
//	AND
//	data[1] slowD <20(for buy)/ >80(for sell) OR slowK > slowD(for buy)
func tradestrategy(data *[]exchange.KLine) TradeStatus {
	d := *data

	const stoch_down float64 = 25
	const stoch_up float64 = 75

	// const rsi_down float64 = 35
	// const rsi_up float64 = 65

	// for buy
	// if (d[1].SlowD < stoch_down || d[1].SlowK < stoch_down) && d[1].RSI < rsi_down {
	if d[1].SlowD < stoch_down || d[1].SlowK < stoch_down {
		// if d[0].SlowK > d[0].SlowD {
		// 	return BUY
		// }
		if d[1].Close <= d[1].BBands_Minus_2K && d[0].Close >= d[0].BBands_Minus_2K { // && d[1].SMASlope < 0 {
			return BUY
		}

	}

	// for sell
	// if (d[1].SlowD > stoch_up || d[1].SlowK > stoch_up) && d[1].RSI > rsi_up {
	if d[1].SlowD > stoch_up || d[1].SlowK > stoch_up {
		// if d[0].SlowK < d[0].SlowD {
		// 	return SELL
		// }
		if d[1].Close >= d[1].BBands_Plus_2K && d[0].Close <= d[0].BBands_Plus_2K { // && d[1].SMASlope > 0 {
			return SELL
		}
	}

	return DoNothing
}
