package bitflyer

import (
	"bytes"
	"time"
)

const baseURL = "https://api.bitflyer.com"

type Bitflyer struct{}

// Execution represents a single execution from bitFlyer
type Execution struct {
	Id       int64      `json:"id"`
	ExecDate CustomTime `json:"exec_date"`
	Price    float64    `json:"price"`
	Size     float64    `json:"size"`
}

type CustomTime struct {
	time.Time
}

var TimeLayout = "2006-01-02T15:04:05.000"

// 尝试解析的时间格式列表
var timeFormats = []string{
	"2006-01-02T15:04:05.000", // 原始格式
	"2006-01-02T15:04:05.00",  // 两位毫秒数
	"2006-01-02T15:04:05.0",   // 一位毫秒数
	"2006-01-02T15:04:05",     // 无毫秒数
}

// UnmarshalJSON 用于自定义时间格式的解析
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	// 去除时间字符串两端的引号
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
	return parseErr // 返回最后一个解析错误
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
