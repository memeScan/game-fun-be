package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"game-fun-be/internal/conf"
	"game-fun-be/internal/constants"
	"game-fun-be/internal/es"
	"game-fun-be/internal/model"
	"game-fun-be/internal/pkg/httpRespone"
	"game-fun-be/internal/pkg/httpUtil"
	"game-fun-be/internal/pkg/util"
	"game-fun-be/internal/redis"
	"game-fun-be/internal/service"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"sync/atomic"

	"github.com/IBM/sarama"
	"github.com/shopspring/decimal"
)

var (
	solPriceReady    atomic.Bool
	solPriceInitOnce sync.Once
)

func ConsumePumpfunTopics() error {
	// 确保首次获取SOL价格
	if err := initSolPrice(); err != nil {
		return fmt.Errorf("failed to initialize SOL price: %v", err)
	}

	// 初始化 TopicConsumer

	topicConsumer, err := NewTopicConsumer(KafkaGroupDexProcessor)
	if err != nil {
		util.Log().Error("Error creating TopicConsumer: %s", err)
		return err
	}

	util.Log().Info("Kafka consumer initialized successfully")

	// 从环境变量获取批处理参数
	batchSize, _ := strconv.Atoi(os.Getenv("KAFKA_BATCH_SIZE"))
	if batchSize == 0 {
		batchSize = 100 // 默认值
	}

	batchTimeout, _ := time.ParseDuration(os.Getenv("KAFKA_BATCH_TIMEOUT"))
	if batchTimeout == 0 {
		batchTimeout = 100 * time.Millisecond // 默认值
	}

	// Log current configuration values
	util.Log().Info("Kafka Configuration:\n"+
		"Batch Size:      %d\n"+
		"Batch Timeout:   %v\n"+
		"Brokers:        %s\n"+
		"Client ID:      %s\n"+
		"Group ID:       %s",
		batchSize,
		batchTimeout,
		os.Getenv("KAFKA_BROKERS"),
		os.Getenv("KAFKA_CLIENT_ID"),
		KafkaGroupDexProcessor)

	// 添加 Pumpfun 相关的处理器
	topicConsumer.AddHandler(TopicPumpCreate, PumpfunImmediateHandler)
	topicConsumer.AddHandler(TopicPumpComplete, PumpfunImmediateHandler)
	topicConsumer.AddHandler(TopicPumpSetParams, PumpfunImmediateHandler)
	topicConsumer.AddBatchHandler(TopicPumpTrade, PumpfunBatchHandler, batchSize, 30, batchTimeout)

	// 添加 Raydium 相关的处理器
	topicConsumer.AddHandler(TopicRayCreate, RaydiumImmediateHandler)
	topicConsumer.AddHandler(TopicRayAddLiquidity, RaydiumImmediateHandler)
	topicConsumer.AddHandler(TopicRayRemoveLiquidity, RaydiumImmediateHandler)
	topicConsumer.AddBatchHandler(TopicRaySwap, RaydiumBatchHandler, batchSize, 30, batchTimeout)

	// 添加未知代币处理器
	topicConsumer.AddHandler(TopicUnknownToken, UnknownTokenHandler)

	// 添加game相关处理器
	topicConsumer.AddHandler(TopicGameOutTrade, gameOutTradeHandler)
	topicConsumer.AddHandler(TopicGameInTrade, gameInTradeHandler)

	// 开始消费主题
	return topicConsumer.ConsumeTopics(AllTopics)
}

func initSolPrice() error {
	var initErr error
	solPriceInitOnce.Do(func() {
		ctx := context.Background()
		initErr = service.FetchAndStoreSolPrice(ctx)
		if initErr == nil {
			solPriceReady.Store(true)
		}
	})
	return initErr
}

func PumpfunImmediateHandler(message []byte, topic string) error {
	util.Log().Info("PumpfunImmediateHandler: Processing message from topic %s", topic)

	switch topic {
	case TopicPumpCreate:
		return handlePumpfunCreate(message)
	case TopicPumpComplete:
		return handlePumpfunComplete(message)
	case TopicPumpSetParams:
		return handlePumpfunSetParams(message)
	default:
		util.Log().Warning("Unknown topic: %s", topic)
		return nil
	}
}

func handlePumpfunCreate(message []byte) error {
	var createMsg model.TokenInfoMessage
	if err := json.Unmarshal(message, &createMsg); err != nil {
		util.Log().Error("Failed to unmarshal pumpfun-create message: %v", err)
		return fmt.Errorf("failed to unmarshal pumpfun-create message: %v", err)
	}

	tokenInfoService := &service.TokenInfoService{}

	// 将 TokenInfoMessage 转换为 TokenInfo 并创建记录
	resp := tokenInfoService.CreateTokenInfoFromMessage(&createMsg)
	if resp.Code != 0 {
		util.Log().Error("创建代币信息失败: %v", resp.Error)
		return fmt.Errorf("failed to create token info: %v", resp.Error)
	}

	util.Log().Info("成功创建代币信息: %s", createMsg.Mint)

	//创建代币流动性池
	liquidityPoolService := &service.TokenLiquidityPoolService{}
	liquidityPool := liquidityPoolService.CreatePumpFunInitialPool(&createMsg)
	liquidityPool.RealNativeReserves, _ = strconv.ParseUint(model.PUMP_INITIAL_REALSOL_TOKEN_RESERVES, 10, 64)
	liquidityPool.RealTokenReserves, _ = strconv.ParseUint(model.PUMP_INITIAL_REAL_TOKEN_RESERVES, 10, 64)
	resp = liquidityPoolService.ProcessTokenLiquidityPoolCreation(liquidityPool)
	if resp.Code != 0 {
		util.Log().Error("创建流动性池失败: %v", resp.Error)
	}

	return nil
}

func handlePumpfunComplete(message []byte) error {
	// 解析消息
	var completeMsg model.PumpFuncCompleteMessage
	if err := json.Unmarshal(message, &completeMsg); err != nil {
		return fmt.Errorf("failed to unmarshal pumpfun-complete message: %v", err)
	}
	// 更新代币信息
	tokenInfoService := &service.TokenInfoService{}
	// 获取代币信息
	tokenInfo, err := tokenInfoService.GetTokenInfoByAddress(completeMsg.Mint, uint8(model.ChainTypeSolana))
	if err != nil {
		return fmt.Errorf("failed to get token info: %v", err)
	}
	if tokenInfo == nil {
		return fmt.Errorf("token not found: %s", completeMsg.Mint)
	}

	// 更新完成状态
	resp := tokenInfoService.UpdateTokenComplete(tokenInfo.TokenAddress, uint8(model.ChainTypeSolana))
	if resp.Code != 200 {
		// 添加更详细的错误信息
		errMsg := "unknown error"
		if resp.Error != "" {
			errMsg = resp.Error
		}
		util.Log().Error("Failed to update token complete status: code=%d, error=%s", resp.Code, errMsg)
		return fmt.Errorf("failed to update token complete status: code=%d, error=%s", resp.Code, errMsg)
	}

	// 清除缓存
	if err := redis.Del(fmt.Sprintf("%s:%s", constants.RedisKeyPrefixTokenInfo, completeMsg.Mint)); err != nil {
		util.Log().Error("Failed to delete Redis cache: %v", err)
		// 注意：这里的错误不返回，因为不是致命错误
	}

	util.Log().Info("Successfully marked token %s as complete", completeMsg.Mint)
	return nil
}

