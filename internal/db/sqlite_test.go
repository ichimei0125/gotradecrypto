package db_test

import (
	"testing"
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/db"
)

func TestOpenDB(t *testing.T) {
	_, err := db.OpenDB()
	if err != nil {
		t.Error(err)
	}
}

func TestInsert(t *testing.T) {
	_db, _ := db.OpenDB()
	data := db.OrderHistory{
		Id:      "testid2",
		RefId:   "testrefid",
		Symbol:  "testsymbol",
		Side:    "testside",
		Size:    100.3,
		Time:    time.Now(),
		SendCnt: 0,
	}
	err := db.Insert(_db, &data)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdate(t *testing.T) {
	_db, _ := db.OpenDB()
	data := db.OrderHistory{
		Id:      "testid2",
		RefId:   "testrefid",
		Symbol:  "testsymbol",
		Side:    "testside",
		Size:    100.3,
		Time:    time.Now(),
		SendCnt: 1,
	}
	err := db.Update(_db, &data)
	if err != nil {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	_db, _ := db.OpenDB()
	data := db.OrderHistory{
		Id:      "testid2",
		RefId:   "testrefid",
		Symbol:  "testsymbol",
		Side:    "testside",
		Size:    100.3,
		Time:    time.Now(),
		SendCnt: 1,
	}
	err := db.Delete(_db, &data)
	if err != nil {
		t.Error(err)
	}
}
