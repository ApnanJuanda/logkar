package utils

import "time"

func GetLocaltime() time.Time {
	location, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(location)
	return now
}
