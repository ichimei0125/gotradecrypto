package exchange

import (
	"fmt"
	"sync"
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/common"
	"github.com/ichimei0125/gotradecrypto/internal/config"
)

type Exchange interface {
	GetInfo() ExchangeInfo
	// Public
	FetchTrades(since time.Time, symbol string) []Trade
	FetchCandleSticks(since time.Time, symbol string, interval time.Duration) []CandleStick
	// Private
	BuyCypto(symbol string, size float64, price float64)
	SellCypto(symbol string, size float64, price float64)
	SellAllCypto()

	GetBalance(b string) (float64, float64)                                   // amount, avaiable (資産総額,発注中以外の金額)
	GetOrderNum(symbol string, status OrderStatus, minues int, side Side) int // number of special orderstatus
	CancelAllOrder(symbol string)

	CheckUnfinishedOrder(symbol string)
}

type ExchangeInfo struct {
	Name       string
	Symbols    []string
	IsDryRun   bool
	IsRealtime bool
	Waittime   time.Duration
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
	mu       sync.Mutex
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
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := common.GetNow()

	// reset
	if now.After(rl.endTime) {
		rl.endTime = now.Add(rl.interval)
		rl.count = 1
		return true, 0
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
