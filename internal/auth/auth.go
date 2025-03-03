package auth

import (
	"my-token-ai-be/internal/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT 验证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 JWT Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		// 去掉 "Bearer " 前缀，获取纯 Token 字符串
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 解析 JWT Token
		claims, err := ParseToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "Invalid token")
			c.Abort()
			return
		}

		// 将 claims 存入上下文
		c.Set("claims", claims)
		c.Next()
	}
}
