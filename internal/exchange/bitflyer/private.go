package bitflyer

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/common"
	"github.com/ichimei0125/gotradecrypto/internal/config"
	"github.com/ichimei0125/gotradecrypto/internal/exchange"
)

func bitFlyerPrivateAPICore(path string, method string, body []byte, is_log ...bool) []byte {

	key, secret := exchange.GetSecret(new(Bitflyer).GetInfo().Name)

	timestamp := fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))

	text := timestamp + method + path + string(body)
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(text))
	sign := hex.EncodeToString(hash.Sum(nil))

	url := baseURL + path
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		panic("cannot get bitflyer private api, maybe limited")
	}

	req.Header.Set("ACCESS-KEY", key)
	req.Header.Set("ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("ACCESS-SIGN", sign)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic("wrong private api response")
	}

	if len(is_log) > 0 && is_log[0] {
		log.Printf("bitflyer private rep: %s", string(responseBody))
	}

	return responseBody
}

type balance struct {
	CurrencyCode string  `json:"currency_code"`
	Amount       float64 `json:"amount"`
	Available    float64 `json:"available"`
}

func (b *Bitflyer) GetBalance(ba string) (float64, float64) {
	balances, err := getbalance()
	if err != nil {
		log.Println("Error:", err)
		return -1.0, -1.0
	}

	baVal := ba
	for _, d := range balances {
		if d.CurrencyCode == baVal {
			return d.Amount, d.Available
		}
	}

	return -1.0, -1.0
}

func getbalance() ([]balance, error) {
	path := "/v1/me/getbalance"
	method := "GET"
	res := bitFlyerPrivateAPICore(path, method, nil)

	var balances []balance
	err := json.Unmarshal([]byte(res), &balances)
	return balances, err
}

func (b *Bitflyer) GetOrderNum(symbol string, status exchange.OrderStatus, minues int, side exchange.Side) int {
	cnt := 0
	childorders := getChildOrders(symbol, getOrderStatus(status))

	time_after := common.GetUTCNow().Add(time.Duration(-minues) * time.Minute)

	for _, order := range childorders {
		if order.ChildOrderDate.After(time_after) && order.Side == getSide(side) {
			cnt += 1
		}
	}
	return cnt
}

func getSide(side exchange.Side) string {
	sideMap := map[exchange.Side]string{
		exchange.BUY:  "BUY",
		exchange.SELL: "SELL",
	}

	if side, exist := sideMap[side]; exist {
		return side
	}

	panic(fmt.Sprintf("bitflyer no side %s", side))
}

func getOrderStatus(status exchange.OrderStatus) string {
	statusMap := map[exchange.OrderStatus]string{
		exchange.ACTIVE:    "ACTIVE",
		exchange.COMPLETED: "COMPLETED",
		exchange.CANCELED:  "CANCELED",
		exchange.REJECTED:  "REJECTED",
		exchange.EXPIRED:   "EXPIRED",
	}

	if s, exist := statusMap[status]; exist {
		return s
	}

	panic(fmt.Sprintf("bitflyer err orderstatus %s", status))
}

// 注文状態
type childOrder struct {
	Id                     int64      `json:"id"`
	ChildOrderID           string     `json:"child_order_id"`
	ProductCode            string     `json:"product_code"`
	Side                   string     `json:"side"`
	ChildOrderType         string     `json:"child_order_type"`
	Price                  float64    `json:"price"`
	AveragePrice           float64    `json:"average_price"`
	Size                   float64    `json:"size"`
	ChildOrderState        string     `json:"child_order_state"`
	ExpireDate             CustomTime `json:"expire_date"`
	ChildOrderDate         CustomTime `json:"child_order_date"`
	ChildOrderAcceptanceID string     `json:"child_order_acceptance_id"`
	OutstandingSize        float64    `json:"outstanding_size"`
	CancelSize             float64    `json:"cancel_size"`
	ExecutedSize           float64    `json:"executed_size"`
	TotalCommission        float64    `json:"total_commission"`
	TimeInForce            string     `json:"time_in_force"`
}

type childOrderAcceptanceID struct {
	ChildOrderAcceptanceID string `json:"child_order_acceptance_id"`
}

