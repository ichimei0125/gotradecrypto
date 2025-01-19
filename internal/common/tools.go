package common

import "time"

func GetNow() time.Time {
	return time.Now()
}

func GetUTCNow() time.Time {
	return time.Now().UTC()
}

func GetUniqueName(exchangeName string, symbol string) string {
	return exchangeName + "_" + symbol
}

func LastDayToDate(last_days int) time.Time {
	return GetNow().Add(time.Duration(-last_days) * 24 * time.Hour)
}
