package bitflyer

import (
	"fmt"
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/common"
	"github.com/ichimei0125/gotradecrypto/internal/db"
	"github.com/ichimei0125/gotradecrypto/internal/exchange"
)

// masterID: 現物買、信用買、信用売 -> 現物売、信用返済買、信用返済売
func insertOrder(acceptID string, symbol exchange.Symbol, side exchange.Side, size float64) {

	record := db.OrderHistory{
		ID:      acceptID,
		Symbol:  string(symbol),
		Side:    string(side),
		Size:    size,
		Time:    common.GetNow(),
		SendCnt: 0,
	}

	_db := db.OpenDB()

	db.Insert(_db, &record)
}

func (b *Bitflyer) CheckUnfinishedOrder(symbol exchange.Symbol) {
	_db := db.OpenDB()
	oh_list := db.GetAllRecords(_db)
	if len(oh_list) <= 0 {
		return
	}

	in_time := common.GetNow().Add(time.Duration(-common.ORDER_WAIT_MINUTE) * time.Minute)

	for _, oh := range oh_list {
		if oh.Time.After(in_time) {
			continue
		}

		child_order := getChildOrderByID(getProductCode(symbol), oh.ID)
		if len(child_order) > 1 {
			msg := ""
			for _, o := range child_order {
				msg += o.ChildOrderID + ", "
			}
			panic(fmt.Sprintf("strange child order get by id %s", msg))
		}

		if child_order[0].ChildOrderState == getOrderStatus(exchange.COMPLETED) {
			db.Delete(_db, &oh)
			continue
		}

		b.SellCypto(symbol, child_order[0].Size, -1.0)
	}

}
