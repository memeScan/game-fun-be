package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"my-token-ai-be/internal/constants"
	"my-token-ai-be/internal/pkg/httpUtil"
	"my-token-ai-be/internal/pkg/util"
	"my-token-ai-be/internal/redis"
)

func FetchAndStoreSolPrice(ctx context.Context) error {
	// 获取 SOL 价格
	solPrice, err := httpUtil.GetSolPrice()
	if err != nil {
		return fmt.Errorf("failed to get SOL price: %w", err)
	}

	// 将价格转换为整数（乘以 10^8 并四舍五入）
	intPrice := int64(math.Round(solPrice * float64(constants.SolPriceMultiplier)))

	// 存储整数价格到 Redis
	err = redis.Set(constants.RedisKeySolLatestPrice, fmt.Sprintf("%d", intPrice), 5*time.Minute)
	if err != nil {
		return fmt.Errorf("failed to store price in Redis: %w", err)
	}
	return nil
}

// GetURIContent 获取 URI 内容的函数，maxRetries 指定最大重试次数
func GetURIContent(uri string, maxRetries int) (string, error) {
	// 使用 FetchURIWithRetry，设置最大重试次数为 maxRetries
	body, err := httpUtil.FetchURIWithRetry(uri, maxRetries)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URI content: %w", err)
	}

	return string(body), nil
}

// HasSocialMedia 检查 URI 内容中是否包含社交媒体信息
func HasSocialMedia(content string) bool {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(content), &data); err != nil {
		util.Log().Error("Failed to parse URI content: %v", err)
		return false
	}

	// 检查社交媒体字段是否存在且不为空
	socialFields := []string{"twitter", "telegram", "website"}
	for _, field := range socialFields {
		if value, exists := data[field].(string); exists && value != "" {
			return true
		}
	}

	return false
}
