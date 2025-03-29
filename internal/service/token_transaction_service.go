package service

import (
	"fmt"
	"game-fun-be/internal/model"
	"game-fun-be/internal/pkg/util"
	"game-fun-be/internal/response"
	"strconv"
	"strings"
	"time"

	"game-fun-be/internal/constants"
	"game-fun-be/internal/redis"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// TokenTransactionService 代币交易服务
type TokenTransactionService struct{}

// CreateTokenTransaction 创建代币交易记录
func (service *TokenTransactionService) CreateTokenTransaction(tx *model.TokenTransaction) (*model.TokenTransaction, error) {
	err := model.CreateTokenTransaction(tx, tx.TransactionTime.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// GetTokenTransactionByID 通过交易ID获取代币交易记录
func (service *TokenTransactionService) GetTokenTransactionByID(date string, transactionID uint64) (*model.TokenTransaction, error) {
	return model.GetTokenTransactionByID(date, transactionID)
}

// GetTokenTransactionByHash 通过交易哈希获取代币交易记录
func (service *TokenTransactionService) GetTokenTransactionByHash(date string, transactionHash string, tokenAddress string) (*model.TokenTransaction, error) {
	return model.GetTokenTransactionByHash(date, transactionHash, tokenAddress)
}

// UpdateTokenTransaction 更新代币交易记录
func (service *TokenTransactionService) UpdateTokenTransaction(tx *model.TokenTransaction) error {
	return model.UpdateTokenTransaction(tx)
}

// ListTokenTransactions 列出代币交易记录
func (service *TokenTransactionService) ListTokenTransactions(limit, offset int) ([]model.TokenTransaction, error) {
	return model.ListTokenTransactions(limit, offset)
}

// GetTokenTransactionsByUserAddress 通过用户地址获取代币交易记录
func (service *TokenTransactionService) GetTokenTransactionsByUserAddress(userAddress string, limit, offset int) ([]model.TokenTransaction, error) {
	return model.GetTokenTransactionsByUserAddress(userAddress, limit, offset)
}

// GetTokenTransactionsByChainAndToken 通过链类型和代币地址获取代币交易记录
func (service *TokenTransactionService) GetTokenTransactionsByChainAndToken(chainType uint8, tokenAddress string, limit, offset int) ([]model.TokenTransaction, error) {
	return model.GetTokenTransactionsByChainAndToken(chainType, tokenAddress, limit, offset)
}

// ProcessTokenTransactionCreation 处理代币交易记录创建
func (service *TokenTransactionService) ProcessTokenTransactionCreation(tx *model.TokenTransaction) response.Response {
	createdTx, err := service.CreateTokenTransaction(tx)
	if err != nil {
		return response.Err(response.CodeDBError, "Failed to create token transaction", err)
	}

	return response.Response{
		Code: 0,
		Data: createdTx,
		Msg:  "Token transaction created successfully",
	}
}

// ProcessTokenTransactionUpdate 处理代币交易记录更新
func (service *TokenTransactionService) ProcessTokenTransactionUpdate(tx *model.TokenTransaction) response.Response {
	err := service.UpdateTokenTransaction(tx)
	if err != nil {
		return response.Err(response.CodeDBError, "Failed to update token transaction", err)
	}

	return response.Response{
		Code: 0,
		Data: tx,
		Msg:  "Token transaction updated successfully",
	}
}

// ProcessTokenTransactionQuery 处理代币交易记录查询
func (service *TokenTransactionService) ProcessTokenTransactionQuery(hash string, date string, tokenAddress string) response.Response {
	tx, err := service.GetTokenTransactionByHash(hash, date, tokenAddress)
	if err != nil {
		return response.Err(response.CodeDBError, "Failed to query token transaction", err)
	}

	return response.Response{
		Code: 0,
		Data: tx,
		Msg:  "Token transaction queried successfully",
	}
}

// ConvertTradeMessageToTransaction 将 TokenTradeMessage 转换为 TokenTransaction
func (service *TokenTransactionService) ConvertTradeMessageToTransaction(msg *model.TokenTradeMessage) *model.TokenTransaction {
	tx := &model.TokenTransaction{}

	// 设置基本���段
	tx.TransactionHash = msg.Signature
	tx.TokenAddress = msg.Mint
	tx.TransactionTime = time.Unix(msg.Timestamp, 0)
	tx.Block = msg.Block
	tx.CreateTime = time.Now()
	tx.UpdateTime = time.Now()
	tx.PlatformType = uint8(model.PlatformTypePump)
	tx.ChainType = uint8(model.ChainTypeSolana)

	// 转换 SolAmount
	solAmount, err := strconv.ParseUint(msg.SolAmount, 10, 64)
	if err != nil {
		util.Log().Error("Error parsing SolAmount: %v", err)
		solAmount = 0
	}
	tx.NativeTokenAmount = solAmount

	// 转换 TokenAmount
	tokenAmount, err := strconv.ParseUint(msg.TokenAmount, 10, 64)
	if err != nil {
		util.Log().Error("Error parsing TokenAmount: %v", err)
		tokenAmount = 0
	}
	tx.TokenAmount = tokenAmount

	tx.UserAddress = msg.User

	// 转换虚拟储备
	virtualSolReserves, err := strconv.ParseUint(msg.VirtualSolReserves, 10, 64)
	if err != nil {
		util.Log().Error("Error parsing VirtualSolReserves: %v", err)
		virtualSolReserves = 0
	}
	tx.VirtualNativeReserves = virtualSolReserves

	virtualTokenReserves, err := strconv.ParseUint(msg.VirtualTokenReserves, 10, 64)
	if err != nil {
		util.Log().Error("Error parsing VirtualTokenReserves: %v", err)
		virtualTokenReserves = 0
	}
	tx.VirtualTokenReserves = virtualTokenReserves
	realSolReserves, err := strconv.ParseUint(msg.RealSolReserves, 10, 64)
	if err != nil {
		util.Log().Error("Error parsing RealSolReserves: %v", err)
		realSolReserves = 0
	}
	tx.RealNativeReserves = realSolReserves

	realTokenReserves, err := strconv.ParseUint(msg.RealTokenReserves, 10, 64)
	if err != nil {
		util.Log().Error("Error parsing RealTokenReserves: %v", err)
		realTokenReserves = 0
	}
	tx.RealTokenReserves = realTokenReserves

	tx.IsBuy = msg.IsBuy

	// 设置进度
	tx.Progress = decimal.NewFromFloat(msg.Progress)
	tx.IsComplete = false // 默认未完成，直到收到代币完成消息

	solPrice, err := getSolPrice()
	if err != nil {
		util.Log().Error("Error getting SOL price: %v", err)
		solPrice = decimal.Zero // 设置为 0
	}
	tx.NativePriceUSD = solPrice

	// SOL精度固定为9
	const solDecimals = 9
	// 获取Token精度（从平台类型定义中获取）
	tokenDecimals := 6 //
	tx.Decimals = uint8(tokenDecimals)

	// 设置价格
	service.calculateAndSetPrices(tx, solDecimals, tokenDecimals, solPrice)

	// 设置市场地址
	tx.MarketAddress = msg.BondingCurve
	// 设置池子地址
	tx.PoolAddress = msg.BondingCurve

	// 设置原生代币地址
	tx.NativeTokenAddress = model.ChainType(tx.ChainType).GetNativeTokenAddress()
	// 设置交易类型
	if msg.IsBuy {
		tx.TransactionType = 1 // 买入
	} else {
		tx.TransactionType = 2 // 卖出
	}

	return tx
}

func (s *TokenTransactionService) GetESDoc(tx *model.TokenTransaction, tokenInfoMap map[string]*model.TokenInfo, poolInfoMap map[string]*model.TokenLiquidityPool) map[string]interface{} {
	doc := make(map[string]interface{})

	// token_transaction表字段
	doc["id"] = strconv.FormatUint(tx.ID, 10)
	doc["transaction_hash"] = tx.TransactionHash
	doc["token_address"] = tx.TokenAddress
	doc["user_address"] = tx.UserAddress
	doc["token_amount"] = strconv.FormatUint(tx.TokenAmount, 10)
	doc["native_token_amount"] = strconv.FormatUint(tx.NativeTokenAmount, 10)
	doc["price"] = tx.Price.InexactFloat64()
	doc["native_price"] = tx.NativePrice.InexactFloat64()
	doc["transaction_time"] = tx.TransactionTime
	doc["create_time"] = tx.CreateTime
	doc["update_time"] = tx.UpdateTime
	doc["chain_type"] = tx.ChainType
	doc["platform_type"] = tx.PlatformType
	doc["is_buy"] = tx.IsBuy
	doc["progress"] = tx.Progress.InexactFloat64()
	doc["is_complete"] = tx.IsComplete
	doc["virtual_native_reserves"] = strconv.FormatUint(tx.VirtualNativeReserves, 10)
	doc["virtual_token_reserves"] = strconv.FormatUint(tx.VirtualTokenReserves, 10)
	doc["real_native_reserves"] = strconv.FormatUint(tx.RealNativeReserves, 10)
	doc["real_token_reserves"] = strconv.FormatUint(tx.RealTokenReserves, 10)
	//设置代币信息相关字段
	doc["pc_symbol"] = "WSOL"
	doc["decimals"] = tx.Decimals

	// Use tokenInfoMap instead of Redis/DB lookup
	if tokenInfo, exists := tokenInfoMap[tx.TokenAddress]; exists {
		updateDocWithTokenInfo(doc, tokenInfo)
	} else {
		// Fallback to default if not found in map
		tokenInfo = s.getDefaultOrFallbackTokenInfo(tx)
		updateDocWithTokenInfo(doc, tokenInfo)
	}

	if poolInfo, exists := poolInfoMap[tx.PoolAddress]; exists {
		updateDocWithPoolInfo(doc, poolInfo)
	} else {
		// Fallback to default if not found in map
		poolInfo = s.getDefaultOrFallbackPoolInfo(tx)
		//取代币创建时间作为池子默认创建时间
		if tokenInfo, exists := tokenInfoMap[tx.TokenAddress]; exists {
			poolInfo.BlockTime = tokenInfo.TransactionTime
		}
		updateDocWithPoolInfo(doc, poolInfo)
	}

	// 移除所有为 nil 的字段
	for key, value := range doc {
		if value == nil {
			delete(doc, key)
		}
	}

	return doc
}

// getDefaultOrFallbackTokenInfo 获取默认或备用的 TokenInfo
func (s *TokenTransactionService) getDefaultOrFallbackTokenInfo(tx *model.TokenTransaction) *model.TokenInfo {
	return &model.TokenInfo{
		TokenAddress:         tx.TokenAddress,
		ChainType:            tx.ChainType,
		TokenName:            "",
		Symbol:               "",
		Creator:              "",
		CreatedPlatformType:  0,
		Decimals:             0,
		TotalSupply:          0,
		CirculatingSupply:    0,
		Block:                0,
		TransactionHash:      "",
		TransactionTime:      time.Time{},
		URI:                  "",
		DevNativeTokenAmount: 0,
		DevTokenAmount:       0,
		Holder:               0,
		CommentCount:         0,
		MarketCap:            decimal.Zero,
		CirculatingMarketCap: decimal.Zero,
		CrownDuration:        0,
		RocketDuration:       0,
		DevStatus:            0,
		IsMedia:              false,
		IsComplete:           false,
		Price:                decimal.Zero,
		NativePrice:          decimal.Zero,
		Liquidity:            decimal.Zero,
		ExtInfo:              "",
		CreateTime:           time.Time{},
		UpdateTime:           time.Time{},
		DevPercentage:        0,
		Top10Percentage:      0,
		BurnPercentage:       0,
		DevBurnPercentage:    0,
		TokenFlags:           0,
	}
}

// updateDocWithTokenInfo 更新文档with token信息
func updateDocWithTokenInfo(doc map[string]interface{}, tokenInfo *model.TokenInfo) {
	if tokenInfo == nil {
		return
	}
	doc["is_complete"] = tokenInfo.IsComplete // 设置代币完成状态
	doc["token_create_time"] = tokenInfo.TransactionTime
	doc["dev_native_token_amount"] = tokenInfo.DevNativeTokenAmount
	doc["holder"] = tokenInfo.Holder
	doc["comment_count"] = tokenInfo.CommentCount
	doc["market_cap"] = tokenInfo.MarketCap.InexactFloat64()
	doc["token_supply"] = strconv.FormatUint(tokenInfo.TotalSupply, 10)
	doc["dev_status"] = tokenInfo.DevStatus
	doc["created_platform_type"] = tokenInfo.CreatedPlatformType
	doc["token_creator"] = tokenInfo.Creator
	// 添加 token_name 和 symbol
	doc["token_name"] = tokenInfo.TokenName
	doc["symbol"] = tokenInfo.Symbol

	// 添加烧池子和DEX广告标记
	if tokenInfo.CreatedPlatformType == uint8(model.CreatedPlatformTypePump) {
		doc["is_burned_lp"] = true
	} else {
		doc["is_burned_lp"] = tokenInfo.HasFlag(model.FLAG_BURNED_LP)
	}
	doc["is_dex_ad"] = tokenInfo.HasFlag(model.FLAG_DXSCR_AD)

	//
	if tokenInfo.CrownDuration != 0 {
		doc["crown_duration"] = tokenInfo.CrownDuration
	}
	if tokenInfo.RocketDuration != 0 {
		doc["rocket_duration"] = tokenInfo.RocketDuration
	}
	doc["is_media"] = tokenInfo.IsMedia
	doc["uri"] = tokenInfo.URI
	doc["ext_info"] = tokenInfo.ExtInfo
	doc["liquidity"] = tokenInfo.Liquidity.InexactFloat64()
	doc["token_flags"] = tokenInfo.TokenFlags
	doc["burn_percentage"] = tokenInfo.BurnPercentage
	doc["dev_burn_percentage"] = tokenInfo.DevBurnPercentage
	doc["dev_percentage"] = tokenInfo.DevPercentage
	doc["top_10_percentage"] = tokenInfo.Top10Percentage
	doc["dev_token_amount"] = strconv.FormatUint(tokenInfo.DevTokenAmount, 10)

}

// getDefaultOrFallbackPoolInfo 获取默认或备用的池子信息
func (s *TokenTransactionService) getDefaultOrFallbackPoolInfo(tx *model.TokenTransaction) *model.TokenLiquidityPool {
	return &model.TokenLiquidityPool{
		PoolAddress:  tx.PoolAddress,
		CoinAddress:  tx.TokenAddress,
		ChainType:    tx.ChainType,
		PlatformType: tx.PlatformType,
		PcAddress:    tx.NativeTokenAddress,
		CreateTime:   time.Time{},
		UpdateTime:   time.Time{},
		BlockTime:    time.Time{},
	}
}

// updateDocWithPoolInfo 更新文档with池子信息
func updateDocWithPoolInfo(doc map[string]interface{}, poolInfo *model.TokenLiquidityPool) {
	if poolInfo == nil {
		return
	}
	// 设置池子相关字段
	doc["pool_address"] = poolInfo.PoolAddress
	doc["native_token_address"] = poolInfo.PcAddress
	doc["block_time"] = poolInfo.BlockTime                //池子创建时间
	doc["creator"] = poolInfo.UserAddress                 //池子创建者地址
	doc["initial_pc_reserve"] = poolInfo.InitialPcReserve //池子初始定价代币总量

}

// BatchCreateTokenTransactions 批量创建代币交易记录
func (service *TokenTransactionService) BatchCreateTokenTransactions(txs []*model.TokenTransaction, date string) error {
	if len(txs) == 0 {
		return nil
	}
	return model.BatchCreateTokenTransactions(txs, date)
}

// ProcessBatchTokenTransactionCreation 处理批量代币交易记录创建
func (service *TokenTransactionService) ProcessBatchTokenTransactionCreation(txs []*model.TokenTransaction, date string) response.Response {
	err := service.BatchCreateTokenTransactions(txs, date)
	if err != nil {
		return response.Err(response.CodeDBError, "Failed to create batch token transactions", err)
	}

	return response.Response{
		Code: 0,
		Msg:  fmt.Sprintf("%d token transactions created successfully", len(txs)),
	}
}

// ConvertTradeMessagesToTransactions 将 TokenTradeMessage 列表转换为 TokenTransaction 列表
func (service *TokenTransactionService) ConvertTradeMessagesToTransactions(messages []*model.TokenTradeMessage) []*model.TokenTransaction {
	transactions := make([]*model.TokenTransaction, 0, len(messages))
	for _, msg := range messages {
		tx := service.ConvertTradeMessageToTransaction(msg)
		transactions = append(transactions, tx)
	}
	return transactions
}

// GetESDocList 通过 TokenTransaction 列表获取 ES 文档列表
func (service *TokenTransactionService) GetESDocList(transactions []*model.TokenTransaction, tokenInfoMap map[string]*model.TokenInfo, poolInfoMap map[string]*model.TokenLiquidityPool) []map[string]interface{} {
	docList := make([]map[string]interface{}, 0, len(transactions))
	for _, tx := range transactions {
		doc := service.GetESDoc(tx, tokenInfoMap, poolInfoMap) // Pass tokenInfoMap to GetESDoc
		docList = append(docList, doc)
	}
	return docList
}

// getSolPrice 获取 SOL 的最新价格
func getSolPrice() (decimal.Decimal, error) {
	priceStr, err := redis.Get(constants.RedisKeySolLatestPrice)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get SOL price from Redis: %w", err)
	}

	// 移除字符串中的引号
	priceStr = strings.Trim(priceStr, "\"")

	intPrice, err := strconv.ParseInt(priceStr, 10, 64)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to parse SOL price: %w", err)
	}

	return decimal.NewFromInt(intPrice).Div(decimal.NewFromInt(constants.SolPriceMultiplier)), nil
}

// ConvertRaydiumSwapMessagesToTransactions 将 RaydiumSwapMessage 列表转换为 TokenTransaction 列表
func (service *TokenTransactionService) ConvertRaydiumSwapMessagesToTransactions(messages []*model.RaydiumSwapMessage) []*model.TokenTransaction {
	transactions := make([]*model.TokenTransaction, 0, len(messages))
	for _, msg := range messages {
		tx := &model.TokenTransaction{}

		// 设置基本字段
		tx.TransactionHash = msg.Signature
		tx.TokenAddress = msg.QuoteToken
		tx.TransactionTime = time.Unix(msg.Timestamp, 0)
		tx.Block = msg.Block
		tx.CreateTime = time.Now()
		tx.UpdateTime = time.Now()
		tx.PlatformType = uint8(model.PlatformTypeRaydium)
		tx.ChainType = uint8(model.ChainTypeSolana)

		// 转换 BaseAmount
		baseAmount, err := strconv.ParseUint(msg.BaseAmount, 10, 64)
		if err != nil {
			util.Log().Error("Error parsing BaseAmount: %v", err)
			baseAmount = 0
		}
		tx.NativeTokenAmount = baseAmount

		// 转换 QuoteAmount
		quoteAmount, err := strconv.ParseUint(msg.QuoteAmount, 10, 64)
		if err != nil {
			util.Log().Error("Error parsing QuoteAmount: %v", err)
			quoteAmount = 0
		}
		tx.TokenAmount = quoteAmount

		tx.UserAddress = msg.User
		tx.IsBuy = msg.IsBuy

		// 设置地址相关字段
		tx.MarketAddress = msg.MarketAddress
		tx.PoolAddress = msg.PoolAddress
		tx.NativeTokenAddress = msg.BaseToken

		// 设置交易类型
		if msg.IsBuy {
			tx.TransactionType = 1 // 买入
		} else {
			tx.TransactionType = 2 // 卖出
		}

		// 设置虚拟储备
		baseReserve, err := decimal.NewFromString(msg.PoolBaseReserve)
		if err != nil {
			baseReserve = decimal.Zero // 如果解析失败，设置为零值
		}
		quoteReserve, err := decimal.NewFromString(msg.PoolQuoteReserve)
		if err != nil {
			quoteReserve = decimal.Zero // 如果解析失败，设置为零值
		}
		tx.VirtualNativeReserves = uint64(baseReserve.IntPart())
		tx.VirtualTokenReserves = uint64(quoteReserve.IntPart())

		// 计算池子流动性
		const baseDecimals = 9
		var quoteDecimals = msg.Decimals
		tx.Decimals = uint8(msg.Decimals)

		// 获取SOL价格
		solPrice, err := getSolPrice()
		if err != nil {
			util.Log().Error("Error getting SOL price: %v", err)
			solPrice = decimal.Zero
		}
		tx.NativePriceUSD = solPrice

		// 使用统一的价格计算方法
		service.calculateAndSetPrices(tx, baseDecimals, quoteDecimals, solPrice)

		// 设置交易完成
		tx.Progress = decimal.NewFromInt(100)
		tx.IsComplete = true

		// 判断是否是回购地址
		if tx.UserAddress == model.TreasuryAddress {
			tx.IsBuyback = true //是回购交易
		} else {
			tx.IsBuyback = false
		}

		// 判断是否是game代理地址
		if msg.ParentInstAddress == model.GameProxyAddress {
			tx.ProxyType = uint8(model.ProxyTypeGame)
		}

		transactions = append(transactions, tx)
	}
	return transactions
}

func (service *TokenTransactionService) ConvertPumpAmmSwapMessagesToTransactions(pumpAmmSwapMessages []*model.PumpAmmSwapMessage) []*model.TokenTransaction {
	transactions := make([]*model.TokenTransaction, 0, len(pumpAmmSwapMessages))
	for _, msg := range pumpAmmSwapMessages {
		tx := &model.TokenTransaction{}

		// 设置基本字段
		tx.TransactionHash = msg.Signature
		tx.TokenAddress = msg.QuoteToken
		tx.TransactionTime = time.Unix(msg.Timestamp, 0)
		tx.Block = msg.Block
		tx.CreateTime = time.Now()
		tx.UpdateTime = time.Now()
		tx.PlatformType = uint8(model.PlatformTypePumpSwap)
		tx.ChainType = uint8(model.ChainTypeSolana)

		// 转换 BaseAmount
		baseAmount, err := strconv.ParseUint(msg.BaseAmount, 10, 64)
		if err != nil {
			util.Log().Error("Error parsing BaseAmount: %v", err)
			baseAmount = 0
		}
		tx.NativeTokenAmount = baseAmount

		// 转换 QuoteAmount
		quoteAmount, err := strconv.ParseUint(msg.QuoteAmount, 10, 64)
		if err != nil {
			util.Log().Error("Error parsing QuoteAmount: %v", err)
			quoteAmount = 0
		}
		tx.TokenAmount = quoteAmount

		tx.UserAddress = msg.User
		tx.IsBuy = msg.IsBuy

		// 设置地址相关字段
		tx.PoolAddress = msg.PoolAddress
		tx.NativeTokenAddress = msg.BaseToken

		// 设置交易类型
		if msg.IsBuy {
			tx.TransactionType = 1 // 买入
		} else {
			tx.TransactionType = 2 // 卖出
		}

		// 设置虚拟储备
		baseReserve, err := decimal.NewFromString(msg.PoolBaseReserve)
		if err != nil {
			baseReserve = decimal.Zero // 如果解析失败，设置为零值
		}
		quoteReserve, err := decimal.NewFromString(msg.PoolQuoteReserve)
		if err != nil {
			quoteReserve = decimal.Zero // 如果解析失败，设置为零值
		}
		tx.VirtualNativeReserves = uint64(baseReserve.IntPart())
		tx.VirtualTokenReserves = uint64(quoteReserve.IntPart())

		// 计算池子流动性
		const baseDecimals = 9
		var quoteDecimals = msg.Decimals
		tx.Decimals = uint8(msg.Decimals)

		// 获取SOL价格
		solPrice, err := getSolPrice()
		if err != nil {
			util.Log().Error("Error getting SOL price: %v", err)
			solPrice = decimal.Zero
		}
		tx.NativePriceUSD = solPrice

		// 使用统一的价格计算方法
		service.calculateAndSetPrices(tx, baseDecimals, quoteDecimals, solPrice)

		// 设置交易完成
		tx.Progress = decimal.NewFromInt(100)
		tx.IsComplete = true

		// 判断是否是回购地址
		if tx.UserAddress == model.TreasuryAddress {
			tx.IsBuyback = true //是回购交易
		} else {
			tx.IsBuyback = false
		}

		// 判断是否是game代理地址
		if msg.ParentInstAddress == model.GameProxyAddress {
			tx.ProxyType = uint8(model.ProxyTypeGame)
		}

		transactions = append(transactions, tx)
	}
	return transactions
}

// 新增方法
func (service *TokenTransactionService) calculateAndSetPrices(tx *model.TokenTransaction, solDecimals, tokenDecimals int, solPrice decimal.Decimal) {
	if tx.TokenAmount == 0 {
		tx.NativePrice = decimal.Zero
		tx.Price = decimal.Zero
		tx.TransactionAmountUSD = decimal.Zero
		util.Log().Warning("Invalid amount for price calculation: solAmount=%d, tokenAmount=%d",
			tx.NativeTokenAmount, tx.TokenAmount)
		return
	}

	actualSolAmount := decimal.NewFromUint64(tx.NativeTokenAmount).Shift(-int32(solDecimals))
	actualTokenAmount := decimal.NewFromUint64(tx.TokenAmount).Shift(-int32(tokenDecimals))
	tx.NativePrice = actualSolAmount.Div(actualTokenAmount)

	if !solPrice.IsZero() {
		tx.Price = tx.NativePrice.Mul(solPrice)
	} else {
		tx.Price = decimal.Zero
		util.Log().Warning("SOL price is zero, unable to calculate USD price")
	}
	tx.TransactionAmountUSD = tx.Price.Mul(actualTokenAmount)
}

// GetLatestTokenTransaction 获取代币最新的一条交易记录
func (service *TokenTransactionService) GetLatestTokenTransaction(tokenAddress string, chainType uint8) (*model.TokenTransaction, error) {
	// 1. 从索引表获取最新记录
	latestIndex, err := model.GetLatestTokenTxIndexByTokenAddress(tokenAddress, chainType)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("get latest token tx index error: %v", err)
	}

	// 2. 用索引信息查询具体交易记录
	dateStr := latestIndex.TransactionDate.Format("20060102")
	tx, err := service.GetTokenTransactionByID(dateStr, latestIndex.TransactionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("get token transaction error: %v", err)
	}

	return tx, nil
}
