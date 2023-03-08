package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"utils/jwt"
)

func JWTAuthMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// 根据实际情况取TOKEN, 这里从request header取
		tokenStr := ctx.Request.Header.Get("Authorization")
		if tokenStr == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 1,
				"msg":  "ERR_AUTH_NULL",
			})
			return
		}

		claims, err := jwt.ParseToken(tokenStr)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 2,
				"msg":  "ERR_AUTH_INVALID",
			})
			return
		}
		
		// 此处已经通过了, 可以把Claims中的有效信息拿出来放入上下文使用
		ctx.Set("username", claims.Username)
		ctx.Next()
	}
}
