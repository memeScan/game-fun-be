package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"game-fun-be/internal/pkg/util"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// Set 将键值对存储到 Redis
func Set(key string, value interface{}, expiration ...time.Duration) error {
	// 将 value 序列化为 JSON
	jsonValue, err := json.Marshal(value)
	if err != nil {
		util.Log().Error("Error marshaling value for Redis: %v", err)
		return err
	}

	// 如果没有传入 expiration，则永久存储
	var exp time.Duration
	if len(expiration) > 0 {
		exp = expiration[0]
	} else {
		exp = 0 // 0 表示永久存储
	}

	// 使用 Redis 客户端设置值
	err = RedisClient.Set(context.Background(), key, jsonValue, exp).Err()
	if err != nil {
		util.Log().Error("Error setting value in Redis: %v", err)
		return err
	}

	util.Log().Info("Successfully set key %s in Redis", key)
	return nil
}

// Get 从 Redis 获取值
func Get(key string) (string, error) {
	ctx := context.Background()
	value, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // 键不存在,返回空字符串和 nil 错误
	} else if err != nil {
		return "", err // 其他错误,返回空字符串和错误
	}
	return value, nil // 返回获取到的值和 nil 错误
}

// Delete 从 Redis 删除键
func Delete(key string) error {
	ctx := context.Background()
	return RedisClient.Del(ctx, key).Err()
}

// SetNX 如果键不存在，则设置键的值
func SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	ctx := context.Background()
	return RedisClient.SetNX(ctx, key, value, expiration).Result()
}

// Exists 检查键是否存在
func Exists(key string) (bool, error) {
	ctx := context.Background()
	n, err := RedisClient.Exists(ctx, key).Result()
	return n > 0, err
}

// Expire 设置键的过期时间
func Expire(key string, expiration time.Duration) (bool, error) {
	ctx := context.Background()
	return RedisClient.Expire(ctx, key, expiration).Result()
}

// GetAndDelete 从 Redis 获取值并立即删除该键
func GetAndDelete(key string) (string, error) {
	ctx := context.Background()

	// 获取值
	value, err := RedisClient.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return "", err // 返回错误
	}

	// 删除键
	_, delErr := RedisClient.Del(ctx, key).Result()
	if delErr != nil {
		return "", delErr // 返回删除错误
	}

	// 返回获取到的值，如果值不存在，返回空字符串
	return value, nil
}

// MSet 批量设置多个键值对到 Redis
func MSet(keyValues map[string]string, expiration time.Duration) error {
	ctx := context.Background()
	pipe := RedisClient.Pipeline()

	for key, value := range keyValues {
		pipe.Set(ctx, key, value, expiration)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		util.Log().Error("Error batch setting values in Redis: %v", err)
		return err
	}

	util.Log().Info("Successfully batch set %d keys in Redis", len(keyValues))
	return nil
}

// MGet 批量获取多个键的值
func MGet(keys []string) ([]string, error) {
	ctx := context.Background()
	values, err := RedisClient.MGet(ctx, keys...).Result()
	if err != nil {
		util.Log().Error("Error batch getting values from Redis: %v", err)
		return nil, err
	}

	// 将 interface{} 切片转换为 string 切片
	results := make([]string, len(values))
	for i, v := range values {
		if v != nil {
			results[i] = v.(string)
		}
	}

	return results, nil
}

// SAdd 将一个或多个成员添加到集合中
func SAdd(key string, members ...interface{}) error {
	ctx := context.Background()
	err := RedisClient.SAdd(ctx, key, members...).Err()
	if err != nil {
		util.Log().Error("Error adding members to set %s: %v", key, err)
		return err
	}
	return nil
}

// SMembers 获取集合中的所有成员
func SMembers(key string) ([]string, error) {
	ctx := context.Background()
	members, err := RedisClient.SMembers(ctx, key).Result()
	if err != nil {
		util.Log().Error("Error getting members from set %s: %v", key, err)
		return nil, err
	}
	return members, nil
}

// SIsMember 判断成员是否在集合中
func SIsMember(key string, member interface{}) (bool, error) {
	ctx := context.Background()
	return RedisClient.SIsMember(ctx, key, member).Result()
}

// SRem 从集合中移除一个或多个成员
func SRem(key string, members ...interface{}) error {
	ctx := context.Background()
	err := RedisClient.SRem(ctx, key, members...).Err()
	if err != nil {
		util.Log().Error("Error removing members from set %s: %v", key, err)
		return err
	}
	return nil
}

// SCard 获取集合中成员数量
func SCard(key string) (int64, error) {
	ctx := context.Background()
	return RedisClient.SCard(ctx, key).Result()
}

