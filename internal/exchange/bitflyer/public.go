package bitflyer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/common"
	"github.com/ichimei0125/gotradecrypto/internal/db"
	"github.com/ichimei0125/gotradecrypto/internal/exchange"
	"github.com/ichimei0125/gotradecrypto/internal/logger"
)

var (
	cache_execution sync.Map
	db_lastest_time sync.Map
)

var (
	cache_trades sync.Map
)

func (b *Bitflyer) FetchCandleSticks(s string, cache *[]exchange.CandleStick) {
	// 	var oldest_id int64 = 0
	// 	var lastest_id int64 = 0
	// 	uniqueName := common.GetUniqueName(new(Bitflyer).Name(), string(s))

	// 	klineMap := make(map[time.Time]exchange.CandleStick)
	// 	for _, kline := range *cache {
	// 		klineMap[kline.OpenTime] = kline
	// 	}

	// 	// load cache
	// 	start_time := common.GetNow().Add(-time.Duration(common.KLINE_INTERVAL*(common.KLINE_LENGTH+5)) * time.Minute)
	// 	_, ok := cache_execution.Load(uniqueName)
	// 	if !ok && len(*cache) <= 0 {
	// 		db_executions := []Execution{}
	// 		db.GetDBTradeAfter(start_time, new(Bitflyer).Name(), string(s)).Find(&db_executions)
	// 		if len(db_executions) > 0 {
	// 			cache_execution.Store(uniqueName, db_executions)
	// 			lastest_id = db_executions[0].ID
	// 			db_lastest_time.Store(uniqueName, db_executions[0].ExecDate.Time)
	// 		}
	// 	}

	// 	is_update := false
	// 	// get data to cache
	// 	for {
	// 		executions := FetchExecution(s, 0, oldest_id, lastest_id)
	// 		if len(executions) <= 0 {
	// 			is_update = true
	// 			break
	// 		}

	// 		// cache
	// 		cached, ok := cache_execution.Load(uniqueName)
	// 		if ok {
	// 			cached = append(executions, cached.([]Execution)...)
	// 		} else {
	// 			cached = executions
	// 		}
	// 		cache_execution.Store(uniqueName, cached)

	// 		oldest_id = executions[len(executions)-1].ID
	// 		lastest_id = 0

	// 		// empty range between db data?
	// 		_db_lastest_time, ok := db_lastest_time.Load(uniqueName)
	// 		if ok {
	// 			if _db_lastest_time.(time.Time).After(executions[len(executions)-1].ExecDate.Time) {
	// 				_cached := cached.([]Execution)
	// 				// rm repeat
	// 				unique_cache := make(map[time.Time]Execution)
	// 				for _, __cached := range _cached {
	// 					unique_cache[__cached.ExecDate.Time] = __cached
	// 				}
	// 				// map -> slice
	// 				_cached = make([]Execution, 0, len(_cached))
	// 				for _, uc := range unique_cache {
	// 					_cached = append(_cached, uc)
	// 				}
	// 				// sort by: new -> old
	// 				sort.Slice(_cached, func(i, j int) bool {
	// 					return _cached[i].ExecDate.Time.After(_cached[j].ExecDate.Time)
	// 				})
	// 				cache_execution.Store(uniqueName, cached)

	// 				break
	// 			}
	// 		}
	// 	}

	// 	if !is_update && len(*cache) > 0 {
	// 		return
	// 	}

	// 	// generate & update kline
	// 	cached, _ := cache_execution.Load(uniqueName)
	// 	new_kline := convertExetutionsToKLine(cached.([]Execution), common.KLINE_INTERVAL)
	// 	// rm repeat kline, old kline's indicators will keep, only append new
	// 	for _, kline := range new_kline {
	// 		if _, exist := klineMap[kline.OpenTime]; !exist {
	// 			klineMap[kline.OpenTime] = kline
	// 		}
	// 	}
	// 	merged := make([]exchange.CandleStick, 0, len(klineMap))
	// 	for _, kline := range klineMap {
	// 		merged = append(merged, kline)
	// 	}
	// 	// order: new -> old
	// 	sort.Slice(merged, func(i, j int) bool {
	// 		return merged[i].OpenTime.After(merged[j].OpenTime)
	// 	})
	// 	*cache = merged[:common.KLINE_LENGTH:common.KLINE_LENGTH]

	// 	// cache
	// 	_executions, loaded := cache_execution.Load(uniqueName)
	// 	if !loaded {
	// 		return
	// 	}
	// 	executions := _executions.([]Execution)
	// 	__db_lastest_time, ok_db_lastest_time := db_lastest_time.Load(uniqueName)
	// 	_db_lastest_time := __db_lastest_time.(time.Time)

	// insert := make([]Execution, 0, len(executions))
	// new_cache := make([]Execution, 0, len(executions))
	//
	//	if !ok_db_lastest_time {
	//		insert = executions
	//	} else {
	//
	//		for _, e := range executions {
	//			if e.ExecDate.Time.After(_db_lastest_time) {
	//				insert = append(insert, e)
	//			}
	//			if e.ExecDate.Time.After(_db_lastest_time.Add(time.Duration(-common.KLINE_INTERVAL*3) * time.Minute)) {
	//				new_cache = append(new_cache, e)
	//			}
	//		}
	//	}
	//
	// // db insert
	//
	//	if len(insert) > 0 {
	//		db.BulkInsertDBTrade(insert, new(Bitflyer).Name(), string(s))
	//		db_lastest_time.Store(uniqueName, insert[0].ExecDate.Time)
	//	}
	//
	// // update cache
	// cache_execution.Store(uniqueName, new_cache)
}

