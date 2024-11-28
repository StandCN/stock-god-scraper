package main

import (
	"math/rand/v2"
	"stock-god-scraper/config"
	"stock-god-scraper/message/feishu"
	"stock-god-scraper/message/telegram"
	"stock-god-scraper/stock"
	"stock-god-scraper/stock/god"
	"stock-god-scraper/stock/market"
	"time"
)

func main() {
	config.Init()

	telegramChan := make(chan stock.SourceData)
	feishuChan := make(chan stock.SourceData)

	go func() {
		for sourceData := range telegramChan {
			telegram.SendMessage(sourceData)
		}
	}()
	go func() {
		for sourceData := range feishuChan {
			feishu.SendMessage(sourceData)
		}
	}()

	run(telegramChan, feishuChan)

	// 创建定时器，每 ${SCRAPER_TIME_DURATION} min执行一次
	ticker := time.NewTicker(time.Duration(config.GetConfig().ScraperTimeDuration()) * time.Minute) // 定时器7min触发一次
	defer ticker.Stop()

	for range ticker.C {
		time.Sleep(time.Duration(rand.N(100)) * time.Millisecond) // 随机等待1min
		run(telegramChan, feishuChan)
	}
}

func run(telegramChan, feishuChan chan stock.SourceData) {
	if market.IsValidDateTime() {
		data, ok := god.FetchSourceData()
		if ok {
			telegramChan <- &data
			feishuChan <- &data
		}
	}
}
