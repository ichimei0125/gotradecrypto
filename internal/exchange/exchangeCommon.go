package exchange

type Exchange interface {
	// Public
	FetchKLine(symbol Symbol, cache *[]KLine)
	// Private
	// BuyCypto(jpy float64)
	// SellCypto(jpy float64) // 0 for sell all

	GetBalance(b Balance) (float64, float64)
	// GetTradeHistory
	// GetOrder()
	// CancelOrder()
}

type Balance string

const (
	JPY = "JPY"
	BTC = "BTC"
	ETH = "ETH"
)

type Symbol string

const (
	BTCJPY    Symbol = "BTC_JPY"
	XRPJPY    Symbol = "XRP_JPY"
	ETHJPY    Symbol = "ETH_JPY"
	XLMJPY    Symbol = "XLM_JPY"
	MONAJPY   Symbol = "MONA_JPY"
	ETHBTC    Symbol = "ETH_BTC"
	BCHBTC    Symbol = "BCH_BTC"
	FX_BTCJPY Symbol = "FX_BTC_JPY"
)
