package global

import (
	ut "github.com/go-playground/universal-translator"
	"go-shop/api/oss_api/config"
)

var (
	Trans ut.Translator

	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	NacosConfig *config.NacosConfig = &config.NacosConfig{}
)