func handlePumpfunSetParams(message []byte) error {
	// TODO: 实现 pumpfun-set-params 的处理逻辑
	var setParamsMsg struct {
		// 定义 pumpfun-set-params 消息的结构
	}
	if err := json.Unmarshal(message, &setParamsMsg); err != nil {
		return fmt.Errorf("failed to unmarshal pumpfun-set-params message: %v", err)
	}

	// 处理 setParamsMsg
	util.Log().Info("Handling pumpfun-set-params message")
	return nil
}

func PumpfunBatchHandler(topic string, messages []sarama.ConsumerMessage, partition int32, goroutineID uint64) error {
	util.Log().Info("Processing batch messages for topic %s: %d messages", topic, len(messages))

	switch topic {
	case TopicPumpTrade:
		return handlePumpTradeMessages(messages)
	default:
		util.Log().Warning("Unknown topic for batch processing: %s", topic)
		return nil
	}
}

func handlePumpTradeMessages(messages []sarama.ConsumerMessage) error {

	var tokenTradeMessages []*model.TokenTradeMessage
	for _, msg := range messages {
		var tokenTradeMsg model.TokenTradeMessage
		err := json.Unmarshal(msg.Value, &tokenTradeMsg)
		if err != nil {
			util.Log().Error("Failed to unmarshal message: %v", err)
			continue
		}
		tokenTradeMessages = append(tokenTradeMessages, &tokenTradeMsg)
	}

	// 如果没有有效的交易消息，直接返回
	if len(tokenTradeMessages) == 0 {
		util.Log().Info("No valid trade messages to process")
		return nil
	}

	// 处理交易消息
	if err := handlePumpfunTokenTransactions(tokenTradeMessages); err != nil {
		util.Log().Error("处理交易消息失败: %v", err)
		return fmt.Errorf("处理交易消息失败: %v", err)
	}

	return nil
}

// 处理交易消息
func handlePumpfunTokenTransactions(tokenTradeMessages []*model.TokenTradeMessage) error {
	tokenTxService := &service.TokenTransactionService{}

	// 转换消息到交易
	tokenTransactions := tokenTxService.ConvertTradeMessagesToTransactions(tokenTradeMessages)

	// 分离新旧数据
	currentDate := time.Now().Format("2006-01-02")
	var todayTxs, oldTxs []*model.TokenTransaction

	for _, tx := range tokenTransactions {
		txDate := tx.TransactionTime.Format("2006-01-02")
		if txDate == currentDate {
			todayTxs = append(todayTxs, tx)
		} else {
			oldTxs = append(oldTxs, tx)
		}
	}

	// 处理今的数据
	if len(todayTxs) > 0 {
		if err := processCurrentDayPumpfunTransactions(todayTxs); err != nil {
			return fmt.Errorf("处理当天交易失败: %v", err)
		}
	}

	// 处理历史数据
	if len(oldTxs) > 0 {
		if err := processHistoricalPumpfunTransactions(oldTxs); err != nil {
			return fmt.Errorf("处理历史交易失败: %v", err)
		}
	}

	return nil
}

// 处理当天交易数据
func processCurrentDayPumpfunTransactions(transactions []*model.TokenTransaction) error {
	tokenTxService := &service.TokenTransactionService{}
	tokenTxIndexService := &service.TokenTxIndexService{}

	// 批量创建交易记录
	start := time.Now()
	currentDate := time.Now().Format("20060102")
	resp := tokenTxService.ProcessBatchTokenTransactionCreation(transactions, currentDate)
	if resp.Code != 0 {
		util.Log().Error("批量创建代币交易失败: %v", resp.Error)
		return fmt.Errorf("failed to create transactions: %v", resp.Error)
	}
	util.Log().Info("批量创建交易记录耗时: %v", time.Since(start))

	// 批量创建交易索引
	start = time.Now()
	indexResp := tokenTxIndexService.BatchCreateIndexFromTransactions(transactions)
	if indexResp.Code != 0 {
		util.Log().Error("批量创建代币交易索引失败: %v", indexResp.Error)
		return fmt.Errorf("failed to create indices: %v", indexResp.Error)
	}
	util.Log().Info("批量创建交易索引耗时: %v", time.Since(start))

	// 更新代币信息
	start = time.Now()
	tokenInfoMap, err := updateTokensInfo(transactions)
	if err != nil {
		util.Log().Error("更新代币信息失败: %v", err)
		// 不返回错误，继续处理其他逻辑
	}
	util.Log().Info("更新代币信息耗时: %v", time.Since(start))

	// 更新池子信息
	start = time.Now()
	poolInfoMap, err := updatePoolsInfo(transactions)
	if err != nil {
		util.Log().Error("更新池子信息失败: %v", err)
		// 不返回错误，继续处理其他逻辑
	}
	util.Log().Info("更新池子信息耗时: %v", time.Since(start))

	// 处理 Elasticsearch 数据
	start = time.Now()
	if err := processElasticsearchData(transactions, tokenInfoMap, poolInfoMap); err != nil {
		return err
	}
	util.Log().Info("处理 Elasticsearch 数据耗时: %v", time.Since(start))

	// 处理 ClickHouse 数据
	start = time.Now()
	processClickHouseData(transactions)
	util.Log().Info("处理 ClickHouse 数据耗时: %v", time.Since(start))

	return nil
}

