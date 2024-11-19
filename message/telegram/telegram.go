package telegram

import (
	"encoding/json"
	"log"
	"stock-god-scraper/config"
	"stock-god-scraper/request"
	"stock-god-scraper/stock"
	"time"
)

func FormatMessage[D stock.SourceData](card D) string {
	return "_当前时间_: " + time.Now().Format("2006-01-02 15:04:05 +0800") + "\n" +
		"_微博发送时间_: " + card.GetDate().Format("2006-01-02 15:04:05 +0800") + "\n" +
		"_微博正文_: \n" +
		"```html" + "\n" +
		card.GetText() + "\n" +
		"```" + "\n" +
		"_微博地址_: [点击查看](" + card.GetUrl() + ")"
}

func SendMessage[D stock.SourceData](card D) {
	var msg = FormatMessage(card)

	if config.GetConfig().Debug() {
		log.Printf("将要发送到telegram的消息为: %s", msg)
	}

	// send message
	resp, err := request.GetClientWithProxy().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{
			"chat_id":           config.GetConfig().TelegramChatId(),
			"text":              msg,
			"message_thread_id": config.GetConfig().TelegramMessageThreadId(),
			"parse_mode":        "Markdown",
		}).
		Post("https://api.telegram.org/bot" + config.GetConfig().TelegramBotToken() + "/sendMessage")

	if err != nil {
		log.Printf("发送消息失败: %v", err)
	}

	// 定义一个map来存储解析的JSON数据
	var responseMap map[string]interface{}

	// 解析响应体为JSON
	if err := json.Unmarshal(resp.Body(), &responseMap); err != nil {
		log.Fatalf("解析JSON失败: %v", err)
	}

	// 获取cards数组
	postResult, ok := responseMap["ok"].(bool)
	if !ok || !postResult {
		log.Fatalf("发送消息失败。消息: %s, response: %v", msg, responseMap)
	}
}
