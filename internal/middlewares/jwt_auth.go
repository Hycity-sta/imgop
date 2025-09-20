package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"imgop/internal/utils"
)

// JWTAuthMiddleware JWT 认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 中获取 Token
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供认证 Token"})
			c.Abort()
			return
		}

		// 去掉 "Bearer " 前缀
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// 解析 Token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的 Token"})
			c.Abort()
			return
		}

		// 将 Claims 存入上下文
		c.Set("claims", claims)
		c.Next()
	}
}
