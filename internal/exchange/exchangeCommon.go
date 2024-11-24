package exchange

import "fmt"

type Exchange interface {
	// Public
	FetchKLine(symbol Symbol, cache *[]KLine)
	// Private
	BuyCypto(symbol Symbol, size float64, price float64)
	SellCypto(symbol Symbol, size float64, price float64)
	SellAllCypto()

	GetBalance(b Balance) (float64, float64)                                  // amount, avaiable (資産総額,発注中以外の金額)
	GetOrderNum(symbol Symbol, status OrderStatus, minues int, side Side) int // number of special orderstatus
	CancelAllOrder(symbol Symbol)
	GetTradeSizeLimit(symbol Symbol) float64
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

// Return: Sell coin, Buy money (虚拟币， 法币)
func GetTradePair(symbol Symbol) (Balance, Balance) {
	tradePairMap := map[Symbol][2]Balance{
		BTCJPY:    {BTC, JPY},
		XRPJPY:    {XRP, JPY},
		ETHJPY:    {ETH, JPY},
		XLMJPY:    {XLM, JPY},
		MONAJPY:   {MONA, JPY},
		ETHBTC:    {ETH, BTC},
		BCHBTC:    {BCH, BTC},
		FX_BTCJPY: {FX_BTC, FX_JPY},
	}

	if pair, exists := tradePairMap[symbol]; exists {
		return pair[0], pair[1]
	}

	panic(fmt.Sprintf("no symbol: %s", symbol))
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
	BUY  Side = "buy"
	SELL Side = "sell"
)
