package bitflyer

import (
	"fmt"
	"time"

	"github.com/ichimei0125/gotradecrypto/internal/common"
	"github.com/ichimei0125/gotradecrypto/internal/db"
	"github.com/ichimei0125/gotradecrypto/internal/exchange"
)

func insertOrder(acceptID string, symbol exchange.Symbol, side exchange.Side, size float64) {

	record := db.OrderHistory{
		ID:       acceptID,
		Exchange: new(Bitflyer).Name(),
		Symbol:   string(symbol),
		Side:     string(side),
		Size:     size,
		SendCnt:  0,
	}

	db.Insert(&record)
}

func (b *Bitflyer) CheckUnfinishedOrder(symbol exchange.Symbol) {
	oh_list := db.GetAllRecords()
	if len(oh_list) <= 0 {
		return
	}

	in_time := common.GetNow().Add(time.Duration(-common.ORDER_WAIT_MINUTE) * time.Minute)

	for _, oh := range oh_list {
		if oh.CreatedAt.Before(in_time) {
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

		if len(child_order) <= 0 || child_order[0].ChildOrderState == getOrderStatus(exchange.COMPLETED) {
			db.Delete(&oh)
			continue
		}

		b.SellCypto(symbol, child_order[0].Size, -1.0)
	}

}
