package util

import "time"

// cst 东八区，用来划"天"的边界
var cst = time.FixedZone("CST", 8*3600)

// DayToMilli 过期天数转过期时间（13位时间戳）
func DayToMilli(day uint) uint {
	return NowMilli() + day*86400*1000
}

// NowMilli 获取当前时间的13位时间戳
func NowMilli() uint {
	return uint(time.Now().UnixMilli())
}

// TodayRangeMilli 返回东八区"今天"的起止时间戳：[今天 00:00:00, 明天 00:00:00)
func TodayRangeMilli() (uint, uint) {
	now := time.Now().In(cst)
	year, month, day := now.Date()
	start := time.Date(year, month, day, 0, 0, 0, 0, cst)
	return uint(start.UnixMilli()), uint(start.Add(24 * time.Hour).UnixMilli())
}
