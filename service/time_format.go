package service

import "time"

func formatExpireTimeUnix(millis int64) string {
	return time.Unix(millis/1000, 0).Format("2006-01-02 15:04")
}
