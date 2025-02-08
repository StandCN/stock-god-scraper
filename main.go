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

type SendMessageFn func(stock.SourceData) error

func main() {
	config.Init()

	telegramChan := buildChannel(telegram.SendMessage)
	feishuChan := buildChannel(feishu.SendMessage)

	run(telegramChan, feishuChan)

	// 创建定时器，每 ${SCRAPER_TIME_DURATION} min执行一次
	ticker := time.NewTicker(time.Duration(config.GetConfig().ScraperTimeDuration()) * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		time.Sleep(time.Duration(rand.N(100)) * time.Millisecond) // 随机等待
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

func buildChannel(fn SendMessageFn) chan stock.SourceData {
	msgChan := make(chan stock.SourceData)
	go func() {
		for sourceData := range msgChan {
			for {
				err := fn(sourceData)
				if err == nil {
					break
				}
			}
		}
	}()
	return msgChan
}
