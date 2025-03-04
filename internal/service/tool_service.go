package service

import (
	"errors"
	"fmt"
	"game-fun-be/internal/constants"
	"game-fun-be/internal/model"
	"game-fun-be/internal/pkg/util"
	"game-fun-be/internal/redis"
	"game-fun-be/internal/response"
	"time"

	"github.com/shopspring/decimal"
)

type ToolService struct{}

// ResetPoolInfo 重置所有代币的流动池相关信息
func (s *ToolService) ResetPoolInfo(startID ...int64) response.Response {
	const batchSize = 1000
	var lastID int64 = 0
	totalProcessed := 0
	tokenInfoService := &TokenInfoService{}
	tokenTxService := &TokenTransactionService{}
	liquidityPoolService := &TokenLiquidityPoolService{}

	if len(startID) > 0 {
		lastID = startID[0]
	}

	for {
		// 使用游标分页获取代币信息
		tokens, err := model.ListTokenInfosByCursor(lastID, uint8(model.ChainTypeSolana), uint8(model.CreatedPlatformTypePump), false, batchSize)
		if err != nil {
			util.Log().Error("获取代币信息失败: %v", err)
			return response.Err(response.CodeDBError, "获取代币信息失败", err)
		}

		if len(tokens) == 0 {
			break
		}

		var tokenInfos []*model.TokenInfo
		var poolsToUpdate []*model.TokenLiquidityPool
		// 处理这一批数据
		for i := range tokens {
			// 获取最新交易记录
			latestTx, err := tokenTxService.GetLatestTokenTransaction(tokens[i].TokenAddress, tokens[i].ChainType)
			if err != nil {
				util.Log().Error("获取最新交易记录失败: %v", err)
				continue
			}
			if latestTx == nil {
				continue
			}

			// 内盘时使用实际 SOL 储备计算流动性（美元价值）
			if latestTx.RealNativeReserves != 0 {
				// 使用实际 SOL 储备计算
				adjustedNativeReserves := decimal.NewFromInt(int64(latestTx.RealNativeReserves)).Shift(-9)
				// 计算 SOL 的美元价值
				solValue := adjustedNativeReserves.Mul(latestTx.NativePriceUSD)
				// 由于是双边流动性，总流动性��� SOL 价值的 2 倍
				tokens[i].Liquidity = solValue.Mul(decimal.NewFromInt(2))
			}
			tokens[i].PoolAddress = latestTx.PoolAddress

			tokenInfos = append(tokenInfos, &tokens[i])

			pool, err := model.GetPoolInfoByAddressOrderByPoolPcReserve(tokens[i].TokenAddress, uint8(model.ChainTypeSolana), uint8(model.CreatedPlatformTypePump))
			if err != nil {
				util.Log().Error("查询代币池子失败: %v", err)
				continue
			}
			if pool != nil {
				if latestTx.RealNativeReserves != 0 {
					pool.RealNativeReserves = latestTx.RealNativeReserves
					pool.RealTokenReserves = latestTx.RealTokenReserves
					pool.UpdateTime = time.Now()
					poolsToUpdate = append(poolsToUpdate, pool)
				}
			}

		}

		// 批量更新到数据库
		resp := tokenInfoService.BatchUpdateTokenInfo(tokenInfos)
		if resp.Code != 0 {
			util.Log().Error("批量更新失败: %v", resp.Error)
			return response.Err(response.CodeDBError, "批量更新失败", errors.New(resp.Error))
		}

		// 记录本批次更新数量
		util.Log().Info("本批次成功更新 %d 个代币信息", len(tokenInfos))
		totalProcessed += len(tokenInfos)
		util.Log().Info("累计已更新 %d 个代币信息", totalProcessed)

		// 批量删除缓存
		var keys []string
		for _, info := range tokenInfos {
			redisKey := fmt.Sprintf("%s:%s", constants.RedisKeyPrefixTokenInfo, info.TokenAddress)
			keys = append(keys, redisKey)
		}
		if err := redis.Del(keys...); err != nil {
			util.Log().Error("批量删除缓存失败: %v", err)
		}

		lastID = tokens[len(tokens)-1].ID

		if len(poolsToUpdate) > 0 {
			resp := liquidityPoolService.BatchUpdatePools(poolsToUpdate)
			if resp.Code != 0 {
				util.Log().Error("批量更新池子信息失败: %v", resp.Error)
			}
			util.Log().Info("成功更新 %d 个池子的信息", len(poolsToUpdate))
		}

		util.Log().Info("已处理 %d 条记录，最后ID: %d", totalProcessed, lastID)
	}

	return response.Response{
		Code: 0,
		Msg:  "重置流动池信息完成",
		Data: map[string]interface{}{
			"total_processed": totalProcessed,
		},
	}
}

// 获取服务实例
func NewToolService() *ToolService {
	return &ToolService{}
}
