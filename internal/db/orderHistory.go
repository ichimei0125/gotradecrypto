package db

import (
	"fmt"
	"time"
)

type OrderHistory struct {
	ID          string    `gorm:"primaryKey"` // 订单id
	Exchange    string    `gorm:"not null"`
	Symbol      string    `gorm:"not null"`
	Side        string    `gorm:"not null"`
	Size        float64   `gorm:"not null"`
	SendCnt     int       `gorm:"not null"` // 发出订单的次数，原始为0， 超过次数后改为 成行
	IsCompleted bool      `gorm:"not null"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
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
