package Utils

import (
	"time"
)

/**
判断当前时间是否已经到第二天
param： time 当前时间（time.Now().UTC()）
result：bool 是否已经到第二天
**/
func CheckIsNextDay(time time.Time) bool{
	if time.Hour() == 0 && time.Minute() == 0 && time.Second() == 0 {
		return true
	}
	return false
}

/**
判断当前时间是否已经到下一小时
param： time 当前时间（time.Now().UTC()）
result：bool 是否已经到下一小时
**/
func CheckIsNextHour(time time.Time) bool{
	if time.Minute() == 0 && time.Second() == 0 {
		return true
	}
	return false
}

/**
判断当前时间是否已经到下一分钟
param： time 当前时间（time.Now().UTC()）
result：bool 是否已经到下一分钟install
**/
func CheckIsNextMinite(time time.Time) bool{
	if time.Second() == 0 {
		return true
	}
	return false
}