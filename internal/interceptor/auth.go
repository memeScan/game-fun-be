package interceptor

import (
	"game-fun-be/internal/auth"
	"game-fun-be/internal/conf"
	"game-fun-be/internal/redis"
	"game-fun-be/internal/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// BearerAuth 验证JWT的中间件
func BearerAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, response.Err(response.CodeUnauthorized, "Authorization header is required", nil))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(401, response.Err(response.CodeUnauthorized, "Authorization header format must be Bearer {token}", nil))
			c.Abort()
			return
		}

		token := parts[1]

		// 解析 token
		claims, err := auth.ParseToken(token)
		if err != nil {
			c.JSON(401, response.Err(response.CodeUnauthorized, "Invalid or expired JWT", err))
			c.Abort()
			return
		}

		// 检查 token 是否存在于 Redis 中
		_, err = redis.Get(response.RedisKeyPrefixToken + claims.Address)
		if err != nil {
			c.JSON(401, response.Err(response.CodeUnauthorized, "Token expired or not found, please reconnect your wallet", nil))
			c.Abort()
			return
		}

		// 将用户信息存储在上下文中
		c.Set("address", claims.Address)
		c.Set("user_id", claims.UserID)

		c.Next()
	}
}

// ApiKeyAuth 验证API Key的中间件
func ApiKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")

		// 如果没有提供API Key
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "API Key is required",
			})
			c.Abort()
			return
		}

		// 验证API Key是否有效
		valid := false
		for _, key := range conf.APIKeys {
			if apiKey == key {
				valid = true
				break
			}
		}

		if !valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid API Key",
			})
			c.Abort()
			return
		}

		// API Key有效，继续处理请求
		c.Next()
	}
}
