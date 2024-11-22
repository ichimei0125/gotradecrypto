package exchange

import "fmt"

type Exchange interface {
	// Public
	FetchKLine(symbol Symbol, cache *[]KLine)
	// Private
	BuyCypto(symbol Symbol, size float64, price float64)
	SellCypto(symbol Symbol, size float64, price float64)
	SellAllCypto()

	GetBalance(b Balance) (float64, float64)                                  // amount, avaiable
	GetOrderNum(symbol Symbol, status OrderStatus, minues int, side Side) int // number of special orderstatus
	CancelAllOrder(symbol Symbol)
}

type Balance string

const (
	JPY    = "JPY"
	FX_JPY = "FX_JPY"
	BTC    = "BTC"
	ETH    = "ETH"
	XRP    = "XRP"
	XLM    = "XLM"
	MONA   = "MONA"
	BCH    = "BCH"
	FX_BTC = "FX_BTC"
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

// Return: Sell coin, Buy money
func GetTradePair(symbol Symbol) (Balance, Balance) {
	switch symbol {
	case BTCJPY:
		return BTC, JPY
	default:
		panic(fmt.Sprintf("no symbol: %s", symbol))
	}
}

type OrderStatus string

const (
	ACTIVE    OrderStatus = "active"
	COMPLETED OrderStatus = "completed"
	CANCELED  OrderStatus = "canceled"
	REJECTED  OrderStatus = "rejected"
)

type Side string

const (
	BUY  Side = "side"
	SELL Side = "sell"
)