// updateTokensInfo 更新代币信息
func updateTokensInfo(transactions []*model.TokenTransaction) (map[string]*model.TokenInfo, error) {
	tokenInfoService := &service.TokenInfoService{}

	// 按代币地址分组，只保留最新交易
	tokenLatestTx := make(map[string]*model.TokenTransaction)
	for _, tx := range transactions {
		if existing, ok := tokenLatestTx[tx.TokenAddress]; !ok ||
			tx.TransactionTime.After(existing.TransactionTime) {
			tokenLatestTx[tx.TokenAddress] = tx
		}
	}
	// 批量构建代币信息
	var tokenInfos []*model.TokenInfo
	// 收集所有需要查询的地址
	addresses := make([]string, 0, len(tokenLatestTx))
	for _, tx := range tokenLatestTx {
		addresses = append(addresses, tx.TokenAddress)
	}
	// 批量查询代币信息
	tokenInfoMap := tokenInfoService.GetExistingTokenInfos(addresses, uint8(model.ChainTypeSolana))

	// 找出查询到的代币
	missingTokens := make([]string, 0, len(addresses)-len(tokenInfoMap))
	for _, addr := range addresses {
		if tokenInfoMap[addr] == nil {
			missingTokens = append(missingTokens, addr)
		}
	}
	if len(missingTokens) > 0 {
		util.Log().Info("Found missing tokens: %v", missingTokens)

		// 准备批量消息
		messages := make([]*sarama.ProducerMessage, 0, len(missingTokens))

		for _, addr := range missingTokens {
			msg := &sarama.ProducerMessage{
				Topic: TopicUnknownToken,
				Value: sarama.StringEncoder(addr),
			}
			messages = append(messages, msg)
		}

		// 批量发送消息
		if len(messages) > 0 {
			err := GetProducer().SendMessages(messages)
			if err != nil {
				util.Log().Error("Failed to send batch unknown token messages: %v", err)
			} else {
				util.Log().Info("Successfully sent %d unknown token messages", len(messages))
			}
		}
	}
	// 找出未检查权限的
	uncheckedTokens := make([]string, 0, len(tokenInfoMap))
	for addr, tokenInfo := range tokenInfoMap {
		if tokenInfo != nil &&
			tokenInfo.CreatedPlatformType != uint8(model.CreatedPlatformTypePump) &&
			!tokenInfo.HasFlag(model.FLAG_AUTHORITY_CHECKED) {
			uncheckedTokens = append(uncheckedTokens, addr)
		}
	}

	// 批量查询权限数据
	var authorityMap map[string]*httpRespone.SafetyData
	if len(uncheckedTokens) > 0 {
		safetyResp, err := httpUtil.GetSafetyCheckData(uncheckedTokens)
		if err != nil {
			// 检查是否是 429 错误
			if strings.Contains(err.Error(), "status code 429") {
				util.Log().Error("Rate limit exceeded for safety check API: %v", err)
			}
			util.Log().Error("Failed to get safety check data: %v", err)
		} else if safetyResp != nil {
			authorityMap = make(map[string]*httpRespone.SafetyData, len(*safetyResp))
			for i := range *safetyResp {
				authorityMap[(*safetyResp)[i].Mint] = &(*safetyResp)[i]
			}
		}
	}

	// 处理创建者交易
	for _, tx := range transactions {
		tokenInfo := tokenInfoMap[tx.TokenAddress]
		if tokenInfo == nil || tokenInfo.Creator != tx.UserAddress {
			continue // 跳过非创建者交易
		}

		// 处理创建者交易
		processCreatorTransaction(tokenInfo, tx)
	}

	for _, tx := range tokenLatestTx {
		tokenInfo, exists := tokenInfoMap[tx.TokenAddress]
		if !exists {
			continue // 跳过这个代币的处理
		}
		tokenInfos = append(tokenInfos, tokenInfo)

		// 计算市值需要考虑代币度
		actualSupply := decimal.NewFromInt(int64(tokenInfo.TotalSupply)).Shift(-int32(tokenInfo.Decimals))

		marketCap := tx.Price.Mul(actualSupply)
		tokenInfo.MarketCap = marketCap

		// 设置流通市值
		actualCirculatingSupply := decimal.NewFromInt(int64(tokenInfo.CirculatingSupply)).Shift(-int32(tokenInfo.Decimals))
		circulatingMarketCap := tx.Price.Mul(actualCirculatingSupply)
		tokenInfo.CirculatingMarketCap = circulatingMarketCap

		tokenInfo.Progress = tx.Progress
		// 如果 PoolAddress 为空，补充池子地址
		if tokenInfo.PoolAddress == "" {
			tokenInfo.PoolAddress = tx.PoolAddress
		}

		// 检查进度并设置皇冠时间
		if tx.Progress.GreaterThanOrEqual(decimal.NewFromInt(42)) && tokenInfo.RocketDuration == 0 {
			tokenInfo.RocketDuration = time.Now().Unix() - tokenInfo.TransactionTime.Unix()
		}
		if tx.Progress.GreaterThanOrEqual(decimal.NewFromInt(84)) && tokenInfo.CrownDuration == 0 {
			tokenInfo.CrownDuration = time.Now().Unix() - tokenInfo.TransactionTime.Unix()
		}
		if tx.PlatformType == uint8(model.PlatformTypeRaydium) {
			tokenInfo.IsComplete = true
		}

		// 设置价格和 sol 位格
		tokenInfo.Price = tx.Price
		tokenInfo.NativePrice = tx.NativePrice

		// 计算流动性
		if tx.PlatformType == uint8(model.PlatformTypePump) {
			// 内盘时使用实际 SOL 储备计算流动性（美元价值）
			if tx.RealNativeReserves != 0 {
				// 使用实际 SOL 储备计算
				adjustedNativeReserves := decimal.NewFromInt(int64(tx.RealNativeReserves)).Shift(-9)
				// 计算 SOL 的美元价值
				solValue := adjustedNativeReserves.Mul(tx.NativePriceUSD)
				// 由于是双边流动性，总流动性为 SOL 价值的 2 倍
				tokenInfo.Liquidity = solValue.Mul(decimal.NewFromInt(2))
			}
		} else if tx.PlatformType == uint8(model.PlatformTypeRaydium) {
			// Raydium外盘流动性计算:
			// 1. 将池子中的SOL储备转换为USD价 (SOL储备量 * SOL格)
			// 2. 由于池子是双边的，总流动性为SOL价值的2倍
			nativeReserves := decimal.NewFromInt(int64(tx.VirtualNativeReserves)).Shift(-9)
			if !nativeReserves.IsZero() && !tx.NativePriceUSD.IsZero() {
				solValue := nativeReserves.Mul(tx.NativePriceUSD)
				tokenInfo.Liquidity = solValue.Mul(decimal.NewFromInt(2))
			}
		}

		// 检查权限
		if data, ok := authorityMap[tokenInfo.TokenAddress]; ok {
			if data.MintAuthority == 1 {
				tokenInfo.SetFlag(model.FLAG_MINT_AUTHORITY)
			}
			if data.FreezeAuthority == 1 {
				tokenInfo.SetFlag(model.FLAG_FREEZE_AUTHORITY)
			}
			tokenInfo.SetFlag(model.FLAG_AUTHORITY_CHECKED)
		}

		tokenInfo.UpdateTime = time.Now()
	}

	// 如果没有需要更新的数据直接返回
	if len(tokenInfos) == 0 {
		util.Log().Info("No token info needs to be updated")
		return tokenInfoMap, nil
	}

	// 批量更新数据库
	resp := tokenInfoService.BatchUpdateTokenInfo(tokenInfos)
	if resp.Code != 0 {
		// 数据库更新失败时，清除相关的 Redis 缓存
		keys := make([]string, 0, len(tokenInfos))
		for _, info := range tokenInfos {
			keys = append(keys, fmt.Sprintf("%s:%s", constants.RedisKeyPrefixTokenInfo, info.TokenAddress))
		}
		if err := redis.Del(keys...); err != nil {
			util.Log().Error("Failed to batch delete Redis cache: %v", err)
		}
		return tokenInfoMap, fmt.Errorf("批量更新代币信息失败: %v", resp.Error)
	}

	// 批量更新缓存
	keyValues := make(map[string]string, len(tokenInfos))
	for _, info := range tokenInfos {
		infoJSON, err := json.Marshal(info)
		if err != nil {
			util.Log().Error("Failed to marshal token info: %v", err)
			continue
		}
		redisKey := fmt.Sprintf("%s:%s", constants.RedisKeyPrefixTokenInfo, info.TokenAddress)
		keyValues[redisKey] = string(infoJSON)
	}

	if len(keyValues) > 0 {
		if err := redis.MSet(keyValues, 1*time.Minute); err != nil {
			util.Log().Error("Failed to batch update Redis cache: %v", err)
		}
	}
	util.Log().Info("成功更新 %d 个代币的信息", len(tokenInfos))
	return tokenInfoMap, nil
}

