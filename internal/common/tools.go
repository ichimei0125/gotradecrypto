package common

import "time"

func GetNow() time.Time {
	return time.Now()
}

func GetUTCNow() time.Time {
	return time.Now().UTC()
}
