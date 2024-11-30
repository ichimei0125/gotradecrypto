package db

import (
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type OrderHistory struct {
	Id      string    `gorm:"primaryKey"` // 订单id
	RefId   string    // 原始的订单id， 现物卖->现物买， 信用返还买->信用返还卖，信用返还卖->信用返还买
	Symbol  string    `gorm:"not null"`
	Side    string    `gorm:"not null"`
	Size    float64   `gorm:"not null"`
	Time    time.Time `gorm:"not null"`
	SendCnt int       `gorm:"not null"` // 发出订单的次数，原始为0， 超过次数后改为 成行
}

func OpenDB() (*gorm.DB, error) {
	err := os.MkdirAll("data", os.ModePerm)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open("data/local.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(new(OrderHistory))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Insert(db *gorm.DB, oh *OrderHistory) error {
	result := db.Create(oh)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func Update(db *gorm.DB, oh *OrderHistory) error {
	result := db.Model(new(OrderHistory)).Where("id = ?", oh.Id).Updates(oh)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func Delete(db *gorm.DB, oh *OrderHistory) error {
	result := db.Delete(oh)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
