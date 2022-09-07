package initialize

import (
	"github.com/gin-gonic/gin"
	"go-shop/api/good-api/middlewares"
	"go-shop/api/good-api/router"
	"net/http"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	//CORS
	Router.Use(middlewares.Cors())

	ApiGroup := Router.Group("/g/v1")
	router.InitGoodsRouter(ApiGroup)
	router.InitCategoryRouter(ApiGroup)
	router.InitBannerRouter(ApiGroup)
	router.InitBrandRouter(ApiGroup)

	return Router
}
