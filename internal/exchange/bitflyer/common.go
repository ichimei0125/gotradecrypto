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
	ID                         int64      `json:"id" gorm:"primaryKey;autoIncrement:false"`
	Side                       string     `json:"side" gorm:"not null"`
	Price                      float64    `json:"price" gorm:"not null"`
	Size                       float64    `json:"size" gorm:"not null"`
	ExecDate                   CustomTime `json:"exec_date" gorm:"not null;index"`
	BuyChildOrderAcceptanceId  string     `json:"buy_child_order_acceptance_id"`
	SellChildOrderAcceptanceId string     `json:"sell_child_order_acceptance_id"`
	IsSync                     bool       `gorm:"default:false"`
	CreatedAt                  time.Time  `gorm:"not null"`
	UpdatedAt                  time.Time  `gorm:"not null"`
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
