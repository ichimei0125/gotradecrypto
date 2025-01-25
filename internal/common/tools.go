package common

import (
	"time"
)

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

// save memory by expand slice's cap first
func AppendBig[T any](arr1, arr2 []T) []T {
	if cap(arr1) < len(arr1)+len(arr2) {
		newArr := make([]T, len(arr1), len(arr1)+len(arr2))
		copy(newArr, arr1)
		arr1 = newArr
	}
	return append(arr1, arr2...)
}