func getChildOrders(product_code string, child_order_status string) []childOrder {
	endpoint := "/v1/me/getchildorders"
	path := endpoint + "?product_code=" + product_code + "&child_order_state=" + child_order_status

	res := bitFlyerPrivateAPICore(path, "GET", nil)
	var childorders []childOrder
	err := json.Unmarshal([]byte(res), &childorders)
	if err != nil {
		panic("ERROR GetChildOrders: " + err.Error())
	}
	return childorders
}

func getChildOrderByID(product_code string, child_order_acceptance_id string) []childOrder {
	endpoint := "/v1/me/getchildorders"
	path := endpoint + "?product_code=" + product_code + "&child_order_acceptance_id=" + child_order_acceptance_id

	res := bitFlyerPrivateAPICore(path, "GET", nil)
	var childorders []childOrder
	err := json.Unmarshal([]byte(res), &childorders)
	if err != nil {
		panic("ERROR GetChildOrders: " + err.Error())
	}
	return childorders
}

func (b *Bitflyer) CancelAllOrder(symbol string) {
	path := "/v1/me/cancelallchildorders"
	body := fmt.Sprintf(`{
		"product_code": "%s"
	}`, symbol)
	bitFlyerPrivateAPICore(path, "POST", []byte(body))
}

type sendchildorder struct {
	ProductCode    string  `json:"product_code"`
	ChildOrderType string  `json:"child_order_type"`
	Side           string  `json:"side"`
	Price          float64 `json:"price,omitempty"`
	Size           float64 `json:"size"`
	MinuteToExpire int     `json:"minute_to_expire"`
	TimeInForce    string  `json:"time_in_force"`
}

// 自定义 JSON 序列化方法
func (o sendchildorder) MarshalJSON() ([]byte, error) {
	type Alias sendchildorder // 避免递归调用
	if o.ChildOrderType == "MARKET" {
		return json.Marshal(&struct {
			ProductCode    string  `json:"product_code"`
			ChildOrderType string  `json:"child_order_type"`
			Side           string  `json:"side"`
			Size           float64 `json:"size"`
			MinuteToExpire int     `json:"minute_to_expire"`
			TimeInForce    string  `json:"time_in_force"`
		}{
			ProductCode:    o.ProductCode,
			ChildOrderType: o.ChildOrderType,
			Side:           o.Side,
			Size:           math.Round(o.Size*1e6) / 1e6, // 保留6位小数
			MinuteToExpire: o.MinuteToExpire,
			TimeInForce:    o.TimeInForce,
		})
	}

	return json.Marshal(&struct {
		Price float64 `json:"price"` // 保留 Price 字段
		Size  float64 `json:"size"`  // 保留 Size 字段
		Alias
	}{
		Price: math.Round(o.Price*100) / 100, // 保留2位小数
		Size:  math.Round(o.Size*1e6) / 1e6,  // 保留6位小数
		Alias: (Alias)(o),
	})
}

// 売買最小単位
// https://bitflyer.com/ja-jp/s/commission
func getTradeSizeLimit(symbol string) float64 {
	limitMap := map[string]float64{
		"BTC_JPY":    0.001,
		"XRP_JPY":    0.1,
		"ETH_JPY":    0.01,
		"XLM_JPY":    0.1,
		"MONA_JPY":   0.1,
		"ETH_BTC":    0.01,
		"BCH_BTC":    0.01,
		"FX_BTC_JPY": 0.001,
	}
	limit, exist := limitMap[symbol]
	if !exist {
		panic(fmt.Sprintf("bitflyer not support symbol: %s", symbol))
	}
	return limit
}

func getTradePair(symbol string) (string, string) {
	tradePairMap := map[string][]string{
		"BTC_JPY":    {"BTC", "JPY"},
		"XRP_JPY":    {"XRP", "JPY"},
		"ETH_JPY":    {"ETH", "JPY"},
		"XLM_JPY":    {"XLM", "JPY"},
		"MONA_JPY":   {"MONA", "JPY"},
		"ETH_BTC":    {"ETH", "BTC"},
		"BCH_BTC":    {"BCH", "BTC"},
		"FX_BTC_JPY": {"BTC", "JPY"},
	}
	tradePair, exist := tradePairMap[symbol]
	if !exist {
		panic(fmt.Sprintf("bitflyer not support symbol: %s", symbol))
	}
	return tradePair[0], tradePair[1]

}

