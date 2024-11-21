package bitflyer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/common"
	"github.com/ichimei0125/gotradecrypto/internal/exchange"
)

func (b *Bitflyer) FetchKLine(s exchange.Symbol, cache *[]exchange.KLine) {
	var oldest_id int64 = 0

	klineMap := make(map[time.Time]exchange.KLine)
	for _, kline := range *cache {
		klineMap[kline.OpenTime] = kline
	}

	for {
		excutions := FetchExecution(s, 0, oldest_id, 0)
		oldest_id = excutions[len(excutions)-1].Id

		new_kline := convertExetutionsToKLine(excutions)

		for _, kline := range new_kline {
			klineMap[kline.OpenTime] = kline
		}

		// 将 map 转换为切片
		merged := make([]exchange.KLine, 0, len(klineMap))
		for _, kline := range klineMap {
			merged = append(merged, kline)
		}

		if len(merged) >= common.KLINE_LENGTH {
			// 按时间倒序排序
			sort.Slice(merged, func(i, j int) bool {
				return merged[i].OpenTime.After(merged[j].OpenTime)
			})

			*cache = merged[:common.KLINE_LENGTH:common.KLINE_LENGTH]

			break
		}
	}
}

func bitflyerPublicAPICore(url string) *http.Response {
	// TODO impl 制限
	resp, err := http.Get(url)
	if err != nil {
		log.Panicln("cannot get bitflyer public api, maybe limited")
	}
	return resp
}

func FetchExecution(s exchange.Symbol, count int, before_id int64, after_id int64) []Execution {

	symbol := getsymbol(s)

	endpoint := "/v1/getexecutions"
	u, err := url.Parse(baseURL + endpoint)
	if err != nil {
		panic("error when parse bitflyer url")
	}
	if count <= 0 {
		count = 500
	}
	// Add query parameters
	q := u.Query()
	q.Set("product_code", symbol)
	q.Set("count", fmt.Sprintf("%d", count))
	if before_id > 0 {
		q.Set("before", fmt.Sprintf("%d", before_id))
	}
	if after_id > 0 {
		q.Set("after", fmt.Sprintf("%d", after_id))
	}
	u.RawQuery = q.Encode()

	resp := bitflyerPublicAPICore(u.String())
	defer resp.Body.Close()

	// // Decode the JSON response into a slice of Executions
	var executions []Execution
	if err := json.NewDecoder(resp.Body).Decode(&executions); err != nil {
		panic("wrong bitflyer public response json")
	}

	return executions
}

type byExecDate []Execution

func (a byExecDate) Len() int           { return len(a) }
func (a byExecDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byExecDate) Less(i, j int) bool { return a[i].ExecDate.Time.Before(a[j].ExecDate.Time) }

func convertExetutionsToKLine(executions []Execution) []exchange.KLine {
	sort.Sort(sort.Reverse(byExecDate(executions)))

	minute_unit := common.KLINE_INTERVAL

	// 将最新的时间依据时间单位，选择最近的区间
	// 比如5分钟的时间单位: 16:00, 16:05, 16:10...
	// 16:00~16:05的数据中，算16:00的始值，终值，高值，低值
	lastest_time := executions[0].ExecDate
	var norm_minute int = (lastest_time.Minute()/minute_unit + 1) * minute_unit
	time_unit := time.Date(lastest_time.Year(), lastest_time.Month(), lastest_time.Day(), lastest_time.Hour(), norm_minute, 0, 0, lastest_time.Location())

	kline := []exchange.KLine{}

	time_unit = time_unit.Add(-time.Minute * time.Duration(minute_unit))

	tmp_kline := exchange.KLine{
		Open:     executions[0].Price,
		Close:    executions[0].Price,
		High:     executions[0].Price,
		Low:      executions[0].Price,
		OpenTime: time_unit,
	}

	for _, execution := range executions[1:] {
		if execution.ExecDate.Time.Before(time_unit) {
			kline = append(kline, tmp_kline)
			time_unit = time_unit.Add(-time.Minute * time.Duration(minute_unit))
			tmp_kline = exchange.KLine{
				Open:     execution.Price,
				Close:    execution.Price,
				High:     execution.Price,
				Low:      execution.Price,
				OpenTime: time_unit,
			}
		} else {
			// 時間単位内
			tmp_kline.Open = execution.Price
			tmp_kline.High = max(tmp_kline.High, execution.Price)
			tmp_kline.Low = min(tmp_kline.Low, execution.Price)
		}
	}
	return kline
}

func getsymbol(symbol exchange.Symbol) string {
	switch symbol {
	case exchange.BTCJPY:
		return "BTC_JPY"
	case exchange.XRPJPY:
		return "XRP_JPY"
	case exchange.ETHJPY:
		return "ETH_JPY"
	case exchange.XLMJPY:
		return "XLM_JPY"
	case exchange.MONAJPY:
		return "MONA_JPY"
	case exchange.ETHBTC:
		return "ETH_BTC"
	case exchange.BCHBTC:
		return "BCH_BTC"
	case exchange.FX_BTCJPY:
		return "FX_BTC_JPY"
	default:
		panic(fmt.Sprintf("bitflyer err symbol %s", symbol))
	}
}
