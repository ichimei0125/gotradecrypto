package common

import "time"

const (
	REFRESH_INTERVAL   int = 1    // データ更新間隔
	KLINE_INTERVAL     int = 3    // Klineの時間単位
	KLINE_LENGTH       int = 1000 // Klineの長さ（cache）
	TRADE_WATCH_MINUTE int = 10   // 10分以内買い/売りがあれば、何もしない
	ORDER_WAIT_MINUTE  int = 3    // 規定時間内注文
)

var (
	ZERO_TIME_LOCAL time.Time = time.Date(1, 1, 1, 0, 0, 0, 0, time.Local)
)