func (b *Bitflyer) BuyCypto(symbol string, size float64, price float64) {
	limit := getTradeSizeLimit(symbol)

	_, money := getTradePair(symbol)
	_, avaiable := b.GetBalance(money)
	c := config.GetConfig()

	if size < limit || avaiable < float64(c.Trade.SafeMoney) {
		// お金不足
		return
	}
	sendChildOrder(symbol, size, price, "BUY", "LIMIT")

}

func (b *Bitflyer) SellCypto(symbol string, size float64, price float64) {
	// TOOD 考虑size的策略
	coin, _ := getTradePair(symbol)
	_, coin_available := b.GetBalance(coin)

	comission := getTradingCommission(symbol)
	_size := coin_available * (1 - comission)

	_limit := getTradeSizeLimit(symbol)
	if _size < _limit {
		// 資産不足
		return
	}

	order_type := "LIMIT"
	if price <= 0 {
		order_type = "MARKET"
	}

	accept_id := sendChildOrder(symbol, _size, price, "SELL", order_type)
	insertOrder(accept_id.ChildOrderAcceptanceID, symbol, exchange.SELL, _size)
}

func (b *Bitflyer) SellAllCypto() {
	balances, err := getbalance()
	if err != nil {
		log.Fatalf("bitflyer cannot sell all %s", err)
	}

	for _, balance := range balances {
		if balance.CurrencyCode == "JPY" {
			continue
		}
		s := balance.CurrencyCode
		if balance.Available <= getTradeSizeLimit(s)*2 {
			// 資産不足
			continue
		}

		size := balance.Available * (1 - getTradingCommission(s))

		sendChildOrder(s, size, 0, "SELL", "MARKET")
	}
}

type comission struct {
	ComissionRate float64 `json:"commission_rate"`
}

// 手数料
func getTradingCommission(product_code string) float64 {
	path := "/v1/me/gettradingcommission" + "?product_code=" + product_code

	res := bitFlyerPrivateAPICore(path, "GET", nil)
	var _comission comission
	err := json.Unmarshal([]byte(res), &_comission)
	if err != nil {
		fmt.Println("Error: GetTradingCommission", err)
		return 0.0
	}
	return _comission.ComissionRate
}

func sendChildOrder(symbol string, size float64, price float64, side string, ordertype string) childOrderAcceptanceID {
	if size <= 0.0 {
		return childOrderAcceptanceID{}
	}

	order := sendchildorder{
		ProductCode:    symbol,
		ChildOrderType: ordertype,
		Side:           side,
		Price:          price,
		Size:           size,
		MinuteToExpire: common.ORDER_WAIT_MINUTE,
		TimeInForce:    "GTC",
	}

	// 生成 JSON
	jsonData, err := json.MarshalIndent(order, "", "  ")
	if err != nil {
		log.Fatalf("Error: %s", err)
		return childOrderAcceptanceID{}
	}

	path := "/v1/me/sendchildorder"
	rep := bitFlyerPrivateAPICore(path, "POST", []byte(jsonData), true)
	var _childOrderAcceptanceID childOrderAcceptanceID
	err = json.Unmarshal([]byte(rep), &_childOrderAcceptanceID)
	if err != nil {
		fmt.Println("Error: GetTradingCommission", err)
		return _childOrderAcceptanceID
	}
	return _childOrderAcceptanceID
}

type position struct {
	ProductCode         string     `json:"product_code"`          // 产品代码
	Side                string     `json:"side"`                  // 买卖方向
	Price               float64    `json:"price"`                 // 成交价格
	Size                float64    `json:"size"`                  // 成交数量
	Commission          float64    `json:"commission"`            // 佣金
	SwapPointAccumulate float64    `json:"swap_point_accumulate"` // 积累的掉期点
	RequireCollateral   float64    `json:"require_collateral"`    // 所需保证金
	OpenDate            CustomTime `json:"open_date"`             // 开仓时间
	Leverage            float64    `json:"leverage"`              // 杠杆
	PNL                 float64    `json:"pnl"`                   // 盈亏
	SFD                 float64    `json:"sfd"`                   // 闪电掉期费
}

func getPositions(product_code string) ([]position, error) {
	endpoint := "/v1/me/getpositions"
	method := "GET"
	path := endpoint + "?product_code=" + product_code
	res := bitFlyerPrivateAPICore(path, method, nil)

	var positions []position
	err := json.Unmarshal([]byte(res), &positions)
	return positions, err
}
