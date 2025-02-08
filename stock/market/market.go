package market

import (
	"stock-god-scraper/holiday"
	"time"
)

var stockMarketZone = time.FixedZone("Stock_Market +0800", 8*60*60)

// 股市上午开盘时间
func getStockMarketMorningOpenTime(currentTime time.Time) time.Time {
	tempTime := currentTime
	year, month, day := tempTime.In(stockMarketZone).Date()
	return time.Date(year, month, day, 9, 30, 0, 0, stockMarketZone)
}

// 股市上午收盘时间
func getStockMarketMorningCloseTime(currentTime time.Time) time.Time {
	tempTime := currentTime
	year, month, day := tempTime.In(stockMarketZone).Date()
	return time.Date(year, month, day, 11, 30, 0, 0, stockMarketZone)
}

// 股市下午开盘时间
func getStockMarketAfternoonOpenTime(currentTime time.Time) time.Time {
	tempTime := currentTime
	year, month, day := tempTime.In(stockMarketZone).Date()
	return time.Date(year, month, day, 13, 0, 0, 0, stockMarketZone)
}

// 股市下午收盘时间
func getStockMarketAfternoonCloseTime(currentTime time.Time) time.Time {
	tempTime := currentTime
	year, month, day := tempTime.In(stockMarketZone).Date()
	return time.Date(year, month, day, 15, 0, 0, 0, stockMarketZone)
}

func isStockMarketOpenTime(currentTime time.Time) bool {
	morningOpenTime := getStockMarketMorningOpenTime(currentTime)
	morningCloseTime := getStockMarketMorningCloseTime(currentTime)
	afternoonOpenTime := getStockMarketAfternoonOpenTime(currentTime)
	afternoonCloseTime := getStockMarketAfternoonCloseTime(currentTime)
	if !currentTime.Before(morningOpenTime) && !currentTime.After(morningCloseTime) {
		return true
	}
	if !currentTime.Before(afternoonOpenTime) && !currentTime.After(afternoonCloseTime) {
		return true
	}
	return false
}

type DayType struct {
	weekday   time.Weekday
	isWorkday bool
	isQueried bool
}

var todayType DayType = DayType{time.Sunday, false, false}

func IsValidDateTime() bool {
	// 获取当前时间
	currentTime := time.Now()
	// 获取当前时间的星期
	weekdayTemp := currentTime.Weekday()
	updateCurrentWeekday(weekdayTemp)
	// 判断是否是周一至周五的工作时间
	if !(todayType.weekday >= time.Monday && todayType.weekday <= time.Friday) {
		return false
	}
	// 判断是否在9:30~11:30或13:00~15:00内内
	isValidTime := isStockMarketOpenTime(currentTime)
	if !isValidTime {
		return false
	}
	// 调用三方api查询当天是否是节假日
	if todayType.isQueried {
		return todayType.isWorkday
	}
	isHoliday, queriedErr := holiday.IsHoliday()
	todayType.isWorkday = !isHoliday
	todayType.isQueried = queriedErr == nil
	return todayType.isWorkday
}

func updateCurrentWeekday(weekdayTemp time.Weekday) {
	if todayType.weekday != weekdayTemp {
		todayType.weekday = weekdayTemp
		todayType.isQueried = false
	}
}
