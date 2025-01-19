package db

import (
	"fmt"
	"log"
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _db *gorm.DB

func InitDB() *gorm.DB {
	var err error
	// err = os.MkdirAll("data", os.ModePerm)
	// if err != nil {
	// 	panic(fmt.Sprintf("db cannot create folder: %s", err.Error()))
	// }

	connectionString := config.GetConfig().ConnectionString
	_db, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{
		CreateBatchSize: 1000,
	})
	if err != nil {
		panic(fmt.Sprintf("db cannot open: %s", err.Error()))
	}

	sqlDB, err := _db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying DB: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)

	// init table
	err = _db.AutoMigrate(new(OrderHistory), new(AppErr))
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
