package exchange


// 主逻辑，接受 Exchange 接口
func GetKLine(exchange Exchange, symbol Symbol, cache *[]KLine) {
	exchange.FetchKLine(symbol, cache)
}
