package middlewares

import (
	"github.com/gin-gonic/gin"
	"go-shop/api/good-api/models"
	"net/http"
)

func IsAdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		currentUser := claims.(*models.CustomClaims)

		if currentUser.AuthorityId != 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg": "無權限",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}

}
