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
	cache_trades sync.Map
)

func (b *Bitflyer) FetchCandleSticks(since time.Time, symbol string, interval time.Duration) []exchange.CandleStick {
	trades := b.FetchTrades(since, symbol)
	// TODO: support seconds, hours, days...
	return exchange.TradesToCandleStickByMinute(trades, int(interval.Minutes()))
}

func (b *Bitflyer) FetchTrades(since time.Time, symbol string) []exchange.Trade {
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
			since = common.GetUTCNow().Add(-30 * 24 * time.Hour)
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
	_newTrades := []exchange.Trade{}
	for i := len(newTrades) - 1; i >= 0; i-- {
		if newTrades[i].ExecutionTime.Before(trades[0].ExecutionTime) {
			_newTrades = append(_newTrades, newTrades[i])
		}
	}
	trades = append(_newTrades, trades...)

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
