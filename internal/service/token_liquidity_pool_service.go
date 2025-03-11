package service

import (
	"game-fun-be/internal/model"
	"game-fun-be/internal/pkg/httpUtil"
	"game-fun-be/internal/pkg/util"
	"game-fun-be/internal/response"
	"hash/fnv"
	"time"

	"github.com/shopspring/decimal"
)

// TokenLiquidityPoolService 代币流动性池服务
type TokenLiquidityPoolService struct{}

// CreateTokenLiquidityPool 创建代币流动性池记录
func (service *TokenLiquidityPoolService) CreateTokenLiquidityPool(pool *model.TokenLiquidityPool) (*model.TokenLiquidityPool, error) {
	err := model.CreateTokenLiquidityPool(pool)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

// GetTokenLiquidityPoolByAddress 通过池子地址和链类型获取代币流动性池记录
func (service *TokenLiquidityPoolService) GetTokenLiquidityPoolByAddress(poolAddress string, chainType uint8) (*model.TokenLiquidityPool, error) {
	return model.GetTokenLiquidityPoolByAddress(poolAddress, chainType)
}

// UpdateTokenLiquidityPool 更新代币流动性池记录
func (service *TokenLiquidityPoolService) UpdateTokenLiquidityPool(pool *model.TokenLiquidityPool) (*model.TokenLiquidityPool, error) {
	err := model.UpdateTokenLiquidityPool(pool)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

// DeleteTokenLiquidityPool 删除代币流动性池记录
func (service *TokenLiquidityPoolService) DeleteTokenLiquidityPool(id uint64) error {
	return model.DeleteTokenLiquidityPool(id)
}

// ListTokenLiquidityPools 列出代币流动性池记录
func (service *TokenLiquidityPoolService) ListTokenLiquidityPools(limit, offset int) ([]model.TokenLiquidityPool, error) {
	return model.ListTokenLiquidityPools(limit, offset)
}

// ProcessTokenLiquidityPoolCreation 处理代币流动性池记录创建
func (service *TokenLiquidityPoolService) ProcessTokenLiquidityPoolCreation(pool *model.TokenLiquidityPool) response.Response {
	_, err := service.CreateTokenLiquidityPool(pool)
	if err != nil {
		return response.Err(response.CodeDBError, "Failed to create token liquidity pool", err)
	}

	return response.Response{
		Code: 0,
		Data: pool,
		Msg:  "Token liquidity pool created successfully",
	}
}

// ProcessTokenLiquidityPoolUpdate 处理代币流动性池记录更新
func (service *TokenLiquidityPoolService) ProcessTokenLiquidityPoolUpdate(pool *model.TokenLiquidityPool) response.Response {
	_, err := service.UpdateTokenLiquidityPool(pool)
	if err != nil {
		return response.Err(response.CodeDBError, "Failed to update token liquidity pool", err)
	}

	return response.Response{
		Code: 0,
		Data: pool,
		Msg:  "Token liquidity pool updated successfully",
	}
}

// ProcessTokenLiquidityPoolQuery 处理代币流动性池记录查询
func (service *TokenLiquidityPoolService) ProcessTokenLiquidityPoolQuery(poolAddress string, chainType uint8) response.Response {
	pool, err := service.GetTokenLiquidityPoolByAddress(poolAddress, chainType)
	if err != nil {
		return response.Err(response.CodeDBError, "Failed to query token liquidity pool", err)
	}

	return response.Response{
		Code: 0,
		Data: pool,
		Msg:  "Token liquidity pool queried successfully",
	}
}

// ProcessTokenLiquidityPoolList 处理代币流动性池记录列表查询
func (service *TokenLiquidityPoolService) ProcessTokenLiquidityPoolList(limit, offset int) response.Response {
	pools, err := service.ListTokenLiquidityPools(limit, offset)
	if err != nil {
		return response.Err(response.CodeDBError, "Failed to list token liquidity pools", err)
	}

	return response.Response{
		Code: 0,
		Data: pools,
		Msg:  "Token liquidity pools listed successfully",
	}
}

// ConvertMessageToLiquidityPool 将 RaydiumCreateMessage 转换为 TokenLiquidityPool
func (service *TokenLiquidityPoolService) ConvertMessageToLiquidityPool(msg *model.RaydiumCreateMessage) *model.TokenLiquidityPool {
	// 转换代币数量为 decimal
	poolPcReserve, err := decimal.NewFromString(msg.PoolBaseReserve)
	if err != nil {
		util.Log().Error("Error parsing PoolPcReserve: %v", err)
		poolPcReserve = decimal.Zero
	}

	poolCoinReserve, err := decimal.NewFromString(msg.PoolQuoteReserve)
	if err != nil {
		util.Log().Error("Error parsing PoolCoinReserve: %v", err)
		poolCoinReserve = decimal.Zero
	}

	// 如果类型是初始化，则使用 ChangePoolBaseAmount 作为初始化定价代币数量
	var initialPcReserve, initialCoinReserve decimal.Decimal
	if msg.PoolState == 0 {
		initialPcReserve, err = decimal.NewFromString(msg.ChangePoolBaseAmount)
		if err != nil {
			util.Log().Error("Error parsing ChangePoolBaseAmount: %v", err)
			initialPcReserve = decimal.Zero
		}
		initialCoinReserve, err = decimal.NewFromString(msg.ChangePoolQuoteAmount)
		if err != nil {
			util.Log().Error("Error parsing ChangePoolQuoteAmount: %v", err)
			initialCoinReserve = decimal.Zero
		}
	}

	pool := &model.TokenLiquidityPool{}

	// 设置链类型和平台类型
	pool.ChainType = uint8(model.ChainTypeSolana)        // Solana 链
	pool.PlatformType = uint8(model.PlatformTypeRaydium) // Raydium 平台

	// 设置地址��关字段
	pool.MarketAddress = msg.MarketAddress // 市场地址
	pool.PoolAddress = msg.PoolAddress     // 池子地址
	pool.PcAddress = msg.BaseToken         // 定价代币地址
	pool.CoinAddress = msg.QuoteToken      // 交易代币地址
	pool.UserAddress = msg.User            // 创建者地址

	// 设置代币储备量
	pool.PoolPcReserve = uint64(poolPcReserve.IntPart())           // 池子中定价代币的当前总量
	pool.PoolCoinReserve = uint64(poolCoinReserve.IntPart())       // 池子中交易代币的当前总量
	pool.InitialPcReserve = uint64(initialPcReserve.IntPart())     // 池子初始定价代币总量
	pool.InitialCoinReserve = uint64(initialCoinReserve.IntPart()) // 池子初始交易代币总量
	// 设置区块信息
	pool.Block = msg.Block                       // 区块高度
	pool.BlockTime = time.Unix(msg.Timestamp, 0) // 区块时间

	// 计算交易对哈希值
	pool.PairHash = calculatePairHash(msg.QuoteToken, msg.BaseToken)

	return pool
}

// CreateLiquidityPoolFromMessage 从 Kafka 消息创建流动性池记录
func (service *TokenLiquidityPoolService) CreateLiquidityPoolFromMessage(msg *model.RaydiumCreateMessage) response.Response {
	pool := service.ConvertMessageToLiquidityPool(msg)

	if err := model.UpsertTokenLiquidityPool(pool); err != nil {
		return response.Err(response.CodeDBError, "Failed to upsert liquidity pool", err)
	}

	return response.Response{
		Code: 0,
		Data: pool,
		Msg:  "Liquidity pool created successfully",
	}
}

// calculatePairHash 计算交易对哈希值
func calculatePairHash(pcAddress, coinAddress string) uint64 {
	h := fnv.New64a()

	// 确保相同的交易对产生相同的哈希值，不受地址顺序影响
	if pcAddress < coinAddress {
		h.Write([]byte(pcAddress))
		h.Write([]byte(coinAddress))
	} else {
		h.Write([]byte(coinAddress))
		h.Write([]byte(pcAddress))
	}

	return h.Sum64()
}

// CreatePumpFunInitialPool 创建 PumpFun 初始流动性池
func (service *TokenLiquidityPoolService) CreatePumpFunInitialPool(msg *model.TokenInfoMessage) *model.TokenLiquidityPool {
	pool := &model.TokenLiquidityPool{}

	// 设置链类型和平台类型
	pool.ChainType = uint8(model.ChainTypeSolana)
	pool.PlatformType = uint8(model.PlatformTypePump)

	// 设置地址相关字段
	pool.PcAddress = model.SolanaWrappedSOLAddress // SOL 代币地址
	pool.CoinAddress = msg.Mint                    // Meme 代币地址
	pool.MarketAddress = ""                        // 市场地址
	pool.PoolAddress = msg.BondingCurve            // 池子地址
	pool.UserAddress = msg.Creator                 // 创建者地址

	// 从 bondingCurve 计算初始流动性
	initialPcReserve, _ := decimal.NewFromString(model.PUMP_INITIAL_REALSOL_TOKEN_RESERVES)
	initialCoinReserve, _ := decimal.NewFromString(model.PUMP_INITIAL_VIRTUAL_TOKEN_RESERVES)
	// 设置代币储备量
	pool.InitialPcReserve = uint64(initialPcReserve.IntPart())
	pool.InitialCoinReserve = uint64(initialCoinReserve.IntPart())

	// 设置区块信息
	pool.BlockTime = time.Unix(msg.Timestamp, 0)
	pool.Block = msg.Block

	// 计算交易对哈希值
	pool.PairHash = calculatePairHash(msg.Mint, pool.PcAddress)

	return pool
}

// GetExistingPools 通过池子地址列表批量查询池子信息
func (service *TokenLiquidityPoolService) GetExistingPools(poolAddresses []string) map[string]*model.TokenLiquidityPool {
	// 初始化返回结果map
	result := make(map[string]*model.TokenLiquidityPool)

	// 调用model层批量查询
	pools, err := model.GetTokenLiquidityPoolsByAddresses(poolAddresses, uint8(model.ChainTypeSolana))
	if err != nil {
		util.Log().Error("批量查询池子信息失败: %v", err)
		return result
	}

	// 将结果转换为map
	for _, pool := range pools {
		result[pool.PoolAddress] = pool
	}

	return result
}

// BatchUpdatePools 批量更新池子信息
func (service *TokenLiquidityPoolService) BatchUpdatePools(pools []*model.TokenLiquidityPool) response.Response {
	if err := model.BatchUpdateTokenLiquidityPools(pools); err != nil {
		return response.Err(response.CodeDBError, "Failed to batch update pools", err)
	}
	return response.Response{Code: 0}
}

// 查询代币池子,并返回最大池子,如果没有则调用链端api补池子数据
func QueryAndCheckPool(tokenAddress string, chainType uint8, platformType uint8) (*model.TokenLiquidityPool, error) {
	pool, err := model.GetPoolInfoByAddressOrderByPoolPcReserve(tokenAddress, chainType, platformType)
	if err != nil {
		util.Log().Error("查询代币池子失败: %v", err)
		return nil, err
	}

	if pool == nil {
		// 调用链端api补池子数据
		poolInfo, err := httpUtil.GetPoolInfo([]string{tokenAddress})
		if err != nil {
			return nil, err
		}
		if len(*poolInfo) == 0 {
			return nil, err
		}

		var raydiumMsg model.RaydiumCreateMessage

		for _, pool := range *poolInfo {
			if pool.Mint == tokenAddress {
				raydiumMsg.PoolAddress = pool.Data.PoolAddress
				raydiumMsg.MarketAddress = pool.Data.ReturnPoolData.MarketId
				raydiumMsg.PoolState = 0
				raydiumMsg.PoolBaseReserve = pool.Data.ReturnPoolData.BaseReserve
				raydiumMsg.PoolQuoteReserve = pool.Data.ReturnPoolData.QuoteReserve
				raydiumMsg.BaseToken = pool.Data.ReturnPoolData.BaseMint
				raydiumMsg.QuoteToken = pool.Data.ReturnPoolData.QuoteMint
				raydiumMsg.User = pool.Data.ReturnPoolData.OpenOrders
				// raydiumMsg.Timestamp = pool.Data.ReturnPoolData.PoolOpenTime
				raydiumMsg.Block = uint64(pool.Data.ReturnPoolData.OrderbookToInitTime)
			}
		}
		var service TokenLiquidityPoolService
		pool = service.ConvertMessageToLiquidityPool(&raydiumMsg)
		if err := model.UpsertTokenLiquidityPool(pool); err != nil {
			return nil, err
		}
	}

	return pool, nil
}

// GetTokenLiquidityPoolsByTokenAddresses 通过代币地址和交易平台类型获取代币流动性池记录
func (service *TokenLiquidityPoolService) GetTokenLiquidityPoolsByTokenAddresses(tokenAddresses []string, platformType uint8) ([]model.TokenLiquidityPool, error) {
	return model.GetTokenLiquidityPoolsByTokenAddresses(tokenAddresses, platformType)
}
