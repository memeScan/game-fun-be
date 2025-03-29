package cron

import (
	"time"

	"game-fun-be/internal/constants"
	"game-fun-be/internal/model"
	"game-fun-be/internal/pkg/util"
	"game-fun-be/internal/redis"
)

const (
	TokenConfigTTL = 6 * time.Minute // TokenConfig在Redis中的过期时间
)

// SyncTokenConfigToRedis 将TokenConfig数据同步到Redis
func SyncTokenConfigToRedis() error {
	start := time.Now()
	util.Log().Info("开始同步TokenConfig数据到Redis")

	// 1. 从数据库获取所有TokenConfig数据
	tokenConfigs, err := model.GetAllTokenConfigs()
	if err != nil {
		util.Log().Error("获取TokenConfig数据失败: %v", err)
		return err
	}

	// 2. 将数据存储到Redis
	if err := redis.Set(constants.TokenConfigRedisKey, tokenConfigs, TokenConfigTTL); err != nil {
		util.Log().Error("存储TokenConfig数据到Redis失败: %v", err)
		return err
	}

	util.Log().Info("同步TokenConfig数据到Redis完成, 耗时: %v", time.Since(start))
	return nil
}