// BatchSAdd 批量将多个成员添加到多个集合中
func BatchSAdd(keyMembers map[string][]interface{}) error {
	ctx := context.Background()
	pipe := RedisClient.Pipeline()

	for key, members := range keyMembers {
		pipe.SAdd(ctx, key, members...)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		util.Log().Error("Error batch adding members to sets: %v", err)
		return err
	}
	return nil
}

// BatchSRem 批量从多个集合中移除多个成员
func BatchSRem(keyMembers map[string][]interface{}) error {
	ctx := context.Background()
	pipe := RedisClient.Pipeline()

	for key, members := range keyMembers {
		pipe.SRem(ctx, key, members...)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		util.Log().Error("Error batch removing members from sets: %v", err)
		return err
	}
	return nil
}

// Z represents sorted set member.
type Z = redis.Z

// ZRangeBy represents sorted set range query options.
type ZRangeBy = redis.ZRangeBy

// ZAdd 添加一个或多个成员到有序集合
func ZAdd(key string, members ...Z) error {
	ctx := context.Background()
	err := RedisClient.ZAdd(ctx, key, members...).Err()
	if err != nil {
		util.Log().Error("Error adding to sorted set %s: %v", key, err)
		return err
	}
	return nil
}

// ZRemRangeByScore 移除有序集合中指定分数范围的成员
func ZRemRangeByScore(key, min, max string) error {
	ctx := context.Background()
	err := RedisClient.ZRemRangeByScore(ctx, key, min, max).Err()
	if err != nil {
		util.Log().Error("Error removing range from sorted set %s: %v", key, err)
		return err
	}
	return nil
}

// ZRangeByScore 获取有序集合中指定分数范围的成员
func ZRangeByScore(key string, opt *ZRangeBy) ([]string, error) {
	ctx := context.Background()
	members, err := RedisClient.ZRangeByScore(ctx, key, opt).Result()
	if err != nil {
		util.Log().Error("Error getting range from sorted set %s: %v", key, err)
		return nil, err
	}
	return members, nil
}

// SAddWithTTL 将成员添加到集合并设置过期时间
func SAddWithTTL(key string, ttl time.Duration, members ...interface{}) error {
	ctx := context.Background()
	pipe := RedisClient.Pipeline()

	// 添加成员到集合
	pipe.SAdd(ctx, key, members...)
	// 设置过期时间
	pipe.Expire(ctx, key, ttl)

	_, err := pipe.Exec(ctx)
	if err != nil {
		util.Log().Error("Error adding members to set %s with TTL: %v", key, err)
		return err
	}
	return nil
}

// Del 批量删除多个键
func Del(keys ...string) error {
	ctx := context.Background()
	return RedisClient.Del(ctx, keys...).Err()
}

// 获取ttl
func TTL(key string) (time.Duration, error) {
	ctx := context.Background()
	return RedisClient.TTL(ctx, key).Result()
}

var ctx = context.Background()

// Lock 尝试获取分布式锁
func Lock(lockKey string, lockValue string, expiration time.Duration, retryInterval time.Duration) (bool, error) {

	timeout := time.After(expiration)

	for {
		// 尝试获取锁
		success, err := RedisClient.SetNX(ctx, lockKey, lockValue, expiration).Result()
		if err != nil {
			return false, err // 返回错误
		}
		if success {
			return true, nil // 成功获取锁
		}

		// 等待重试间隔
		select {
		case <-timeout:
			return false, nil // 超时，返回失败
		case <-time.After(retryInterval):
			// 继续重试
		}
	}
}

// Unlock 释放分布式锁
func Unlock(key string, lockValue string) error {
	// 使用 Lua 脚本释放锁
	luaScript := `
    if redis.call("GET", KEYS[1]) == ARGV[1] then
        return redis.call("DEL", KEYS[1])
    else
        return 0 -- Lock not released
    end
    `

	// 执行 Lua 脚本
	result, err := RedisClient.Eval(ctx, luaScript, []string{key}, lockValue).Result()
	if err != nil {
		return err // 错误处理
	}

	// 检查返回值
	if result.(int64) == 0 {
		return fmt.Errorf("lock not released, it may not belong to the current client")
	}

	return nil // 成功释放锁
}

// AddToken 添加代币地址，设置过期时间
func AddToken(ctx context.Context, key string, tokenAddress string, expireTime int64) error {
	return RedisClient.ZAdd(ctx, key, redis.Z{
		Score:  float64(expireTime),
		Member: tokenAddress,
	}).Err()
}

