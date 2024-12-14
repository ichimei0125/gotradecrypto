package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/common"
	"github.com/ichimei0125/gotradecrypto/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var _db *gorm.DB

func InitDB() *gorm.DB {
	var err error
	err = os.MkdirAll("data", os.ModePerm)
	if err != nil {
		panic(fmt.Sprintf("db cannot create folder: %s", err.Error()))
	}

	db_path := config.GetEnvVar(common.ENV_DB_PATH[0], common.ENV_DB_PATH[1])
	_db, err = gorm.Open(sqlite.Open(db_path), &gorm.Config{
		CreateBatchSize: 1000,
	})
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
