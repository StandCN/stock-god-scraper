package god

import (
	"encoding/json"
	"log"
	"slices"
	"stock-god-scraper/request"
	"stock-god-scraper/stock"
	"strings"
	"time"
)

type WeiboCardData struct {
	CreatedAt time.Time
	Text      string
	Id        string
}

// SourceData 接口的实现
func (w *WeiboCardData) GetText() string {
	return w.Text
}

func (w *WeiboCardData) GetDate() time.Time {
	return w.CreatedAt
}

func (w *WeiboCardData) GetUrl() string {
	return "https://m.weibo.cn/detail/" + w.Id
}

var _ stock.SourceData = (*WeiboCardData)(nil)

func FetchSourceData() (WeiboCardData, bool) {
	// 获取页面信息
	card, err := fetchPageInfo()
	if err != nil {
		log.Printf("抓取页面失败: %v", err)
		return WeiboCardData{}, false
	}
	// 存储id
	if card.Id != "" && strings.Compare(card.Id, lastWeiboId) != 0 {
		lastWeiboId = card.Id
	} else {
		// 表示消息已经发送过或没有有效消息
		return WeiboCardData{}, false
	}
	log.Printf("获取到的微博信息: %v", card)
	return card, true
}

var lastWeiboId string

const (
	// 微博api时间格式 (layout)，对应时间字符串的格式
	layout = "Mon Jan 2 15:04:05 -0700 2006"
)

// 获取页面信息的函数
func fetchPageInfo() (WeiboCardData, error) {

	resp, err := request.GetClient().
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
	processedCards := processCards(cards, func(card map[string]interface{}) WeiboCardData {
		mblog, ok := card["mblog"].(map[string]interface{})
		if !ok {
			return WeiboCardData{CreatedAt: time.Now(), Text: "", Id: ""}
		}
		createdAtString, _ := mblog["created_at"].(string)
		parsedTime, _ := time.Parse(layout, createdAtString)
		text, _ := mblog["text"].(string)
		id, _ := mblog["id"].(string)
		return WeiboCardData{CreatedAt: parsedTime, Text: text, Id: id}
	})

	filtedCards := slices.DeleteFunc(processedCards, func(data WeiboCardData) bool {
		return !strings.Contains(data.Text, "进场") && !strings.Contains(data.Text, "离场")
	})

	slices.SortFunc(filtedCards, func(left WeiboCardData, right WeiboCardData) int {
		return left.CreatedAt.Compare(right.CreatedAt) * (-1)
	})

	if len(filtedCards) == 0 {
		return WeiboCardData{}, nil
	} else {
		return filtedCards[0], nil
	}
}

// processCards 函数用于处理 cards
func processCards(cards []interface{}, processFunc func(map[string]interface{}) WeiboCardData) []WeiboCardData {
	var result []WeiboCardData
	for _, card := range cards {
		cardMap, ok := card.(map[string]interface{})
		if !ok {
			continue
		}
		card := processFunc(cardMap)
		if !card.CreatedAt.IsZero() && card.Text != "" {
			result = append(result, card)
		}
	}
	return result
}
