package exchange

import (
	"fmt"
	"sync"
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/common"
	"github.com/ichimei0125/gotradecrypto/internal/config"
)

type Exchange interface {
	Name() string
	// Public
	FetchCandleSticks(symbol Symbol, cache *[]CandleStick)
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

// -----Limitiation-----
var (
	LimitMap sync.Map
)

type RateLimiter struct {
	endTime  time.Time
	interval time.Duration
	count    int
	maxCount int
}

func GetRateLimiter(name string, interval time.Duration, maxCount int) *RateLimiter {
	if rl, exists := LimitMap.Load(name); exists {
		return rl.(*RateLimiter)
	}

	rl := &RateLimiter{
		endTime:  common.GetNow().Add(interval),
		interval: interval,
		count:    0,
		maxCount: maxCount,
	}
	LimitMap.Store(name, rl)
	return rl
}

// Call this, before HTTP request
func (rl *RateLimiter) Allow() (bool, time.Duration) {
	now := common.GetNow()

	// reset
	if now.After(rl.endTime) {
		rl.endTime = now.Add(rl.interval)
		rl.count = 0
	}

	// count (in limit)
	if rl.count < rl.maxCount {
		rl.count++
		return true, 0
	}

	// over limit
	waitTime := rl.endTime.Sub(now)
	return false, waitTime
}

func CleanupRateLimiters(maxIdleTime time.Duration) {
	ticker := time.NewTicker(maxIdleTime)
	defer ticker.Stop()

	for range ticker.C {
		now := common.GetNow()
		LimitMap.Range(func(key, value interface{}) bool {
			rl := value.(*RateLimiter)

			if now.After(rl.endTime.Add(maxIdleTime)) {
				LimitMap.Delete(key)
			}

			return true
		})
	}
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
	ACTIVE    OrderStatus = "active"
	COMPLETED OrderStatus = "completed"
	CANCELED  OrderStatus = "canceled"
	REJECTED  OrderStatus = "rejected"
	EXPIRED   OrderStatus = "expired"
)

type Side string

const (
	BUY  Side = "buy"
	SELL Side = "sell"
	NONE Side = ""
)
