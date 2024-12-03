package db_test

import (
	"testing"

	"github.com/ichimei0125/gotradecrypto/internal/db"
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
