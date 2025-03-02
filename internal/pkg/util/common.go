package util

import (
	"math/rand"
	"os"
	"strconv"
	"time"
)

// RandStringRunes 返回随机字符串
func RandStringRunes(n int) string {
	var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// GetEnvAsDuration 从环境变量中获取 duration 类型的值
func GetEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return time.Duration(intValue) * time.Second
		}
	}
	return defaultValue
}

// GetEnvAsInt 从环境变量中获取 int 类型的值
func GetEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
