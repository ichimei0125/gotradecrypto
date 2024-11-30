package db_test

import (
	"testing"
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/db"
)

func TestOpenDB(t *testing.T) {
	db.OpenDB()
}

func TestInsert(t *testing.T) {
	_db := db.OpenDB()
	data := db.OrderHistory{
		ID:      "testid2",
		Symbol:  "testsymbol",
		Side:    "testside",
		Size:    100.3,
		Time:    time.Now(),
		SendCnt: 0,
	}
	db.Insert(_db, &data)
}

func TestUpdate(t *testing.T) {
	_db := db.OpenDB()
	data := db.OrderHistory{
		ID:      "testid2",
		Symbol:  "testsymbol",
		Side:    "testside",
		Size:    100.3,
		Time:    time.Now(),
		SendCnt: 1,
	}
	db.Update(_db, &data)
}

func TestDelete(t *testing.T) {
	_db := db.OpenDB()
	data := db.OrderHistory{
		ID:      "testid2",
		Symbol:  "testsymbol",
		Side:    "testside",
		Size:    100.3,
		Time:    time.Now(),
		SendCnt: 1,
	}
	db.Delete(_db, &data)
}

func TestGetAll(t *testing.T) {
	_db := db.OpenDB()
	oh_list := db.GetAllRecords(_db)
	for _, oh := range oh_list {
		t.Log(oh.ID)
	}
}

func TestDeleteByID(t *testing.T) {
	_db := db.OpenDB()
	db.DeleteByID(_db, "testid2")
}
