package db_test

import (
	"testing"
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/db"
	"github.com/ichimei0125/gotradecrypto/internal/exchange"
	"github.com/ichimei0125/gotradecrypto/internal/exchange/bitflyer"
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

func TestDBExecution(t *testing.T) {
	db.InitDB()
	defer db.CloseDB()
	d := bitflyer.CustomTime{
		Time: time.Now(),
	}

	exe := []bitflyer.Execution{{
		ID:       425436,
		Side:     "TestSide",
		Size:     543.5,
		Price:    101.1,
		ExecDate: d,
	}}

	db.BulkInsertDBExecution(exe, new(bitflyer.Bitflyer).Name(), string(exchange.XRPJPY))
}
