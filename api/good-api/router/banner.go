package router

import (
	"github.com/gin-gonic/gin"
	"go-shop/api/good-api/api/banners"
	"go-shop/api/good-api/middlewares"
)

func InitBannerRouter(Router *gin.RouterGroup) {
	BannerRouter := Router.Group("banners")
	{
		BannerRouter.GET("", banners.List)                                                            // 輪播圖清單
		BannerRouter.DELETE("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), banners.Delete) // 刪除輪播圖
		BannerRouter.POST("", middlewares.JWTAuth(), middlewares.IsAdminAuth(), banners.New)          // 新建輪播圖
		BannerRouter.PUT("/:id", middlewares.JWTAuth(), middlewares.IsAdminAuth(), banners.Update)    //修改輪播圖資訊
	}
}
