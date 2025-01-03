package exchange

import (
	"fmt"

	"github.com/ichimei0125/gotradecrypto/internal/config"
)

type Exchange interface {
	Name() string
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

	CheckUnfinishedOrder(symbol Symbol)
}

type Balance string

const (
	JPY  = "JPY"
	BTC  = "BTC"
	ETH  = "ETH"
	XRP  = "XRP"
	XLM  = "XLM"
	MONA = "MONA"
	BCH  = "BCH"
)

type Symbol string

const (
	BTCJPY  Symbol = "BTC_JPY"
	XRPJPY  Symbol = "XRP_JPY"
	ETHJPY  Symbol = "ETH_JPY"
	XLMJPY  Symbol = "XLM_JPY"
	MONAJPY Symbol = "MONA_JPY"
	// ETHBTC    Symbol = "ETH_BTC"
	// BCHBTC    Symbol = "BCH_BTC"
	FX_BTCJPY Symbol = "FX_BTC_JPY"
)

// Return: Sell coin, Buy money (虚拟币， 法币)
func (s *Symbol) GetTradePair() (Balance, Balance) {
	tradePairMap := map[Symbol][2]Balance{
		BTCJPY:  {BTC, JPY},
		XRPJPY:  {XRP, JPY},
		ETHJPY:  {ETH, JPY},
		XLMJPY:  {XLM, JPY},
		MONAJPY: {MONA, JPY},
		// ETHBTC:    {ETH, BTC},
		// BCHBTC:    {BCH, BTC},
		FX_BTCJPY: {BTC, JPY},
	}

	if pair, exists := tradePairMap[*s]; exists {
		return pair[0], pair[1]
	}

	panic(fmt.Sprintf("no symbol: %s", *s))
}

func (s *Symbol) IsDryRun(exchangeName string) bool {
	dry_run := config.GetConfig().DryRun
	symbols, exist := dry_run[exchangeName]
	if !exist {
		panic(fmt.Sprintf("no exchange in config.yaml %s", exchangeName))
	}

	if is_dry_run, _exist := symbols[string(*s)]; _exist {
		return is_dry_run
	} else {
		panic(fmt.Sprintf("no symbol in config.yaml %s, exchange %s", string(*s), exchangeName))
	}
}

func (s *Symbol) IsMargin() bool {
	return *s == FX_BTCJPY
}

func GetSecret(exchangeName string) (string, string) {
	secrets := config.GetConfig().Secrets
	if secret, exist := secrets[exchangeName]; exist {
		return secret.ApiKey, secret.ApiSecret
	} else {
		panic(fmt.Sprintf("no exchange in config.yaml, %s", exchangeName))
	}
}

type OrderStatus string

const (
	ACTIVE    OrderStatus = "active"    // オープンな注文
	COMPLETED OrderStatus = "completed" // 全額が取引完了した注文
	CANCELED  OrderStatus = "canceled"  // お客様がキャンセルした注文
	REJECTED  OrderStatus = "rejected"  // 有効期限に到達したため取り消された注文
	EXPIRED   OrderStatus = "expired"   // 失敗した注文
)

type Side string

const (
	BUY  Side = "buy"
	SELL Side = "sell"
)
