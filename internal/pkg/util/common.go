package util

import (
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

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

func GenerateInviteCode(address string) string {

	// 定义可用字符集（去掉 I、L、O，使用大写字母）
	const baseChars = "ABCDEFGHJKMNPQRSTUVWXYZ0123456789"

	// 创建一个6字节的随机数
	var codeBuilder strings.Builder
	codeBuilder.WriteString("game-") // 添加固定前缀

	// 创建一个6字节的随机数
	b := make([]byte, 6)
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		_, err := crand.Read(b) // 尝试读取随机字节
		if err == nil {
			// 生成6位随机字符
			for i := 0; i < 6; i++ {
				index := int(b[i]) % len(baseChars)
				codeBuilder.WriteByte(baseChars[index])
			}
			return codeBuilder.String()
		}
		// 如果失败，打印错误并在最后一次失败时返回错误
		fmt.Printf("Attempt %d failed: %v, retrying...\n", i+1, err)
		if i == maxRetries-1 {
			return ""
		}
		time.Sleep(10 * time.Millisecond) // 短暂休眠，避免过于频繁重试
	}
	// 理论上不会到达这里，但为了编译通过添加返回
	return ""
}
