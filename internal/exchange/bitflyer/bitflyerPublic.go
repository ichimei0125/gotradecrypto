package bitflyer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ichimei0125/gotradecrypto/internal/exchange"
)

type Bitflyer struct{}

// FetchData retrieves market data from bitFlyer
// func (b *BitflyerPublic) FetchData(productCode string, count int, before_id int64, after_id int64) []Execution {
func (b *Bitflyer) FetchPrice(s exchange.Symbol) (float64, error) {
	symbol := getsymbol(s)

	baseURL := "https://api.bitflyer.com"
	endpoint := "/v1/getexecutions"
	u, err := url.Parse(baseURL + endpoint)
	if err != nil {
		panic("error when parse bitflyer url")
	}

	// Add query parameters
	q := u.Query()
	q.Set("product_code", symbol)
	q.Set("count", fmt.Sprintf("%d", 500))
	// if before_id > 0 {
	// 	q.Set("before", fmt.Sprintf("%d", before_id))
	// }
	// if after_id > 0 {
	// 	q.Set("after", fmt.Sprintf("%d", after_id))
	// }
	u.RawQuery = q.Encode()

	// Make the request
	resp, err := http.Get(u.String())
	if err != nil {
		panic("cannot get bitflyer public api, maybe limited")
	}
	defer resp.Body.Close()

	// // Decode the JSON response into a slice of Executions
	var executions []Execution
	if err := json.NewDecoder(resp.Body).Decode(&executions); err != nil {
		panic("wrong bitflyer public response json")
	}

	return executions[0].Price, nil
}

func getsymbol(symbol exchange.Symbol) string {
	// TODO
	return "BTC_JPY"
}
