package router

import (
	"github.com/gin-gonic/gin"

	"go-shop/api/order_api/api/order"
	"go-shop/api/order_api/middlewares"
)

func InitOrderRouter(Router *gin.RouterGroup) {
	OrderRouter := Router.Group("orders").Use(middlewares.JWTAuth())
	{
		OrderRouter.GET("", order.List)       // 訂單清單
		OrderRouter.POST("", order.New)       // 新建訂單
		OrderRouter.GET("/:id", order.Detail) // 訂單資訊
	}

}
