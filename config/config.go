package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config interface {
	TelegramBotToken() string
	TelegramChatId() string
	TelegramMessageThreadId() string
	ProxyUrl() string
	Debug() bool
	ScraperTimeDuration() uint16
}
type config struct {
	telegramBotToken        string
	telegramChatId          string
	telegramMessageThreadId string
	proxyUrl                string
	debug                   bool
	scraperTimeDuration     uint16
}

// 确保config实现了Config接口
var _ Config = (*config)(nil)

func (c config) TelegramBotToken() string {
	return c.telegramBotToken
}

func (c config) TelegramChatId() string {
	return c.telegramChatId
}

func (c config) TelegramMessageThreadId() string {
	return c.telegramMessageThreadId
}

func (c config) ProxyUrl() string {
	return c.proxyUrl
}

func (c config) Debug() bool {
	return c.debug
}

func (c config) ScraperTimeDuration() uint16 {
	return c.scraperTimeDuration
}

func Init() {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// 加载config.toml
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	if viper.GetBool("DEBUG") {
		log.Println("Service RUN on DEBUG mode")
		log.Printf("viper.AllSettings(): %v\n", viper.AllSettings())
	}
	instance = &config{
		telegramBotToken:        viper.GetString("TELEGRAM_BOT_TOKEN"),
		telegramChatId:          viper.GetString("TELEGRAM_CHAT_ID"),
		telegramMessageThreadId: viper.GetString("TELEGRAM_MESSAGE_THREAD_ID"),
		proxyUrl:                viper.GetString("PROXY_URL"),
		debug:                   viper.GetBool("DEBUG"),
		scraperTimeDuration:     uint16(viper.GetInt("SCRAPER_TIME_DURATION")),
	}
}

var instance *config

func GetConfig() Config {
	return instance
}
