package global

import (
	ut "github.com/go-playground/universal-translator"
	"go-shop/api/good-api/config"
	"go-shop/api/good-api/proto"
)

var (
	Trans ut.Translator

	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	NacosConfig *config.NacosConfig = &config.NacosConfig{}

	GoodsSrvClient proto.GoodsClient
)
