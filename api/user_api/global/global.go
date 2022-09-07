package global

import (
	ut "github.com/go-playground/universal-translator"
	"go-shop/api/user_api/config"
	"go-shop/api/user_api/proto"
)

var (
	Trans ut.Translator

	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	NacosConfig *config.NacosConfig = &config.NacosConfig{}

	UserSrvClient proto.UserClient
)