// GetValidTokens 获取未过期的代币列表
func GetValidTokens(ctx context.Context, key string) ([]string, error) {
	currentTime := time.Now().Unix()
	return RedisClient.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: strconv.FormatInt(currentTime, 10), // 将时间戳转换为字符串
		Max: "+inf",
	}).Result()
}

// RemoveToken 移除代币
func RemoveToken(ctx context.Context, key string, tokenAddress string) error {
	return RedisClient.ZRem(ctx, key, tokenAddress).Err()
}

// UpdateTokenExpireTime 更新代币过期时间
func UpdateTokenExpireTime(ctx context.Context, key string, tokenAddress string, newExpireTime int64) error {
	return RedisClient.ZAdd(ctx, key, redis.Z{
		Score:  float64(newExpireTime),
		Member: tokenAddress,
	}).Err()
}

// IsTokenValid 检查代币是否有效（未过期）
func IsTokenValid(ctx context.Context, key string, tokenAddress string) (bool, error) {
	score, err := RedisClient.ZScore(ctx, key, tokenAddress).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return score > float64(time.Now().Unix()), nil
}

// CreateSortedSet 创建一个新的 Sorted Set
// ttl <= 0 表示永久存储
// CreateSortedSet 创建一个新的 Sorted Set
func CreateSortedSet(key string, ttl int64) error {
	ctx := context.Background()

	// 1. 检查 key 是否已存在
	exists, err := RedisClient.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("check key exists error: %v", err)
	}

	if exists > 0 {
		return nil // key 已存在，不需要创建
	}

	// 2. 创建事务
	pipe := RedisClient.TxPipeline()

	// 3. 创建 sorted set 时需要至少添加一个成员
	pipe.ZAdd(ctx, key, redis.Z{
		Score:  0,
		Member: "init", // 添加一个初始成员
	})

	// 4. 如果需要设置过期时间
	if ttl > 0 {
		pipe.Expire(ctx, key, time.Duration(ttl)*time.Second)
	}

	// 5. 执行事务
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("create sorted set error: %v", err)
	}

	// 6. 删除初始成员（如果需要）
	err = RedisClient.ZRem(ctx, key, "init").Err()
	if err != nil {
		return fmt.Errorf("remove init member error: %v", err)
	}

	return nil
}

// SafeCleanExpiredTokens 安全清理过期数据
func SafeCleanExpiredTokens(ctx context.Context, key string) error {
	// 使用 MULTI/EXEC 确保原子性
	pipe := RedisClient.TxPipeline()
	now := time.Now().Unix()

	// 删除分数小于当前时间的成员
	pipe.ZRemRangeByScore(ctx, key, "-inf", fmt.Sprintf("%d", now))

	_, err := pipe.Exec(ctx)
	return err
}

// SafeBatchAddTokens 安全批量添加数据
func SafeBatchAddTokens(ctx context.Context, key string, tokens map[string]int64, ttl int64) error {
	pipe := RedisClient.TxPipeline()

	// 1. 添加数据
	members := make([]redis.Z, 0, len(tokens))

	for token, _ := range tokens {
		members = append(members, redis.Z{
			Score:  float64(ttl),
			Member: token,
		})
	}

	// 2. 批量添加到 sorted set
	if len(members) > 0 {
		pipe.ZAdd(ctx, key, members...)
	}

	// 4. 执行 pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("execute pipeline error: %v", err)
	}

	return nil
}

// EnsureKeyTTL 确保 key 有正确的过期时间
func EnsureKeyTTL(key string, ttl int64) error {
	ctx := context.Background()

	// 1. 先检查当前 TTL
	currentTTL, err := RedisClient.TTL(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("get TTL error: %v", err)
	}

	// 2. 如果 TTL 小于预期值，则更新
	if currentTTL.Seconds() < float64(ttl) {
		err = RedisClient.Expire(ctx, key, time.Duration(ttl)*time.Second).Err()
		if err != nil {
			return fmt.Errorf("set expire error: %v", err)
		}
	}

	return nil
}

func GetToken(key string) (string, bool, error) {
	ctx := context.Background()
	value, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", false, nil
	} else if err != nil {
		return "", false, err
	}

	var parsedValue string
	if err := json.Unmarshal([]byte(value), &parsedValue); err == nil {
		return parsedValue, true, nil
	}

	return value, true, nil
}

// Unmarshal 将 JSON 字符串解析为指定对象
func Unmarshal(value string, target interface{}) error {
	if value == "" {
		return fmt.Errorf("value is empty")
	}

	if err := json.Unmarshal([]byte(value), target); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}
	return nil
}
