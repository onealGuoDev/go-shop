package global

import (
	ut "github.com/go-playground/universal-translator"
	"go-shop/api/order_api/config"
	"go-shop/api/order_api/proto"
)

var (
	Trans ut.Translator

	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	NacosConfig *config.NacosConfig = &config.NacosConfig{}

	GoodsSrvClient proto.GoodsClient

	OrderSrvClient proto.OrderClient

	InventorySrvClient proto.InventoryClient
)
