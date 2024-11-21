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
	"net/http"
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/config"
	"github.com/ichimei0125/gotradecrypto/internal/exchange"
)

func bitFlyerPrivateAPICore(path string, method string, body []byte) []byte {

	config := config.GetConfig()
	key := config.Bitflyer.APIKey
	secret := config.Bitflyer.APISecret

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

	return responseBody
}

type balance struct {
	CurrencyCode string  `json:"currency_code"`
	Amount       float64 `json:"amount"`
	Available    float64 `json:"available"`
}

func (b *Bitflyer) GetBalance(ba exchange.Balance) (float64, float64) {
	path := "/v1/me/getbalance"
	method := "GET"

	baVal := getbalancevalue(ba)

	res := bitFlyerPrivateAPICore(path, method, nil)

	var balances []balance
	err := json.Unmarshal([]byte(res), &balances)
	if err != nil {
		log.Println("Error:", err)
		return -1.0, -1.0
	}

	for _, d := range balances {
		if d.CurrencyCode == baVal {
			return d.Amount, d.Available
		}
	}

	return -1.0, -1.0
}

func getbalancevalue(b exchange.Balance) string {
	switch b {
	case exchange.JPY:
		return "JPY"
	case exchange.BTC:
		return "BTC"
	case exchange.ETH:
		return "ETH"
	default:
		panic(fmt.Sprintf("bitflyer err balance %s", b))
	}
}
