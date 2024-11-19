package main

import (
	"math/rand/v2"
	"stock-god-scraper/config"
	"stock-god-scraper/message/telegram"
	"stock-god-scraper/stock/god"
	"stock-god-scraper/stock/market"
	"time"
)

func main() {
	config.Init()
	run()

	// 创建定时器，每 ${SCRAPER_TIME_DURATION} min执行一次
	ticker := time.NewTicker(time.Duration(config.GetConfig().ScraperTimeDuration()) * time.Minute) // 定时器7min触发一次
	defer ticker.Stop()

	for range ticker.C {
		time.Sleep(time.Duration(rand.N(100)) * time.Millisecond) // 随机等待1min
		run()
	}
}

func run() {
	if market.IsValidDateTime() {
		data, ok := god.FetchSourceData()
		if ok {
			telegram.SendMessage(&data)
		}
	}
}
