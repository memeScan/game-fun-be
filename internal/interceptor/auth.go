package interceptor

import (
	"my-token-ai-be/internal/model"
	"my-token-ai-be/internal/redis"
	"my-token-ai-be/internal/auth"
	"github.com/gin-gonic/gin"
	"strings"
	"fmt"
	"my-token-ai-be/internal/response"
)

// CurrentUser 获取登录用户
func CurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.Next()
			return
		}

		claims, err := auth.ParseToken(token)
		if err != nil {
			c.Next()
			return
		}

		fmt.Println("Address: ", claims.Address)

		user, err := model.GetUserByAddress(claims.Address)

		if err != nil {
			c.Next()
			return
		}


		c.Set("user", user)
		c.Set("address", user.Address)
		c.Next()
	}
}

// AuthRequired 需要登录
func AuthRequired() gin.HandlerFunc {
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

		c.Next()
	}
}
