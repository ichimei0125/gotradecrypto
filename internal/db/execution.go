package db

import (
	"fmt"
	"log"
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/common"
	"gorm.io/gorm"
)

func getDBExecutionTableName(exchangeName, symbol string) string {
	return common.GetUniqueName(exchangeName, symbol) + "_execution"
}

func BulkInsertDBExecution(executions interface{}, exchangeName, symbol string) {
	tableName := getDBExecutionTableName(exchangeName, symbol)

	if !_db.Migrator().HasTable(executions) {
		err := _db.Table(tableName).AutoMigrate(executions)
		if err != nil {
			log.Fatalf("Failed to create table for model %T: %v", executions, err)
		}
	}

	result := _db.Table(tableName).Create(executions)
	if result.Error != nil {
		panic(fmt.Sprintf("db cannot insert execution: %s", result.Error.Error()))
	}
}

func GetDBExecutionAfter(t time.Time, exchangeName, symbol string) *gorm.DB {
	tableName := getDBExecutionTableName(exchangeName, symbol)
	res := _db.Table(tableName).Where("exec_date >= ?", t).Order("exec_date DESC")
	return res
}
