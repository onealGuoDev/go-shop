package router

import (
	"github.com/gin-gonic/gin"
	"go-shop/api/order_api/api/shop_cart"
	"go-shop/api/order_api/middlewares"
)

func InitShopCartRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("shopcarts").Use(middlewares.JWTAuth())
	{
		GoodsRouter.GET("", shop_cart.List)          //購物車列表
		GoodsRouter.DELETE("/:id", shop_cart.Delete) //删除購物車商品
		GoodsRouter.POST("", shop_cart.New)          //新增商品到購物車
		GoodsRouter.PATCH("/:id", shop_cart.Update)  //修改購物車內容
	}
}
