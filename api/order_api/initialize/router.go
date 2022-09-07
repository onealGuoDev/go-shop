package initialize

import (
	"github.com/gin-gonic/gin"
	"go-shop/api/order_api/middlewares"
	"go-shop/api/order_api/router"
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

	ApiGroup := Router.Group("/o/v1")
	router.InitOrderRouter(ApiGroup)
	router.InitShopCartRouter(ApiGroup)

	return Router
}
