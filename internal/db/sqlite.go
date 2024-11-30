package db

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type OrderHistory struct {
	Id string `gorm:"primaryKey"` // 订单id
	// RefId   string    // 原始的订单id， 现物卖->现物买， 信用返还买->信用返还卖，信用返还卖->信用返还买
	Symbol  string    `gorm:"not null"`
	Side    string    `gorm:"not null"`
	Size    float64   `gorm:"not null"`
	Time    time.Time `gorm:"not null"`
	SendCnt int       `gorm:"not null"` // 发出订单的次数，原始为0， 超过次数后改为 成行
}

func OpenDB() *gorm.DB {
	err := os.MkdirAll("data", os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("db cannot create folder: %s", err.Error()))
	}

	db, err := gorm.Open(sqlite.Open("data/local.db"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("db cannot open: %s", err.Error()))
	}

	err = db.AutoMigrate(new(OrderHistory))
	if err != nil {
		panic(fmt.Sprintf("db cannot auto migrate: %s", err.Error()))
	}
	return db
}

func Insert(db *gorm.DB, oh *OrderHistory) {
	result := db.Create(oh)
	if result.Error != nil {
		panic(fmt.Sprintf("db cannot insert: %s, id: %s", result.Error.Error(), oh.Id))
	}
}

func GetAllRecords(db *gorm.DB) []OrderHistory {
	var records []OrderHistory
	result := db.Find(&records)
	if result.Error != nil {
		panic(fmt.Sprintf("db cannot read all: %s", result.Error.Error()))
	}
	return records
}

func Update(db *gorm.DB, oh *OrderHistory) {
	result := db.Model(new(OrderHistory)).Where("id = ?", oh.Id).Updates(oh)
	if result.Error != nil {
		panic(fmt.Sprintf("db cannot update: %s, id: %s", result.Error.Error(), oh.Id))
	}
}

func Delete(db *gorm.DB, oh *OrderHistory) {
	result := db.Delete(oh)
	if result.Error != nil {
		panic(fmt.Sprintf("db cannot delete: %s, id: %s", result.Error.Error(), oh.Id))
	}
}