// processCreatorTransaction 处理创建者交易，更新代币状态
func processCreatorTransaction(tokenInfo *model.TokenInfo, tx *model.TokenTransaction) {
	// 记录更新前的余额，用于后续判断状态
	oldBalance := tokenInfo.DevTokenAmount

	// 更新开发者持有的代币数量
	switch tx.TransactionType {
	case uint8(model.TransactionTypeSell):
		// 卖出：减少开发者余额
		tokenInfo.DevTokenAmount -= tx.TokenAmount
		// 更新状态
		if tokenInfo.DevTokenAmount == 0 {
			tokenInfo.DevStatus = uint8(model.DevStatusClear) // 清仓
		} else {
			tokenInfo.DevStatus = uint8(model.DevStatusSell) // 卖出
		}

	case uint8(model.TransactionTypeBuy):
		// 买入：增加开发者余额
		tokenInfo.DevTokenAmount += tx.TokenAmount
		// 更新状态
		if oldBalance == 0 {
			tokenInfo.DevStatus = uint8(model.DevStatusHold) // 持有
		} else {
			tokenInfo.DevStatus = uint8(model.DevStatusIncrease) // 加仓
		}
	}

}

// 处理历史交易数据
func processHistoricalPumpfunTransactions(transactions []*model.TokenTransaction) error {
	tokenTxService := &service.TokenTransactionService{}
	tokenTxIndexService := &service.TokenTxIndexService{}

	// 按日期分组
	txsByDate := make(map[string][]*model.TokenTransaction)
	for _, tx := range transactions {
		date := tx.TransactionTime.Format("20060102")
		txsByDate[date] = append(txsByDate[date], tx)
	}

	// 按日期处理历史数据
	for date, txs := range txsByDate {
		util.Log().Info("处理 %s 的历史数据，共 %d 条", date, len(txs))

		// 检查日期表是否存在，不存在则创建
		if err := model.CreateTableForDate(date); err != nil {
			util.Log().Error("确保表存在失败: %v", err)
			continue
		}

		// 批量创建历史交易记录
		resp := tokenTxService.ProcessBatchTokenTransactionCreation(txs, date)
		if resp.Code != 0 {
			util.Log().Error("创建 %s 的历史交易失败: %v", date, resp.Error)
			continue
		}

		// 批量创建历史交易索引
		indexResp := tokenTxIndexService.BatchCreateIndexFromTransactions(txs)
		if indexResp.Code != 0 {
			util.Log().Error("创建 %s 的历史交易索引失败: %v", date, indexResp.Error)
			continue
		}

		util.Log().Info("成功处理 %s 的历史数据", date)
	}

	return nil
}

func RaydiumBatchHandler(topic string, messages []sarama.ConsumerMessage, partition int32, goroutineID uint64) error {
	util.Log().Info("Processing batch messages for topic %s: %d messages", topic, len(messages))

	switch topic {
	case TopicRaySwap:
		return handleRaydiumSwapMessages(messages)
	default:
		util.Log().Warning("Unknown topic for batch processing: %s", topic)
		return nil
	}
}

func handleRaydiumSwapMessages(messages []sarama.ConsumerMessage) error {
	var raydiumSwapMessages []*model.RaydiumSwapMessage
	for _, msg := range messages {
		var swapMsg model.RaydiumSwapMessage
		err := json.Unmarshal(msg.Value, &swapMsg)
		if err != nil {
			util.Log().Error("Failed to unmarshal Raydium swap message: %v", err)
			continue
		}
		raydiumSwapMessages = append(raydiumSwapMessages, &swapMsg)
	}

	// 如果没有有效的交易消息,直接返回
	if len(raydiumSwapMessages) == 0 {
		util.Log().Info("No valid Raydium swap messages to process")
		return nil
	}

	// 处理 Raydium 交易消息
	if err := handleRaydiumSwapTransactions(raydiumSwapMessages); err != nil {
		util.Log().Error("处理 Raydium 交易消息失败: %v", err)
		return fmt.Errorf("处理 Raydium 交易消息失败: %v", err)
	}

	return nil
}

// 处理 Raydium 交易消息
func handleRaydiumSwapTransactions(swapMessages []*model.RaydiumSwapMessage) error {
	tokenTxService := &service.TokenTransactionService{}

	// 换消息到交易
	tokenTransactions := tokenTxService.ConvertRaydiumSwapMessagesToTransactions(swapMessages)

	// 分离新旧数据
	currentDate := time.Now().Format("2006-01-02")
	var todayTxs, oldTxs []*model.TokenTransaction

	for _, tx := range tokenTransactions {
		txDate := tx.TransactionTime.Format("2006-01-02")
		if txDate == currentDate {
			todayTxs = append(todayTxs, tx)
		} else {
			oldTxs = append(oldTxs, tx)
		}
	}

	// 处理今天的数据
	if len(todayTxs) > 0 {
		if err := processCurrentDayRaydiumTransactions(todayTxs); err != nil {
			return fmt.Errorf("处理当天 Raydium 交易失败: %v", err)
		}
	}

	// 处理历史数据
	if len(oldTxs) > 0 {
		if err := processHistoricalRaydiumTransactions(oldTxs); err != nil {
			return fmt.Errorf("处理历史 Raydium 交易失败: %v", err)
		}
	}

	return nil
}

