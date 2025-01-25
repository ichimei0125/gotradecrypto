package db

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/common"
	"github.com/ichimei0125/gotradecrypto/internal/exchange"
)

func getDBTradeTableName(exchangeName, symbol string) string {
	return strings.ToLower(common.GetUniqueName(exchangeName, symbol) + "_trade")
}

func BulkInsertDBTrade(trades []exchange.Trade, exchangeName, symbol string) {
	tableName := getDBTradeTableName(exchangeName, symbol)

	if !_db.Migrator().HasTable(tableName) {
		err := _db.Table(tableName).AutoMigrate(trades)
		if err != nil {
			log.Fatalf("Failed to create table for model %T: %v", trades, err)
		}
	}

	if len(trades) <= 0 {
		return
	}

	batch_size := 1000

	for i := 0; i < len(trades); i += batch_size {
		end_index := i + batch_size
		if i+batch_size > len(trades) {
			end_index = len(trades)
		}
		insert_trades := trades[i:end_index]

		// gene sql
		sql := "INSERT IGNORE INTO " + tableName + " (id, side, price, size, execution_time) VALUES "
		vals := []interface{}{}

		for _, trade := range insert_trades {
			sql += "(?, ?, ?, ?, ?),"
			vals = append(vals, trade.ID, trade.Side, trade.Price, trade.Size, trade.ExecutionTime)
		}

		// rm last comma
		sql = sql[:len(sql)-1]

		result := _db.Exec(sql, vals...)
		if result.Error != nil {
			log.Fatalf("Failed to insert trades with IGNORE: %v", result.Error)
		}
	}

}

func GetDBTradeAfter(t time.Time, exchangeName, symbol string) ([]exchange.Trade, error) {
	var res []exchange.Trade
	tableName := getDBTradeTableName(exchangeName, symbol)
	err := _db.Table(tableName).Where("execution_time >= ?", t).Order("execution_time ASC").Find(&res).Error
	return res, err
}

// return: new, old
func GetDBTradeMaxMinExecTime(exchangeName, symbol string) (time.Time, time.Time, error) {
	tableName := getDBTradeTableName(exchangeName, symbol)
	var res struct {
		MaxTime time.Time `gorm:"column:max_time"`
		MinTime time.Time `gorm:"column:min_time"`
	}

	err := _db.Raw(
		fmt.Sprintf("SELECT MAX(execution_time) AS max_time, MIN(execution_time) AS min_time FROM %s", tableName)).
		Scan(&res).Error

	return res.MaxTime, res.MinTime, err
}
