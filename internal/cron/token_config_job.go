package cron

import (
	"encoding/json"
	"time"

	"game-fun-be/internal/constants"
	"game-fun-be/internal/model"
	"game-fun-be/internal/pkg/util"
	"game-fun-be/internal/redis"
)

const (
	TokenConfigTTL = 10 * time.Minute // TokenConfig在Redis中的过期时间
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

	// 2. 将TokenConfig数据序列化为JSON
	tokenConfigsJSON, err := json.Marshal(tokenConfigs)
	if err != nil {
		util.Log().Error("序列化TokenConfig数据失败: %v", err)
		return err
	}

	// 3. 将数据存储到Redis
	if err := redis.Set(constants.TokenConfigRedisKey, string(tokenConfigsJSON), TokenConfigTTL); err != nil {
		util.Log().Error("存储TokenConfig数据到Redis失败: %v", err)
		return err
	}

	util.Log().Info("同步TokenConfig数据到Redis完成, 耗时: %v", time.Since(start))
	return nil
}