// 处理当天 Raydium 交易数据
func processCurrentDayRaydiumTransactions(transactions []*model.TokenTransaction) error {
	totalStart := time.Now()
	tokenTxService := &service.TokenTransactionService{}
	tokenTxIndexService := &service.TokenTxIndexService{}

	// 1. 批量创建交易记录
	txCreateStart := time.Now()
	currentDate := time.Now().Format("20060102")
	resp := tokenTxService.ProcessBatchTokenTransactionCreation(transactions, currentDate)
	if resp.Code != 0 {
		util.Log().Error("批量创建代币交易失败: %v", resp.Error)
		return fmt.Errorf("failed to create transactions: %v", resp.Error)
	}
	util.Log().Info("1. 批量创建交易记录耗时: %v, 交易数量: %d",
		time.Since(txCreateStart),
		len(transactions))

	// 2. 批量创建交易索引
	indexStart := time.Now()
	indexResp := tokenTxIndexService.BatchCreateIndexFromTransactions(transactions)
	if indexResp.Code != 0 {
		util.Log().Error("批量创建代币交易索引失败: %v", indexResp.Error)
		return fmt.Errorf("failed to create indices: %v", indexResp.Error)
	}
	util.Log().Info("2. 批量创建交易索引耗时: %v", time.Since(indexStart))

	// 3. 更新代币信息
	tokenStart := time.Now()
	tokenInfoMap, err := updateTokensInfo(transactions)
	if err != nil {
		util.Log().Error("更新代币信息失败: %v", err)
		// 不返回错误，继续处理其他逻辑
	}
	util.Log().Info("3. 更新代币信息耗时: %v, 更新代币数量: %d",
		time.Since(tokenStart),
		len(tokenInfoMap))

	// 4. 更新池子信息
	poolStart := time.Now()
	poolInfoMap, err := updatePoolsInfo(transactions)
	if err != nil {
		util.Log().Error("更新池子信息失败: %v", err)
		// 不返回错误，继续处理其他逻辑
	}
	util.Log().Info("4. 更新池子信息耗时: %v, 更新池子数量: %d",
		time.Since(poolStart),
		len(poolInfoMap))

	// 5. 处理 Elasticsearch 数据
	esStart := time.Now()
	if err := processElasticsearchData(transactions, tokenInfoMap, poolInfoMap); err != nil {
		return err
	}
	util.Log().Info("5. 处理 Elasticsearch 数据耗时: %v", time.Since(esStart))

	// 6. 处理 ClickHouse 数据
	chStart := time.Now()
	processClickHouseData(transactions)
	util.Log().Info("6. 处理 ClickHouse 数据耗时: %v", time.Since(chStart))

	// 总耗时统计
	totalTime := time.Since(totalStart)
	util.Log().Info("Raydium交易处理总耗时: %v, 处理交易数: %d (平均: %v/笔)",
		totalTime,
		len(transactions),
		totalTime/time.Duration(len(transactions)))

	return nil
}

// 处理历史 Raydium 交易数据
func processHistoricalRaydiumTransactions(transactions []*model.TokenTransaction) error {
	tokenTxService := &service.TokenTransactionService{}
	tokenTxIndexService := &service.TokenTxIndexService{}

	// 按日期分组
	txsByDate := make(map[string][]*model.TokenTransaction)
	for _, tx := range transactions {
		date := time.Unix(tx.TransactionTime.Unix(), 0).Format("20060102")
		txsByDate[date] = append(txsByDate[date], tx)
	}

	// 按日期处理历史数据
	for date, txs := range txsByDate {
		util.Log().Info("处理 %s 的 Raydium 历史数据,共 %d 条", date, len(txs))

		// 检查日期表是否存在,不存在则创建
		if err := model.CreateTableForDate(date); err != nil {
			util.Log().Error("确保表存在失败: %v", err)
			continue
		}

		// 批量创建历史交易记录
		resp := tokenTxService.ProcessBatchTokenTransactionCreation(txs, date)
		if resp.Code != 0 {
			util.Log().Error("创建 %s 的历史交易失败: %v", date, resp.Error)
			continue
		}

		// 批量创建历史交易索引
		indexResp := tokenTxIndexService.BatchCreateIndexFromTransactions(txs)
		if indexResp.Code != 0 {
			util.Log().Error("创建 %s 的历史交易索引失败: %v", date, indexResp.Error)
			continue
		}

		util.Log().Info("成功处理 %s 的 Raydium 历史数据", date)
	}

	return nil
}

func RaydiumImmediateHandler(message []byte, topic string) error {
	util.Log().Info("RaydiumImmediateHandler: Processing message from topic %s", topic)

	switch topic {
	case TopicRayCreate:
		return handleRaydiumCreate(message)
	case TopicRayAddLiquidity:
		return handleRaydiumAddLiquidity(message)
	case TopicRayRemoveLiquidity:
		return handleRaydiumRemoveLiquidity(message)
	default:
		util.Log().Warning("Unknown topic: %s", topic)
		return nil
	}
}

func handleRaydiumCreate(message []byte) error {
	var createMsg model.RaydiumCreateMessage
	if err := json.Unmarshal(message, &createMsg); err != nil {
		util.Log().Error("Failed to unmarshal raydium-create message: %v", err)
		return fmt.Errorf("failed to unmarshal raydium-create message: %v", err)
	}

	// 创建流动性池记录
	liquidityPoolService := &service.TokenLiquidityPoolService{}
	resp := liquidityPoolService.CreateLiquidityPoolFromMessage(&createMsg)
	if resp.Code != 0 {
		util.Log().Error("创建流动性池失败: %v", resp.Error)
		return fmt.Errorf("failed to create liquidity pool: %v", resp.Error)
	}

	util.Log().Info("成功创建流动性池: %s", createMsg.PoolAddress)

	return nil
}

func handleRaydiumAddLiquidity(message []byte) error {

	var createMsg model.RaydiumCreateMessage
	if err := json.Unmarshal(message, &createMsg); err != nil {
		util.Log().Error("Failed to unmarshal raydium-create message: %v", err)
		return fmt.Errorf("failed to unmarshal raydium-create message: %v", err)
	}

	// 更新流动性池记录
	liquidityPoolService := &service.TokenLiquidityPoolService{}
	pool := new(model.TokenLiquidityPool)
	pool.PoolAddress = createMsg.PoolAddress
	pool.ChainType = uint8(model.ChainTypeSolana)
	baseReserve, _ := strconv.ParseUint(createMsg.PoolBaseReserve, 10, 64)
	quoteReserve, _ := strconv.ParseUint(createMsg.PoolQuoteReserve, 10, 64)
	pool.PoolPcReserve = baseReserve
	pool.PoolCoinReserve = quoteReserve
	pool.UpdateTime = time.Now()
	resp := liquidityPoolService.ProcessTokenLiquidityPoolUpdate(pool)
	if resp.Code != 0 {
		util.Log().Error("更新流动性池失败: %v", resp.Error)
		return fmt.Errorf("failed to update liquidity pool: %v", resp.Error)
	}
	util.Log().Info("成功更新流动性池: %s", createMsg.PoolAddress)
	return nil
}

func handleRaydiumRemoveLiquidity(message []byte) error {

	var createMsg model.RaydiumCreateMessage
	if err := json.Unmarshal(message, &createMsg); err != nil {
		util.Log().Error("Failed to unmarshal raydium-create message: %v", err)
		return fmt.Errorf("failed to unmarshal raydium-create message: %v", err)
	}

	// 更新流动性池记录
	liquidityPoolService := &service.TokenLiquidityPoolService{}
	pool := new(model.TokenLiquidityPool)
	pool.PoolAddress = createMsg.PoolAddress
	pool.ChainType = uint8(model.ChainTypeSolana)
	baseReserve, _ := strconv.ParseUint(createMsg.PoolBaseReserve, 10, 64)
	quoteReserve, _ := strconv.ParseUint(createMsg.PoolQuoteReserve, 10, 64)
	pool.PoolPcReserve = baseReserve
	pool.PoolCoinReserve = quoteReserve
	pool.UpdateTime = time.Now()
	resp := liquidityPoolService.ProcessTokenLiquidityPoolUpdate(pool)
	if resp.Code != 0 {
		util.Log().Error("更新流动性池失败: %v", resp.Error)
		return fmt.Errorf("failed to update liquidity pool: %v", resp.Error)
	}
	util.Log().Info("成功更新流动性池: %s", createMsg.PoolAddress)
	return nil
}

