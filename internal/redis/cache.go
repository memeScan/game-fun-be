package redis

import (
	"context"
	"game-fun-be/internal/pkg/util"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

// RedisClient Redis缓存客户端单例
var RedisClient *redis.Client

// Redis 在中间件中初始化redis链接
func Redis() {
	db, _ := strconv.ParseUint(os.Getenv("REDIS_DB"), 10, 64)
	client := redis.NewClient(&redis.Options{
		Addr:       os.Getenv("REDIS_ADDR"),
		Username:   os.Getenv("REDIS_USER"),
		Password:   os.Getenv("REDIS_PW"),
		DB:         int(db),
		MaxRetries: 1,
	})

	_, err := client.Ping(context.Background()).Result()

	if err != nil {
		util.Log().Panic("连接Redis不成功: %v", err)
	}

	// 打印Redis版本
	version, err := client.Info(context.Background(), "server").Result()
	if err == nil {
		util.Log().Info("Redis connected successfully. Version info: %v", version)
	}

	RedisClient = client
}