func FetchTrades(since time.Time, symbol string) []exchange.Trade {
	var trades []exchange.Trade
	_trades, ok := cache_trades.Load(symbol)
	if ok {
		// load from cache
		trades = _trades.([]exchange.Trade)
	} else {
		// load from db if not cached
		var err error
		trades, err = db.GetDBTradeAfter(since, new(Bitflyer).GetInfo().Name, string(symbol))
		// get max history data for first time
		if err != nil {
			since = common.GetNow().Add(-30 * 24 * time.Hour).UTC()
		}
	}
	if len(trades) > 0 {
		// update since if cached or loaded from db
		since = trades[0].ExecutionTime
	}

	executions := []Execution{}
	var before_id int64
	for {
		_executions := FetchExecution(symbol, -1, before_id, -1)
		before_id = _executions[len(_executions)-1].ID
		_since := _executions[len(_executions)-1].ExecDate.Time

		executions = append(executions, _executions...)
		if _since.Before(since) {
			break
		}
	}
	// to trades
	newTrades := convertExecutionsToTrades(&executions)
	// sort trades ASC, old -> new
	slices.SortFunc(newTrades, func(a, b exchange.Trade) int {
		if a.ExecutionTime.Before(b.ExecutionTime) {
			return -1
		} else if a.ExecutionTime.After(b.ExecutionTime) {
			return 1
		} else {
			// sort by id
			a_id, _ := strconv.ParseUint(a.ID, 10, 64)
			b_id, _ := strconv.ParseUint(b.ID, 10, 64)
			if a_id < b_id {
				return -1
			}
			return 1
		}
	})
	// rm duplicated trades
	for _, newTrade := range newTrades {
		if newTrade.ExecutionTime.After(since) {
			trades = append(trades, newTrade)
		}
	}

	// update cache
	cache_trades.Store(symbol, trades)
	// insert db
	// TODO: reduce insert frequency
	db.BulkInsertDBTrade(newTrades, new(Bitflyer).GetInfo().Name, symbol)

	return trades
}

func bitflyerPublicAPICore(url string) (*http.Response, error) {
	rl := exchange.GetRateLimiter(new(Bitflyer).GetInfo().Name+"_Public", 5*time.Minute, (500 - 20))
	for {
		ok, waitTime := rl.Allow()
		if ok {
			resp, err := http.Get(url)
			if err != nil {
				panic(fmt.Sprintf("cannot get bitflyer public api, maybe limited \n %s", err.Error()))
			}
			return resp, nil
		}
		logger.Info(new(Bitflyer), "", fmt.Sprintf("over HTTP limit, wait %d, url %s", waitTime, url))
		time.Sleep(waitTime)
	}
}

func FetchExecution(symbol string, count int, before_id int64, after_id int64) []Execution {

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

	resp, err := bitflyerPublicAPICore(u.String())
	if err != nil {
		return []Execution{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("Failed to read response body: %s", err))
	}

	// Decode JSON
	var executions []Execution
	if err := json.Unmarshal(body, &executions); err != nil {
		panic(fmt.Sprintf("Wrong bitflyer public response JSON:\n%s", string(body)))
	}

	return executions
}

// type byExecDate []Execution

// func (a byExecDate) Len() int           { return len(a) }
// func (a byExecDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
// func (a byExecDate) Less(i, j int) bool { return a[i].ExecDate.Time.Before(a[j].ExecDate.Time) }

func convertExetutionsToKLine(executions []Execution, minute_unit int) []exchange.CandleStick {
	// sort.Sort(sort.Reverse(byExecDate(executions)))
	// use sorted ?

	// 将最新的时间依据时间单位，选择最近的区间
	// 比如5分钟的时间单位: 16:00, 16:05, 16:10...
	// 16:00~16:05的数据中，算16:00的始值，终值，高值，低值
	lastest_time := executions[0].ExecDate
	var norm_minute int = (lastest_time.Minute()/minute_unit + 1) * minute_unit
	time_unit := time.Date(lastest_time.Year(), lastest_time.Month(), lastest_time.Day(), lastest_time.Hour(), norm_minute, 0, 0, lastest_time.Location())

	kline := []exchange.CandleStick{}

	time_unit = time_unit.Add(-time.Minute * time.Duration(minute_unit))

	tmp_kline := exchange.CandleStick{
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
			tmp_kline = exchange.CandleStick{
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

func convertExecutionsToTrades(exections *[]Execution) []exchange.Trade {
	data := *exections
	res := []exchange.Trade{}
	for _, execution := range data {
		res = append(res, exchange.Trade{
			ID:            fmt.Sprintf("%d", execution.ID),
			Side:          execution.Side,
			Price:         execution.Price,
			Size:          execution.Size,
			ExecutionTime: execution.ExecDate.Time,
		})
	}
	return res
}