// 处理 ClickHouse 数据
func processClickHouseData(transactions []*model.TokenTransaction) {
	transactionCkService := &service.TransactionCkServiceImpl{}

	// 转换并处理交易数据
	transactionCks := transactionCkService.ConvertToTransactionCks(transactions)
	if err := transactionCkService.BatchProcessTransactions(transactionCks); err != nil {
		util.Log().Error("Failed to process transactions in ClickHouse: %v", err)
		// 这里只记录错误但不返回，因为 ClickHouse 的错误不应影响主流
	}
}

// 处理 Elasticsearch 数据
func processElasticsearchData(transactions []*model.TokenTransaction, tokenInfoMap map[string]*model.TokenInfo, poolInfoMap map[string]*model.TokenLiquidityPool) error {
	totalStart := time.Now()
	tokenTxService := &service.TokenTransactionService{}

	// 1. 获取 ES 文档列表
	prepareStart := time.Now()
	esDocList := tokenTxService.GetESDocList(transactions, tokenInfoMap, poolInfoMap)
	prepareTime := time.Since(prepareStart)
	util.Log().Info("ES文档准备耗时: %v, 文档数量: %d", prepareTime, len(esDocList))

	// 2. 批量索引文档
	indexStart := time.Now()
	resp, err := es.BulkIndexDocuments(es.ES_INDEX_TOKEN_TRANSACTIONS_ALIAS, esDocList)
	if err != nil {
		util.Log().Error("Failed to bulk index documents in Elasticsearch: %v", err)
		return fmt.Errorf("failed to index to ES: %v", err)
	}
	indexTime := time.Since(indexStart)

	// 3. 记录批量索引结果
	if resp != nil {
		util.Log().Info("ES批量索引结果 - 成功: %d, 耗时: %v",
			len(resp.Items),
			indexTime)

		// 如果有失败的文档，记录详细信息
		if len(resp.Failed()) > 0 {
			for _, item := range resp.Failed() {
				util.Log().Error("ES索引失败 - Index: %s, Type: %s, ID: %s, Error: %v",
					item.Index,
					item.Type,
					item.Id,
					item.Error)
			}
		}
	}

	totalTime := time.Since(totalStart)
	util.Log().Info("ES处理总耗时: %v (准备: %v, 索引: %v), 总文档数: %d",
		totalTime,
		prepareTime,
		indexTime,
		len(esDocList))

	return nil
}

// updatePoolsInfo 更新代币池子信息
func updatePoolsInfo(transactions []*model.TokenTransaction) (map[string]*model.TokenLiquidityPool, error) {
	liquidityPoolService := &service.TokenLiquidityPoolService{}

	// 按池子地址分组，只保留最新交易
	poolLatestTx := make(map[string]*model.TokenTransaction)
	for _, tx := range transactions {
		if existing, ok := poolLatestTx[tx.PoolAddress]; !ok ||
			tx.TransactionTime.After(existing.TransactionTime) {
			poolLatestTx[tx.PoolAddress] = tx
		}
	}

	// 收集所有要查询的池子地址
	addresses := make([]string, 0, len(poolLatestTx))
	for _, tx := range poolLatestTx {
		addresses = append(addresses, tx.PoolAddress)
	}

	// 批量查询池子信息
	poolInfoMap := liquidityPoolService.GetExistingPools(addresses)

	// 找出未查询到的池子
	missingPools := make([]string, 0, len(addresses)-len(poolInfoMap))
	for _, addr := range addresses {
		if poolInfoMap[addr] == nil {
			missingPools = append(missingPools, addr)
		}
	}
	if len(missingPools) > 0 {
		util.Log().Info("Found missing pools: %v", missingPools)
	}

	// 迭代交易更新池子信息
	var poolsToUpdate []*model.TokenLiquidityPool
	for poolAddress, tx := range poolLatestTx {
		pool, exists := poolInfoMap[poolAddress]
		if !exists {
			continue
		}

		if tx.MarketAddress != "" && pool.MarketAddress == "" {
			pool.MarketAddress = tx.MarketAddress
		}
		// 更新池子储备量
		pool.PoolPcReserve = tx.VirtualNativeReserves
		pool.PoolCoinReserve = tx.VirtualTokenReserves

		if tx.PlatformType == uint8(model.PlatformTypePump) {
			pool.RealNativeReserves = tx.RealNativeReserves
			pool.RealTokenReserves = tx.RealTokenReserves
		}

		pool.UpdateTime = time.Now()

		poolsToUpdate = append(poolsToUpdate, pool)
	}

	// 批量更新数据库
	if len(poolsToUpdate) > 0 {
		resp := liquidityPoolService.BatchUpdatePools(poolsToUpdate)
		if resp.Code != 0 {
			util.Log().Error("批量更新池子信息失败: %v", resp.Error)
			return poolInfoMap, fmt.Errorf("批量更新池子信息失败: %v", resp.Error)
		}
		util.Log().Info("成功更新 %d 个池子的信息", len(poolsToUpdate))
	}

	return poolInfoMap, nil
}

