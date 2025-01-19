package db_test

import (
	"testing"
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/db"
	"github.com/ichimei0125/gotradecrypto/internal/exchange"
)

func TestOpenDB(t *testing.T) {
	db.InitDB()
	defer db.CloseDB()
}

func TestInsert(t *testing.T) {
	db.InitDB()
	defer db.CloseDB()
	data := db.OrderHistory{
		ID:     "testid2",
		Symbol: "testsymbol",
		Side:   "testside",
		Size:   100.3,
	}
	db.Insert(&data)
}

func TestUpdate(t *testing.T) {
	db.InitDB()
	defer db.CloseDB()
	data := db.OrderHistory{
		ID: "testid2",
	}
	db.Update(&data)
}

func TestDelete(t *testing.T) {
	db.InitDB()
	defer db.CloseDB()
	data := db.OrderHistory{
		ID:      "testid2",
		Symbol:  "testsymbol",
		Side:    "testside",
		Size:    100.3,
		SendCnt: 1,
	}
	db.Delete(&data)
}

func TestGetAll(t *testing.T) {
	db.InitDB()
	defer db.CloseDB()
	oh_list := db.GetAllRecords()
	for _, oh := range oh_list {
		t.Log(oh.ID)
	}
}

func TestDeleteByID(t *testing.T) {
	db.InitDB()
	defer db.CloseDB()
	db.DeleteByID("testid2")
}

func TestInsertErr(t *testing.T) {
	db.InitDB()
	defer db.CloseDB()
	db.InsertErr("")
}

func TestDeleteErr(t *testing.T) {
	_t := time.Date(2024, 12, 3, 19, 55, 0, 0, time.Local)
	db.InitDB()
	defer db.CloseDB()
	db.DeleteErrAfter(_t)
}

func TestGetErr(t *testing.T) {
	_t := time.Date(2024, 12, 3, 19, 50, 0, 0, time.Local)
	db.InitDB()
	defer db.CloseDB()
	errs := db.GetErrAfter(_t)
	for _, err := range errs {
		t.Log(err.CreatedAt)
	}
}

func TestDBTrade(t *testing.T) {
	db.InitDB()
	defer db.CloseDB()

	_trade1 := exchange.Trade{
		ID:            "test1",
		Side:          "",
		Size:          12.4,
		Price:         100.1,
		ExecutionTime: time.Now().UTC(),
	}
	// _trade2 := db.Trade{
	// 	ID:            "test2",
	// 	Side:          "",
	// 	Size:          14.5,
	// 	Price:         98.1,
	// 	ExecutionTime: time.Now().UTC(),
	// }
	_trade3 := exchange.Trade{
		ID:            "test3",
		Side:          "",
		Size:          15.6,
		Price:         130.1,
		ExecutionTime: time.Now().UTC(),
	}
	_trades := []exchange.Trade{_trade1, _trade3}

	db.BulkInsertDBTrade(_trades, "testExchange", "testSymbol")
}

func TestGetDBTrade(t *testing.T) {
	db.InitDB()
	defer db.CloseDB()

	_t := time.Date(2024, 12, 5, 0, 0, 0, 0, time.UTC)

	// db_trades, err := db.GetDBTradeAfter(_t, "no_existed_exchange", "no_existed_symbol")
	db_trades, err := db.GetDBTradeAfter(_t, "bitflyer", "btc_jpy")
	if err != nil {
		t.Error(err)
	}
	t.Log(db_trades[len(db_trades)-1].ExecutionTime)
}

func TestGetDBTradeFitstLastTime(t *testing.T) {
	db.InitDB()
	defer db.CloseDB()

	new, old, _ := db.GetDBTradeMaxMinExecTime("bitflyer", "btc_jpy")
	t.Log(new, old)
}
