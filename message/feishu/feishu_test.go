package feishu_test

import (
	"stock-god-scraper/config"
	"stock-god-scraper/message/feishu"
	"stock-god-scraper/stock/god"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func init() {
	viper.AddConfigPath("../../config")
	godotenv.Load("../../.env")
	// 在测试开始前，初始化配置
	config.Init()
}

// 创建一个东八区时区（+0800）
var zone = time.FixedZone("CST", 8*60*60) // UTC+8 时区，8 小时 * 60 分钟 * 60 秒

func TestSendMessage(t *testing.T) {
	card := god.WeiboCardData{
		CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).In(zone),
		Text:      "test",
		Id:        "test",
	}
	feishu.SendMessage(&card)
}