// 新增处器函数
func UnknownTokenHandler(message []byte, topic string) error {
	var tokenAddress string = string(message)
	// util.Log().Info("Processing unknown token: %s", tokenAddress)

	// 测试环境跳过处理
	if conf.IsTest() {
		// util.Log().Info("Skip processing unknown token in test environment: %s", tokenAddress)
		// return nil
	}

	// 使用 SETNX 进行原子检查和设置
	lockKey := fmt.Sprintf("processing:unknown_token:%s", tokenAddress)
	locked, err := redis.SetNX(lockKey, "1", 30*time.Second)
	if err != nil {
		util.Log().Error("Failed to check/set processing status for token %s: %v", tokenAddress, err)
		return fmt.Errorf("failed to check/set processing status: %v", err)
	}

	// 如果返回 false，说明已经在处理中
	if !locked {
		util.Log().Info("Token %s is already being processed, skipping", tokenAddress)
		return nil // 直接返回成功，消息会被确认
	}

	// 确保处理完成后删除标记
	defer func() {
		if err := redis.Del(lockKey); err != nil {
			util.Log().Error("Failed to remove processing status for token %s: %v", tokenAddress, err)
		}
	}()

	// 先检查数据库是否已存在
	tokenInfoService := &service.TokenInfoService{}
	existingToken, err := tokenInfoService.GetTokenInfoByAddress(tokenAddress, uint8(model.ChainTypeSolana))
	if err != nil {
		util.Log().Error("Failed to check existing token: %v", err)
		return nil
	} else if existingToken != nil {
		util.Log().Info("Token already exists in database: %s", tokenAddress)
		return nil
	}

	// 调用接口获取代币信息
	resp, err := httpUtil.GetTokenFullInfo([]string{tokenAddress}, "sol")
	if err != nil {
		util.Log().Error("Failed to get token info: %v", err)
		return fmt.Errorf("failed to get token info: %v", err)
	}

	if resp == nil || len(resp.Data) == 0 {
		util.Log().Warning("No token info found for address: %s", tokenAddress)
		return nil
	}

	// 转换为 TokenInfo 模型
	tokenInfo := &model.TokenInfo{}

	// 设置基本信息
	tokenInfo.TokenAddress = tokenAddress
	tokenInfo.ChainType = uint8(model.ChainTypeSolana)

	// 设置代币名称和符号
	tokenInfo.TokenName = resp.Data[0].Name
	tokenInfo.Symbol = resp.Data[0].Symbol

	// 设置精度和供应量
	tokenInfo.Decimals = uint8(resp.Data[0].Decimals)
	tokenInfo.TotalSupply, _ = strconv.ParseUint(resp.Data[0].Supply, 10, 64)
	tokenInfo.CirculatingSupply = tokenInfo.TotalSupply
	// 设置创建者信息
	tokenInfo.Creator = resp.Data[0].Creator

	// 设置权限标记
	if resp.Data[0].FreezeAuthority == 1 {
		tokenInfo.SetFlag(model.FLAG_FREEZE_AUTHORITY)
	}
	if resp.Data[0].MintAuthority == 1 {
		tokenInfo.SetFlag(model.FLAG_MINT_AUTHORITY)
	}
	// 标记已检查权限
	tokenInfo.SetFlag(model.FLAG_AUTHORITY_CHECKED)

	// 设置平台类型
	if resp.Data[0].PlatformType == 1 {
		tokenInfo.CreatedPlatformType = uint8(model.CreatedPlatformTypePump)
		tokenInfo.PoolAddress = resp.Data[0].BondingCurveAddress
	} else {
		tokenInfo.CreatedPlatformType = uint8(model.CreatedPlatformTypeUnknown)
		tokenInfo.IsComplete = true
	}

	// 设置时间信息
	tokenInfo.TransactionTime = time.Unix(resp.Data[0].Timestamp, 0)
	tokenInfo.TransactionHash = resp.Data[0].Signature
	tokenInfo.Block = uint64(resp.Data[0].Block)
	tokenInfo.IsMedia = false                         // 初始媒体类型状态
	tokenInfo.DevStatus = uint8(model.DevStatusClear) // 初始开发者状态

	// 设置 URI
	tokenInfo.URI = resp.Data[0].URI
	if resp.Data[0].URI != "" {
		content, err := service.GetURIContent(resp.Data[0].URI, 2) // 设置 2 次重试
		if err != nil {
			util.Log().Error("Error fetching URI content: %v", err)
		} else {
			util.Log().Info("Successfully fetched URI content: %s", content)
			hasSocial := service.HasSocialMedia(content)
			tokenInfo.IsMedia = hasSocial
			tokenInfo.ExtInfo = content
		}
	}
	tokenInfo.UpdateTime = time.Now()

	// 保存到数据库
	serviceResp := tokenInfoService.ProcessTokenInfoCreation(tokenInfo)
	if serviceResp.Code != 0 {
		util.Log().Error("保存代币信息失败: %v", serviceResp.Error)
		// return fmt.Errorf("保存代币信息失败: %v", serviceResp.Error)
	}

	util.Log().Info("Successfully processed unknown token: %s", tokenAddress)

	//补充池子信息
	if tokenInfo.CreatedPlatformType == uint8(model.CreatedPlatformTypePump) {
		// 先检查数据库是否已存在
		liquidityPoolService := &service.TokenLiquidityPoolService{}
		existingPool, err := liquidityPoolService.GetTokenLiquidityPoolsByTokenAddresses([]string{tokenAddress}, uint8(model.PlatformTypePump))
		if err != nil {
			util.Log().Error("Failed to check existing pool: %v", err)
			return nil
		}
		if len(existingPool) > 0 {
			util.Log().Info("Pool already exists in database: %s", tokenAddress)
			return nil
		}
		// 如果池子不存在，从 GetBondingCurves 口查询
		bondingResp, err := httpUtil.GetBondingCurves([]string{tokenInfo.PoolAddress})
		if err != nil {
			util.Log().Error("Failed to get bonding curves info: %v", err)
			return nil
		}

		if bondingResp == nil || len(bondingResp.Data) == 0 {
			util.Log().Warning("No bonding curves data found for address: %s", tokenInfo.PoolAddress)
			return nil
		}

		//创建代币流动性池
		var createMsg model.TokenInfoMessage
		createMsg.Mint = tokenAddress
		createMsg.BondingCurve = tokenInfo.PoolAddress
		createMsg.Creator = tokenInfo.Creator
		createMsg.Timestamp = tokenInfo.TransactionTime.Unix()
		createMsg.Block = tokenInfo.Block
		liquidityPool := liquidityPoolService.CreatePumpFunInitialPool(&createMsg)
		liquidityPool.PoolCoinReserve, _ = strconv.ParseUint(bondingResp.Data[0].Data.VirtualTokenReserves, 10, 64)
		liquidityPool.PoolPcReserve, _ = strconv.ParseUint(bondingResp.Data[0].Data.VirtualSolReserves, 10, 64)
		liquidityPool.RealNativeReserves, _ = strconv.ParseUint(bondingResp.Data[0].Data.RealSolReserves, 10, 64)
		liquidityPool.RealTokenReserves, _ = strconv.ParseUint(bondingResp.Data[0].Data.RealTokenReserves, 10, 64)
		serviceResponse := liquidityPoolService.ProcessTokenLiquidityPoolCreation(liquidityPool)
		if serviceResponse.Code != 0 {
			util.Log().Error("创建流动性池失败: %v", serviceResponse.Error)
			return fmt.Errorf("failed to create liquidity pool: %v", serviceResponse.Error)
		}
		util.Log().Info("Successfully processed unknown pump pool: %s", liquidityPool.PoolAddress)
		saveRaydiumPool(tokenAddress)
		return nil
	} else {
		// 先检查数据库是否已存在
		liquidityPoolService := &service.TokenLiquidityPoolService{}
		existingPool, err := liquidityPoolService.GetTokenLiquidityPoolsByTokenAddresses([]string{tokenAddress}, uint8(model.PlatformTypeRaydium))
		if err != nil {
			util.Log().Error("Failed to check existing pool: %v", err)
			return nil
		}
		if len(existingPool) > 0 {
			util.Log().Info("Pool already exists in database: %s", tokenAddress)
			return nil
		}

		// 调用接口获取池子信息
		resp, err := httpUtil.GetPoolInfo([]string{tokenAddress})
		if err != nil {
			util.Log().Error("Failed to get pool info: %v", err)
			return fmt.Errorf("failed to get pool info: %v", err)
		}
		if resp == nil || len(*resp) == 0 {
			util.Log().Warning("No pool info found for address: %s", tokenAddress)
			return nil
		}

		// 构 RaydiumCreateMessage
		var raydiumMsg model.RaydiumCreateMessage

		raydiumMsg.PoolAddress = (*resp)[0].Data.PoolAddress
		raydiumMsg.MarketAddress = (*resp)[0].Data.ReturnPoolData.MarketId
		raydiumMsg.PoolState = 0
		raydiumMsg.PoolBaseReserve = (*resp)[0].Data.ReturnPoolData.BaseReserve
		raydiumMsg.PoolQuoteReserve = (*resp)[0].Data.ReturnPoolData.QuoteReserve
		raydiumMsg.BaseToken = (*resp)[0].Data.ReturnPoolData.BaseMint
		raydiumMsg.QuoteToken = (*resp)[0].Data.ReturnPoolData.QuoteMint
		raydiumMsg.User = (*resp)[0].Data.ReturnPoolData.OpenOrders
		// raydiumMsg.Timestamp = (*resp)[0].Data.ReturnPoolData.PoolOpenTime

		// 转换并保存到数据库
		pool := liquidityPoolService.ConvertMessageToLiquidityPool(&raydiumMsg)

		serviceResponse := liquidityPoolService.ProcessTokenLiquidityPoolCreation(pool)
		if serviceResponse.Code != 0 {
			util.Log().Error("Failed to save pool info: %v", serviceResponse.Error)
			return fmt.Errorf("failed to save pool info: %v", serviceResponse.Error)
		}

		util.Log().Info("Successfully processed unknown raydium pool: %s", pool.PoolAddress)
		return nil
	}

}

