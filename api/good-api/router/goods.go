package router

import (
	"github.com/gin-gonic/gin"
	"go-shop/api/good-api/middlewares"

	"go-shop/api/good-api/api/goods"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("goods")
	{
		GoodsRouter.GET("", goods.List)                                                                 //商品列表
		GoodsRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.New)               //改接口需要管理员权限
		GoodsRouter.GET("/:id", goods.Detail)                                                           //獲取商品資訊
		GoodsRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Delete)      //刪除商品
		GoodsRouter.GET("/:id/stocks", goods.Stocks)                                                    //取得商品的庫存
		GoodsRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.Update)         //更新商品
		GoodsRouter.PATCH("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), goods.UpdateStatus) //更改商品上下架狀態
	}
}
