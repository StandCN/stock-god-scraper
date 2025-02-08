// holiday/holiday.go
package holiday

import (
	"fmt"
	"stock-god-scraper/request"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	holidayApi        = "http://api.haoshenqi.top/holiday/today"
	isHolidayResponse = "工作"
)

func IsHoliday() (bool, error) {
	resp, err := request.GetClient().Get(holidayApi)
	if err != nil {
		logrus.Errorln(fmt.Sprintf("查询法定节假日请求失败: %v", err))
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
