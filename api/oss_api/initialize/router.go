package initialize

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-shop/api/oss_api/middlewares"
	"go-shop/api/oss_api/router"
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

	Router.LoadHTMLFiles(fmt.Sprintf("oss_api/templates/index.html"))

	Router.StaticFS("/static", http.Dir(fmt.Sprintf("oss_api/static")))
	Router.GET("", func(c *gin.Context) {

		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "posts/index",
		})
	})

	//配置跨域
	Router.Use(middlewares.Cors())

	ApiGroup := Router.Group("/oss/v1")
	router.InitOssRouter(ApiGroup)

	return Router
}
