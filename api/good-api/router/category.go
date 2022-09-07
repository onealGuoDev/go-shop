package router

import (
	"github.com/gin-gonic/gin"
	"go-shop/api/good-api/api/category"
)

func InitCategoryRouter(Router *gin.RouterGroup) {

	CategoryRouter := Router.Group("categorys")
	{
		CategoryRouter.GET("", category.List)          // 商品類別清單
		CategoryRouter.DELETE("/:id", category.Delete) // 删除商品分類
		CategoryRouter.GET("/:id", category.Detail)    // 獲取商品分類
		CategoryRouter.POST("", category.New)          //新建商品分類
		CategoryRouter.PUT("/:id", category.Update)    //修改商品類別
	}
}
