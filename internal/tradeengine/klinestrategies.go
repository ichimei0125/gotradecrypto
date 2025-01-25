package tradeengine

import "github.com/ichimei0125/gotradecrypto/internal/exchange"

// 策略1
//
//	data[2].bbands in, data[1].bbands out
//	AND
//	data[1] slowD <20(for buy)/ >80(for sell) OR slowK > slowD(for buy)
func Tradestrategy(data []exchange.CandleStick) TradeStatus {

	const stoch_down float64 = 25
	const stoch_up float64 = 75

	var last = len(data) - 2
	var now = len(data) - 1

	// const rsi_down float64 = 35
	// const rsi_up float64 = 65

	// for buy
	// if (d[1].SlowD < stoch_down || d[1].SlowK < stoch_down) && d[1].RSI < rsi_down {
	if data[last].SlowD < stoch_down || data[last].SlowK < stoch_down {
		// if d[0].SlowK > d[0].SlowD {
		// 	return BUY
		// }
		if data[last].Close <= data[last].BBands_Minus_2K && data[now].Close >= data[now].BBands_Minus_2K { // && d[last].SMASlope < now {
			return BUY
		}

	}

	// for sell
	// if (d[1].SlowD > stoch_up || d[1].SlowK > stoch_up) && d[1].RSI > rsi_up {
	if data[last].SlowD > stoch_up || data[last].SlowK > stoch_up {
		// if data[now].SlowK < data[now].SlowD {
		// 	return SELL
		// }
		if data[last].Close >= data[last].BBands_Plus_2K && data[now].Close <= data[now].BBands_Plus_2K { // && data[last].SMASlope > now {
			return SELL
		}
	}

	return HOLD
}
