// holiday/holiday.go
package holiday

import (
	"log"
	"stock-god-scraper/request"
	"strings"
	"time"
)

const (
	holidayApi        = "http://api.haoshenqi.top/holiday/today"
	isHolidayResponse = "休息"
)

func IsHoliday() (bool, error) {
	resp, err := request.GetClient().Get(holidayApi)
	if err != nil {
		log.Fatalf("查询法定节假日请求失败: %v", err)
		return true, err
	}
	if strings.Contains(string(resp.Body()), isHolidayResponse) {
		return false, nil
	}
	return true, nil
}

type TodayType struct {
	isWorkday bool
	isQueried bool
	weekday   time.Weekday
}

var todayType TodayType

func UpdateCurrentWeekday(weekdayTemp time.Weekday) {
	if todayType.weekday != weekdayTemp {
		todayType.weekday = weekdayTemp
		todayType.isQueried = false
	}
}
