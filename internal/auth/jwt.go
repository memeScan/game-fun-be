package auth

import (
	"errors"
	"fmt"
	"game-fun-be/internal/redis"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// Claims 自定义的 JWT Claims
type Claims struct {
	Address string `json:"address"`
	UserID  string `json:"user_id"` // 用户信息
	jwt.RegisteredClaims
}

// GenerateJWT 生成 JWT token
func GenerateJWT(address string, userID string, expireDuration time.Duration) (string, time.Time, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(expireDuration)

	// 创建声明
	claims := Claims{
		Address: address,
		UserID:  userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			NotBefore: jwt.NewNumericDate(nowTime),
			Issuer:    "gmgn.exchange",
		},
	}

	// 使用 HS256 签名方法创建 token
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并返回 token
	token, err := tokenClaims.SignedString(jwtSecret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return token, expireTime, nil
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
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// CacheToken 缓存 JWT token
func CacheToken(address, token string, expireDuration time.Duration) error {
	// 将 token 缓存到 Redis
	key := fmt.Sprintf("gmgn:exchange:token:%s", address)
	err := redis.Set(key, token, expireDuration)
	if err != nil {
		return fmt.Errorf("failed to cache token: %w", err)
	}
	return nil
}
