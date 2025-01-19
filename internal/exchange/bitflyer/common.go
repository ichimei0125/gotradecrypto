package bitflyer

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"time"
)

const baseURL = "https://api.bitflyer.com"

type Bitflyer struct{}

func (b *Bitflyer) Name() string {
	return "bitflyer"
}

type Execution struct {
	ID                         int64      `json:"id"`
	Side                       string     `json:"side"`
	Price                      float64    `json:"price"`
	Size                       float64    `json:"size"`
	ExecDate                   CustomTime `json:"exec_date"`
	BuyChildOrderAcceptanceId  string     `json:"buy_child_order_acceptance_id"`
	SellChildOrderAcceptanceId string     `json:"sell_child_order_acceptance_id"`
}

type CustomTime struct {
	time.Time
}

var TimeLayout = "2006-01-02T15:04:05.000"

var timeFormats = []string{
	"2006-01-02T15:04:05.000",
	"2006-01-02T15:04:05.00",
	"2006-01-02T15:04:05.0",
	"2006-01-02T15:04:05",
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := bytes.Trim(b, "\"")
	if len(s) == 0 {
		return nil
	}
	var parseErr error
	for _, layout := range timeFormats {
		t, err := time.Parse(layout, string(s))
		if err == nil {
			ct.Time = t
			return nil
		}
		parseErr = err
	}
	return parseErr
}

// for gorm Value interface
func (ct CustomTime) Value() (driver.Value, error) {
	return ct.Time.Format(TimeLayout), nil
}
func (ct *CustomTime) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		ct.Time = v
		return nil
	case string:
		t, err := time.Parse(TimeLayout, v)
		if err != nil {
			return fmt.Errorf("invalid time string: %v", err)
		}
		ct.Time = t
		return nil
	case []byte:
		t, err := time.Parse(TimeLayout, string(v))
		if err != nil {
			return fmt.Errorf("invalid time byte array: %v", err)
		}
		ct.Time = t
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
}

func CustomTimeToString(t CustomTime) string {
	str := t.Time.Format(TimeLayout)
	return str
}

func StringToCustomTime(s string) CustomTime {
	t, _ := time.Parse(TimeLayout, s)
	return CustomTime{
		Time: t,
	}
}
