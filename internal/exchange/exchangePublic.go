package exchange

type Exchange interface {
	FetchKLine(symbol Symbol) []KLine
}

// 主逻辑，接受 Exchange 接口
func GetKLineAndIndicators(exchange Exchange, symbol Symbol) []KLine {
	kline := exchange.FetchKLine(symbol)
	return kline
}
