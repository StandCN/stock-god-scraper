package main

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// 时间格式 (layout)，对应时间字符串的格式
const layout = "Mon Jan 2 15:04:05 -0700 2006"

// 查询节假日api
const holidayApi, isHolidayResponse = "http://api.haoshenqi.top/holiday/today", "休息"

// 创建一个 resty 客户端
var client = resty.New()

// 获取页面信息的函数
func fetchPageInfo() (CardData, error) {

	resp, err := client.R().
		ForceContentType("application/json; charset=utf-8").
		Get("https://m.weibo.cn/api/container/getIndex?luicode=20000061&lfid=5044453196435127&type=uid&value=7551230054&containerid=1076037551230054")

	if err != nil {
		log.Fatalf("请求失败: %v", err)
	}

	if resp.StatusCode() != 200 {
		log.Fatalf("请求失败: %v", resp.Status())
	}

	// 定义一个map来存储解析的JSON数据
	var responseMap map[string]interface{}

	// 解析响应体为JSON
	if err := json.Unmarshal(resp.Body(), &responseMap); err != nil {
		log.Fatalf("解析JSON失败: %v", err)
	}

	// 获取cards数组
	dataMap, ok := responseMap["data"].(map[string]interface{})
	if !ok {
		log.Fatalf("解析 data 字段失败")
	}

	cards, ok := dataMap["cards"].([]interface{})
	if !ok {
		log.Fatalf("解析 cards 字段失败")
	}

	// 使用自定义函数进行"流式处理"
	processedCards := processCards(cards, func(card map[string]interface{}) CardData {
		mblog, ok := card["mblog"].(map[string]interface{})
		if !ok {
			return CardData{time.Now(), "", ""}
		}
		createdAtString, _ := mblog["created_at"].(string)
		parsedTime, _ := time.Parse(layout, createdAtString)
		text, _ := mblog["text"].(string)
		id, _ := mblog["id"].(string)
		return CardData{parsedTime, text, id}
	})

	filtedCards := slices.DeleteFunc(processedCards, func(data CardData) bool {
		return !strings.Contains(data.text, "进场") && !strings.Contains(data.text, "离场")
	})

	slices.SortFunc(filtedCards, func(left CardData, right CardData) int {
		return left.createdAt.Compare(right.createdAt) * (-1)
	})

	if len(filtedCards) == 0 {
		return CardData{}, fmt.Errorf("没有找到符合条件的数据")
	} else {
		return filtedCards[0], nil
	}
}

type CardData struct {
	createdAt time.Time
	text      string
	id        string
}

// processCards 函数用于处理 cards
func processCards(cards []interface{}, processFunc func(map[string]interface{}) CardData) []CardData {
	var result []CardData
	for _, card := range cards {
		cardMap, ok := card.(map[string]interface{})
		if !ok {
			continue
		}
		card := processFunc(cardMap)
		if !card.createdAt.IsZero() && card.text != "" {
			result = append(result, card)
		}
	}
	return result
}

func main() {
	if isValidDateTime() {
		fetchAndLogPageInfo()
	}
	// // 创建定时器，每7min执行一次
	// ticker := time.NewTicker(7 * time.Minute) // 定时器7min触发一次
	// defer ticker.Stop()

	// for range ticker.C {
	// 	time.Sleep(time.Duration(rand.N(100)) * time.Millisecond) // 随机等待1min
	// 	if isValidDateTime() {
	// 		fetchAndLogPageInfo()
	// 	}
	// }
}

type DayType struct {
	weekday   time.Weekday
	isWorkday bool
	isQueried bool
}

var todayType DayType = DayType{time.Sunday, false, false}

func isValidDateTime() bool {
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
	currentHour := currentTime.Hour()
	isValidTime := isStockMarketOpenTime(currentHour, currentTime)
	if !isValidTime {
		return false
	}
	// 调用三方api查询当天是否是节假日
	if todayType.isQueried {
		return todayType.isWorkday
	}
	queryHoliday()
	return todayType.isWorkday
}

func queryHoliday() {
	resp, err := client.R().Get(holidayApi)
	if err != nil {
		log.Fatalf("查询法定节假日请求失败: %v", err)
		todayType.isWorkday = true
		return
	}
	if strings.Contains(string(resp.Body()), isHolidayResponse) {
		todayType.isWorkday = false
	} else {
		todayType.isWorkday = true
	}
	todayType.isQueried = true
}

func isStockMarketOpenTime(currentHour int, currentTime time.Time) bool {
	var isValidTime bool = true
	switch currentHour {
	case 9:
		if currentTime.Minute() >= 30 {
			isValidTime = true
		}
	case 11:
		if currentTime.Minute() < 30 {
			isValidTime = true
		}
	case 10, 13, 14:
		isValidTime = true
	default:
	}
	return isValidTime
}

func updateCurrentWeekday(weekdayTemp time.Weekday) {
	if todayType.weekday != weekdayTemp {
		todayType.weekday = weekdayTemp
		todayType.isQueried = false
	}
}

var lastWeiboId string

func fetchAndLogPageInfo() {
	// 获取页面信息
	card, err := fetchPageInfo()
	if err != nil {
		log.Printf("抓取页面失败: %v", err)
		return
	}
	// 存储id
	if card.id != "" && strings.Compare(card.id, lastWeiboId) != 0 {
		lastWeiboId = card.id
	} else {
		// 表示消息已经发送过或没有有效消息
		return
	}
	// 在这里添加后续处理逻辑，比如存储信息、数据分析等
	sendMessage(&card)
}

func sendMessage(card *CardData) {

	// 打印抓取的内容
	log.Printf("获取的页面信息: 时间: %s, 文本: %s, 微博地址: %s", card.createdAt.Format(layout), card.text, "https://m.weibo.cn/detail/"+card.id)

	// TODO send message
}