func saveRaydiumPool(tokenAddress string) error {
	// 先检查数据库是否已存在
	liquidityPoolService := &service.TokenLiquidityPoolService{}
	existingPool, err := liquidityPoolService.GetTokenLiquidityPoolsByTokenAddresses([]string{tokenAddress}, uint8(model.PlatformTypeRaydium))
	if err != nil {
		util.Log().Error("Failed to check existing pool: %v", err)
		return nil
	}
	if len(existingPool) > 0 {
		util.Log().Info("Pool already exists in database: %s", tokenAddress)
		return nil
	}

	// 调用接口获取池子信息
	resp, err := httpUtil.GetPoolInfo([]string{tokenAddress})
	if err != nil {
		util.Log().Error("Failed to get pool info: %v", err)
		return fmt.Errorf("failed to get pool info: %v", err)
	}
	if resp == nil || len(*resp) == 0 {
		util.Log().Warning("No pool info found for address: %s", tokenAddress)
		return nil
	}

	// 构 RaydiumCreateMessage
	var raydiumMsg model.RaydiumCreateMessage

	raydiumMsg.PoolAddress = (*resp)[0].Data.PoolAddress
	raydiumMsg.MarketAddress = (*resp)[0].Data.ReturnPoolData.MarketId
	raydiumMsg.PoolState = 0
	raydiumMsg.PoolBaseReserve = (*resp)[0].Data.ReturnPoolData.BaseReserve
	raydiumMsg.PoolQuoteReserve = (*resp)[0].Data.ReturnPoolData.QuoteReserve
	raydiumMsg.BaseToken = (*resp)[0].Data.ReturnPoolData.BaseMint
	raydiumMsg.QuoteToken = (*resp)[0].Data.ReturnPoolData.QuoteMint
	raydiumMsg.User = (*resp)[0].Data.ReturnPoolData.OpenOrders
	// raydiumMsg.Timestamp = (*resp)[0].Data.ReturnPoolData.PoolOpenTime

	// 转换并保存到数据库
	pool := liquidityPoolService.ConvertMessageToLiquidityPool(&raydiumMsg)

	serviceResponse := liquidityPoolService.ProcessTokenLiquidityPoolCreation(pool)
	if serviceResponse.Code != 0 {
		util.Log().Error("Failed to save pool info: %v", serviceResponse.Error)
		return fmt.Errorf("failed to save pool info: %v", serviceResponse.Error)
	}

	util.Log().Info("Successfully processed unknown raydium pool: %s", pool.PoolAddress)
	return nil
}

// gameOutTradeHandler 处理代理合约外盘买卖事件
func gameOutTradeHandler(message []byte, topic string) error {
	util.Log().Info("gameOutTradeHandler: Processing message from topic %s", topic)

	var tradeMsg model.GameOutTradeMessage
	if err := json.Unmarshal(message, &tradeMsg); err != nil {
		util.Log().Error("Failed to unmarshal game-out-trade message: %v", err)
		return fmt.Errorf("failed to unmarshal game-out-trade message: %v", err)
	}

	discount, _ := strconv.ParseUint(os.Getenv("DISCOUNT"), 10, 64)
	coefficient, _ := strconv.ParseUint(os.Getenv("COEFFICIENT"), 10, 64)
	PoolQuoteReserve, _ := strconv.ParseUint(tradeMsg.PoolQuoteReserve, 10, 64)
	PoolBaseReserve, _ := strconv.ParseUint(tradeMsg.PoolBaseReserve, 10, 64)
	FeeBaseAmount, _ := strconv.ParseUint(tradeMsg.FeeBaseAmount, 10, 64)

	discount_value := (100 - discount) / 100
	point := coefficient * FeeBaseAmount / (PoolQuoteReserve * discount_value / PoolBaseReserve)
	util.Log().Info("point: %d", point)

	pointRecordsRepo := model.NewPointRecordsRepo()
	userInfoRepo := model.NewUserInfoRepo()
	pointsService := service.NewPointsServiceImpl(userInfoRepo, pointRecordsRepo)
	err := pointsService.PointsSave(tradeMsg.User, point, tradeMsg.Signature, string(message))
	if err != nil {
		util.Log().Error("Failed to save points: %v", err)
		return fmt.Errorf("failed to save points: %v", err)
	}

	return nil
}

// gameInTradeHandler 处理代理合约内盘买事件（积分兑换买）
func gameInTradeHandler(message []byte, topic string) error {
	util.Log().Info("gameInTradeHandler: Processing message from topic %s", topic)

	var tradeMsg model.GameInTradeMessage
	if err := json.Unmarshal(message, &tradeMsg); err != nil {
		util.Log().Error("Failed to unmarshal game-in-trade message: %v", err)
		return fmt.Errorf("failed to unmarshal game-in-trade message: %v", err)
	}
	//TODO: 更新积分记录表（类型兑换并购买）

	pointRecordsRepo := model.NewPointRecordsRepo()
	userInfoRepo := model.NewUserInfoRepo()
	pointsService := service.NewPointsServiceImpl(userInfoRepo, pointRecordsRepo)
	err := pointsService.CreatePointRecord(tradeMsg.User, uint64(tradeMsg.PointsAmount*model.PointsDecimal), tradeMsg.Signature, string(message), 4, true)
	if err != nil {
		return fmt.Errorf("failed to save points: %v", err)
	}

	return nil
}
