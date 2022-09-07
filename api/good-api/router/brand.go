package router

import (
	"github.com/gin-gonic/gin"
	"go-shop/api/good-api/api/brands"
)

func InitBrandRouter(Router *gin.RouterGroup) {
	BrandRouter := Router.Group("brands")
	{
		BrandRouter.GET("", brands.BrandList)          // 品牌清單
		BrandRouter.DELETE("/:id", brands.DeleteBrand) // 删除品牌
		BrandRouter.POST("", brands.NewBrand)          //新建品牌
		BrandRouter.PUT("/:id", brands.UpdateBrand)    //修改品牌
	}

	CategoryBrandRouter := Router.Group("categorybrands")
	{
		CategoryBrandRouter.GET("", brands.CategoryBrandList)          //商品類別與品牌關係清單
		CategoryBrandRouter.DELETE("/:id", brands.DeleteCategoryBrand) //删除商品類別與品牌關係
		CategoryBrandRouter.POST("", brands.NewCategoryBrand)          //新建商品類別與品牌關係
		CategoryBrandRouter.PUT("/:id", brands.UpdateCategoryBrand)    //修改商品類別與品牌關係
		CategoryBrandRouter.GET("/:id", brands.GetCategoryBrandList)   //獲取商品類別與品牌關係
	}
}
