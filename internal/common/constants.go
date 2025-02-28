package common

import "time"

// TODO: remove this
const (
	REFRESH_INTERVAL   int = 1    // データ更新間隔
	KLINE_INTERVAL     int = 3    // Klineの時間単位
	KLINE_LENGTH       int = 1000 // Klineの長さ（cache）
	TRADE_WATCH_MINUTE int = 10   // 10分以内買い/売りがあれば、何もしない
	ORDER_WAIT_MINUTE  int = 3    // 規定時間内注文
)

var (
	NULLDATE time.Time = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
)

func IsNullDate(t time.Time) bool {
	return t.Equal(NULLDATE)
}

var (
	ENV_CONFIG_PATH [2]string = [2]string{"GOTRADECRYPTO_CONFIG_PATH", "config.yaml"}
	ENV_LOG_PATH    [2]string = [2]string{"GOTRADECRYPTO_LOG_PATH", "log"}
)
