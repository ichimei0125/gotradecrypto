package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var _db *gorm.DB

type OrderHistory struct {
	ID string `gorm:"primaryKey"` // 订单id
	// RefId   string    // 原始的订单id， 现物卖->现物买， 信用返还买->信用返还卖，信用返还卖->信用返还买
	Symbol  string    `gorm:"not null"`
	Side    string    `gorm:"not null"`
	Size    float64   `gorm:"not null"`
	Time    time.Time `gorm:"not null"`
	SendCnt int       `gorm:"not null"` // 发出订单的次数，原始为0， 超过次数后改为 成行
}

func InitDB() *gorm.DB {
	var err error
	err = os.MkdirAll("data", os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("db cannot create folder: %s", err.Error()))
	}

	_db, err = gorm.Open(sqlite.Open("data/local.db"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("db cannot open: %s", err.Error()))
	}

	// 配置连接池
	sqlDB, err := _db.DB() // 获取底层 *sql.DB 对象
	if err != nil {
		log.Fatalf("Failed to get underlying DB: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)                  // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)                 // 最大打开连接数
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)  // 空闲连接最大存活时间
	sqlDB.SetConnMaxLifetime(60 * time.Minute) // 连接最大生命周期

	// init table
	err = _db.AutoMigrate(new(OrderHistory))
	if err != nil {
		panic(fmt.Sprintf("db cannot auto migrate: %s", err.Error()))
	}
	return _db
}

func CloseDB() {
	sqlDB, err := _db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying DB for close: %v", err)
	}
	sqlDB.Close()
}

func Insert(oh *OrderHistory) {
	result := _db.Create(oh)
	if result.Error != nil {
		panic(fmt.Sprintf("db cannot insert: %s, id: %s", result.Error.Error(), oh.ID))
	}
}

func GetAllRecords() []OrderHistory {
	var records []OrderHistory
	result := _db.Find(&records)
	if result.Error != nil {
		panic(fmt.Sprintf("db cannot read all: %s", result.Error.Error()))
	}
	return records
}

func Update(oh *OrderHistory) {
	result := _db.Model(new(OrderHistory)).Where("id = ?", oh.ID).Updates(oh)
	if result.Error != nil {
		panic(fmt.Sprintf("db cannot update: %s, id: %s", result.Error.Error(), oh.ID))
	}
}

func Delete(oh *OrderHistory) {
	result := _db.Delete(oh)
	if result.Error != nil {
		panic(fmt.Sprintf("db cannot delete: %s, id: %s", result.Error.Error(), oh.ID))
	}
}

func DeleteByID(id string) {
	// result := db.Delete(&OrderHistory{}, id)
	result := _db.Where("id = ?", id).Delete(new(OrderHistory))

	if result.Error != nil {
		panic(fmt.Sprintf("db cannot delete: %s, id: %s", result.Error.Error(), id))
	}
}
