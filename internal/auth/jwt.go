package auth

import (
	"os"
	"time"
	"github.com/golang-jwt/jwt/v4"
	"errors"
	"my-token-ai-be/internal/redis"
	"my-token-ai-be/internal/response"
	"fmt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// Claims 自定义的 JWT Claims
type Claims struct {
	Address string `json:"address"`
	jwt.RegisteredClaims
}

// GenerateJWT 生成 JWT token
func GenerateJWT(address string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(2 * time.Hour)

	// 创建声明
	claims := Claims{
		Address: address,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			NotBefore: jwt.NewNumericDate(nowTime),
			Issuer:    "MytokenAI",
		},
	}

	// 使用 HS256 签名方法创建 token
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并返回 token
	token, err := tokenClaims.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ParseToken 解析 JWT token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证 token 的签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err // 返回解析错误
	}

	if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
		return claims, nil // 返回解析后的 claims
	}

	return nil, errors.New("invalid token") // 返回无效的 token 错误
}

// GenerateAndCacheJWT 生成 JWT token 并缓存
func GenerateAndCacheJWT(address string) (string, error) {
	token, err := GenerateJWT(address)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	// 解析 token 以获取过期时间
	claims, err := ParseToken(token)
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	// 计算 token 的有效期
	duration := time.Until(claims.ExpiresAt.Time)

	// 将 token 缓存到 Redis
	err = redis.Set(response.RedisKeyPrefixToken + address, token, duration)
	
	if err != nil {
		return "", fmt.Errorf("failed to cache token: %w", err)
	}

	return token, nil
}
