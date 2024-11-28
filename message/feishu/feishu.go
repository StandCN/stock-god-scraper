package feishu

import (
	"encoding/json"
	"log"
	"stock-god-scraper/config"
	"stock-god-scraper/request"
	"stock-god-scraper/stock"
	"time"
)

func FormatMessage[D stock.SourceData](card D) string {
	msg := map[string]interface{}{
		"msg_type": "interactive",
		"card": map[string]interface{}{
			"elements": []map[string]interface{}{
				{
					"tag": "div",
					"text": map[string]string{
						"content": "<at id=all></at>", //取值须使用 open_id 或 user_id 来 @ 指定人
						"tag":     "lark_md",
					},
				},
				{
					"tag": "div",
					"text": map[string]interface{}{
						"tag":     "lark_md",
						"content": "**当前时间**: \n" + time.Now().Format("2006-01-02 15:04:05 +0800"),
					},
				},
				{
					"tag": "div",
					"text": map[string]interface{}{
						"tag":     "lark_md",
						"content": "**微博发送时间**: \n" + card.GetDate().Format("2006-01-02 15:04:05 +0800"),
					},
				},
				{
					"tag": "div",
					"text": map[string]interface{}{
						"tag":     "lark_md",
						"content": "**微博正文**:\n" + card.GetText(),
					},
				},
				{
					"tag": "div",
					"text": map[string]interface{}{
						"tag":        "plain_text",
						"content":    "毛哥牛逼! ",
						"text_color": "carmine-300",
					},
				},
				{
					"tag": "action",
					"actions": []map[string]interface{}{
						{
							"tag":  "button",
							"type": "laser",
							"icon": map[string]string{
								"tag":   "standard_icon",
								"token": "weibo_filled",
							},
							"text": map[string]string{
								"content": "点击查看",
								"tag":     "lark_md",
							},
							"url":   card.GetUrl(),
							"value": map[string]string{},
						},
					},
				},
			},
			"header": map[string]interface{}{
				"title": map[string]string{
					"content": "股神微博",
					"tag":     "plain_text",
				},
				"template": "blue",
			},
		},
	}
	jsonString, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("序列化消息失败: %v", err)
	}
	return string(jsonString)

}

func SendMessage[D stock.SourceData](card D) {
	var msg = FormatMessage(card)

	if config.GetConfig().Debug() {
		log.Printf("将要发送到telegram的消息为: %s", msg)
	}

	// send message
	resp, err := request.GetClient().
		SetHeader("Content-Type", "application/json").
		SetBody(msg).
		Post("https://open.feishu.cn/open-apis/bot/v2/hook/" + config.GetConfig().FeishuBotHookToken())

	if err != nil {
		log.Printf("发送消息失败: %v", err)
	}

	// 定义一个map来存储解析的JSON数据
	var responseMap map[string]interface{}

	// 解析响应体为JSON
	if err := json.Unmarshal(resp.Body(), &responseMap); err != nil {
		log.Fatalf("解析JSON失败: %v. 消息: %v", err, string(resp.Body()))
	}

	// 发送消息
	postResult, ok := responseMap["code"].(float64)
	if !ok || postResult != 0 {
		log.Fatalf("发送消息失败。消息: %s, response: %v", msg, responseMap)
	}
}
