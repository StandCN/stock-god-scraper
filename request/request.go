package request

import (
	"stock-god-scraper/config"

	"github.com/go-resty/resty/v2"
)

var (
	client          *resty.Client
	clientWithProxy *resty.Client
)

func GetClient() *resty.Request {
	if client == nil {
		client = resty.New()
		client.SetDebug(config.GetConfig().Debug())
	}
	return client.R()
}

func GetClientWithProxy() *resty.Request {
	if clientWithProxy == nil {
		clientWithProxy = resty.New()
		clientWithProxy.SetDebug(config.GetConfig().Debug())
		clientWithProxy.SetProxy(config.GetConfig().ProxyUrl())
	}
	return clientWithProxy.R()
}
