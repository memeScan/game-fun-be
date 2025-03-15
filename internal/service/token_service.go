package service

import (
	"context"
	"encoding/json"
	"fmt"
	"game-fun-be/internal/conf"
	"game-fun-be/internal/constants"
	"game-fun-be/internal/es"
	"game-fun-be/internal/model"

	"game-fun-be/internal/pkg/util"
	"game-fun-be/internal/redis"
	"game-fun-be/internal/response"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

// TokenInfoService 币种信息服务
type TokenInfoService struct{}

// 批量插入代币信息
func (service *TokenInfoService) BatchInsertTokenInfo(tokenInfos []model.TokenInfo) response.Response {
	err := model.BatchInsertTokenInfo(tokenInfos)
	if err != nil {
		return response.Err(response.CodeDBError, "Failed to batch insert token info", err)
	}
	return response.Response{
		Code: 0,
		Msg:  "Token info batch inserted successfully",
	}
}

// GetTokenInfoByAddress 通过代币地址和链类型获取币种信息
func (service *TokenInfoService) GetTokenInfoByAddress(tokenAddress string, chainType uint8) (*model.TokenInfo, error) {
	return model.GetTokenInfoByAddress(tokenAddress, chainType)
}

// GetTokenInfoByAddress 通过代币地址和链类型获取币种信息
func (service *TokenInfoService) GetTokenInfoByAddresses(tokenAddresses []string, chainType uint8) ([]model.TokenInfo, error) {
	return model.GetTokenInfoByAddresses(tokenAddresses, chainType)
}

// UpdateTokenInfo 更新币种信息记录
func (service *TokenInfoService) UpdateTokenInfo(info *model.TokenInfo) error {
	return model.UpdateTokenInfo(info)
}

// DeleteTokenInfo 删除币种信息记录
func (service *TokenInfoService) DeleteTokenInfo(id int64) error {
	return model.DeleteTokenInfo(id)
}

// ListTokenInfos 列出币种信息记录
func (service *TokenInfoService) ListTokenInfos(limit, offset int) ([]model.TokenInfo, error) {
	return model.ListTokenInfos(limit, offset)
}

// ProcessTokenInfoCreation 处理币种信息记录创建
func (service *TokenInfoService) ProcessTokenInfoCreation(info *model.TokenInfo) response.Response {
	if err := service.CreateTokenInfoWithES(info); err != nil {
		return response.Err(response.CodeDBError, "Failed to create token info", err)
	}

	return response.Response{
		Code: 0,
		Data: info,
		Msg:  "Token info created successfully",
	}
}

// ProcessTokenInfoUpdate 处理币种信息记录更新
func (service *TokenInfoService) ProcessTokenInfoUpdate(info *model.TokenInfo) response.Response {
	err := service.UpdateTokenInfo(info)
	if err != nil {
		return response.Err(response.CodeDBError, "Failed to update token info", err)
	}

	return response.Response{
		Code: 0,
		Data: info,
		Msg:  "Token info updated successfully",
	}
}

// ProcessTokenInfoDeletion 处理币种信息记录删除
func (service *TokenInfoService) ProcessTokenInfoDeletion(id int64) response.Response {
	err := service.DeleteTokenInfo(id)
	if err != nil {
		return response.Err(response.CodeDBError, "Failed to delete token info", err)
	}

	return response.Response{
		Code: 0,
		Msg:  "Token info deleted successfully",
	}
}

// ProcessTokenInfoQuery 处理币种信息记录查询
func (service *TokenInfoService) ProcessTokenInfoQuery(tokenAddress string, chainType uint8) response.Response {
	info, err := service.GetTokenInfoByAddress(tokenAddress, chainType)
	if err != nil {
		return response.Err(response.CodeDBError, "Failed to query token info", err)
	}

	return response.Response{
		Code: 0,
		Data: info,
		Msg:  "Token info queried successfully",
	}
}

// ProcessTokenInfoList 处理币种信息记录列表查询
func (service *TokenInfoService) ProcessTokenInfoList(limit, offset int) response.Response {
	infos, err := service.ListTokenInfos(limit, offset)
	if err != nil {
		return response.Err(response.CodeDBError, "Failed to list token infos", err)
	}

	return response.Response{
		Code: 0,
		Data: infos,
		Msg:  "Token infos listed successfully",
	}
}

// ConvertTokenTransactionToInfo 将 TokenTransaction 转换为 TokenInfo
func (service *TokenInfoService) ConvertTokenTransactionToInfo(tx *model.TokenTransaction) *model.TokenInfo {
	return &model.TokenInfo{
		TokenAddress:        tx.TokenAddress,
		ChainType:           tx.ChainType,
		CreatedPlatformType: tx.PlatformType,
		TransactionTime:     tx.TransactionTime,
	}
}

// CreateOrUpdateTokenInfo 创建或更新币种信息
func (service *TokenInfoService) CreateOrUpdateTokenInfo(tx *model.TokenTransaction) response.Response {
	info, err := service.GetTokenInfoByAddress(tx.TokenAddress, tx.ChainType)
	if err != nil {
		// 如果不存在，创建新记录
		newInfo := service.ConvertTokenTransactionToInfo(tx)
		return service.ProcessTokenInfoCreation(newInfo)
	}

	// 如果存在，更新记录
	info.TransactionTime = tx.TransactionTime
	// 更新其他需要更新的字段
	return service.ProcessTokenInfoUpdate(info)
}

// ConvertMessageToTokenInfo 将 Kafka 消息转换为 TokenInfo
func (service *TokenInfoService) ConvertMessageToTokenInfo(msg *model.TokenInfoMessage) *model.TokenInfo {
	tokenInfo := &model.TokenInfo{}

	// 设置基本信息
	tokenInfo.TokenName = msg.Name
	tokenInfo.Symbol = msg.Symbol
	tokenInfo.Creator = msg.Creator
	tokenInfo.TokenAddress = msg.Mint

	// 设置链和创建平台类型
	tokenInfo.ChainType = uint8(model.ChainTypeSolana)
	tokenInfo.CreatedPlatformType = uint8(model.CreatedPlatformTypePump)

	// 设置代币相关数值
	tokenInfo.Decimals = model.CreatedPlatformType(tokenInfo.CreatedPlatformType).GetDecimals()
	tokenInfo.TotalSupply = uint64(decimal.NewFromInt(1000000000000000).IntPart())
	tokenInfo.CirculatingSupply = tokenInfo.TotalSupply

	// 设置交易相关信息
	tokenInfo.Block = msg.Block
	tokenInfo.TransactionHash = msg.Signature
	tokenInfo.TransactionTime = time.Unix(msg.Timestamp, 0)
	tokenInfo.PoolAddress = msg.BondingCurve
	tokenInfo.Progress = decimal.NewFromFloat(0.00)

	// 设置初始状态值
	tokenInfo.DevNativeTokenAmount = uint64(decimal.Zero.IntPart())
	tokenInfo.Holder = 0                              // 初始持有人数
	tokenInfo.CommentCount = 0                        // 初始评论数
	tokenInfo.MarketCap = decimal.Zero                // 始市值
	tokenInfo.CrownDuration = 0                       // 初始皇冠持续时间
	tokenInfo.RocketDuration = 0                      // 初始火箭持续时间
	tokenInfo.DevStatus = uint8(model.DevStatusClear) // 初始开发者状态
	tokenInfo.URI = msg.URI
	tokenInfo.IsMedia = false // 初始媒体类型状态

	// 设置 URI
	tokenInfo.URI = msg.URI
	if msg.URI != "" {
		content, err := GetURIContent(msg.URI, 1) // 设置 1 次重试
		if err != nil {
			util.Log().Error("Error fetching URI content: %v", err)
		} else {
			util.Log().Info("Successfully fetched URI content: %s", content)
			hasSocial := HasSocialMedia(content)
			tokenInfo.IsMedia = hasSocial
			tokenInfo.ExtInfo = content
		}
	}

	return tokenInfo
}

// CreateTokenInfoFromMessage 从 Kafka 消息创建 TokenInfo 并保存到数据库
func (service *TokenInfoService) CreateTokenInfoFromMessage(msg *model.TokenInfoMessage) response.Response {
	tokenInfo := service.ConvertMessageToTokenInfo(msg)
	return service.ProcessTokenInfoCreation(tokenInfo)
}

// BatchUpdateTokenInfo 批量更新币种信息记录
func (service *TokenInfoService) BatchUpdateTokenInfo(tokenInfos []*model.TokenInfo) response.Response {
	err := model.BatchUpdateTokenInfo(tokenInfos)
	if err != nil {
		return response.Err(response.CodeDBError, "批量更新代币信息失败", err)
	}

	return response.Response{
		Code: 0,
		Msg:  "批量更新代币信息成功",
	}
}

// GetExistingTokenInfos 批量获取代币信息（优先从Redis获取，未命中则从数据库查询）
func (service *TokenInfoService) GetExistingTokenInfos(addresses []string, chainType uint8) map[string]*model.TokenInfo {
	result := make(map[string]*model.TokenInfo)
	var missedAddresses []string

	// 构建 Redis keys
	redisKeys := make([]string, len(addresses))
	for i, address := range addresses {
		redisKeys[i] = fmt.Sprintf("%s:%s", constants.RedisKeyPrefixTokenInfo, address)
	}

	// 批量从 Redis 获取
	values, err := redis.MGet(redisKeys)
	if err != nil {
		util.Log().Error("Failed to batch get from Redis: %v", err)
		missedAddresses = addresses
	} else {
		// 处理返回结果
		for i, value := range values {
			if value == "" {
				missedAddresses = append(missedAddresses, addresses[i])
				continue
			}
			var tokenInfo model.TokenInfo
			if err := json.Unmarshal([]byte(value), &tokenInfo); err != nil {
				util.Log().Error("Failed to unmarshal token info from Redis: %v", err)
				missedAddresses = append(missedAddresses, addresses[i])
				continue
			}
			result[addresses[i]] = &tokenInfo
		}
	}

	// 对未命中的地址批量查询数据库
	if len(missedAddresses) > 0 {
		tokenInfos, err := model.GetTokenInfoByAddresses(missedAddresses, chainType)
		if err == nil {
			for _, info := range tokenInfos {
				result[info.TokenAddress] = &info
			}
		}
	}

	return result
}

// UpdateTokenComplete 更新代币的完成状态
func (s *TokenInfoService) UpdateTokenComplete(tokenAddress string, chainType uint8) *response.Response {
	if err := model.UpdateTokenComplete(tokenAddress, chainType); err != nil {
		util.Log().Error("Failed to update token complete status: %v", err)
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "Failed to update token complete status",
		}
	}

	return &response.Response{
		Code: http.StatusOK,
		Msg:  "Token complete status updated successfully",
	}
}

// TokenInfoToESDoc 将 TokenInfo 转换为 ES 文档结构
func (service *TokenInfoService) TokenInfoToESDoc(info *model.TokenInfo) (map[string]interface{}, error) {
	// 构建 ES 文档
	doc := map[string]interface{}{
		"id":                      info.ID,
		"token_name":              info.TokenName,
		"symbol":                  info.Symbol,
		"creator":                 info.Creator,
		"token_address":           info.TokenAddress,
		"chain_type":              uint8(info.ChainType),           // byte类型
		"created_platform_type":   uint8(info.CreatedPlatformType), // byte类型
		"decimals":                uint8(info.Decimals),            // byte类型
		"total_supply":            int64(info.TotalSupply),         // long类型
		"circulating_supply":      int64(info.CirculatingSupply),   // long类型
		"block":                   int64(info.Block),               // long类型
		"transaction_hash":        info.TransactionHash,
		"transaction_time":        info.TransactionTime, // date类型
		"uri":                     info.URI,
		"dev_native_token_amount": int64(info.DevNativeTokenAmount),           // long类型
		"dev_token_amount":        int64(info.DevTokenAmount),                 // long类型
		"holder":                  int32(info.Holder),                         // integer类型
		"comment_count":           int32(info.CommentCount),                   // integer类型
		"market_cap":              info.MarketCap.InexactFloat64(),            // double类型
		"circulating_market_cap":  info.CirculatingMarketCap.InexactFloat64(), // double类型
		"crown_duration":          int64(info.CrownDuration),                  // long类型
		"rocket_duration":         int64(info.RocketDuration),                 // long类型
		"dev_status":              uint8(info.DevStatus),                      // byte类型
		"is_media":                info.IsMedia,                               // boolean类型
		"is_complete":             info.IsComplete,                            // boolean类型
		"price":                   info.Price.InexactFloat64(),                // double类型
		"native_price":            info.NativePrice.InexactFloat64(),          // double类型
		"liquidity":               info.Liquidity.InexactFloat64(),            // double类型
		"ext_info":                info.ExtInfo,                               // text类型
		"create_time":             info.CreateTime,                            // date类型
		"update_time":             info.UpdateTime,                            // date类型
		"dev_percentage":          info.DevPercentage,                         // double类型
		"top10_percentage":        info.Top10Percentage,                       // double类型
		"burn_percentage":         info.BurnPercentage,                        // double类型
		"dev_burn_percentage":     info.DevBurnPercentage,                     // double类型
		"token_flags":             int32(info.TokenFlags),                     // integer类型
		"progress":                info.Progress.InexactFloat64(),             // double类型
		"pool_address":            info.PoolAddress,
	}

	return doc, nil
}

// CreateTokenInfoWithES 创建代币信息并同步到 ES
func (service *TokenInfoService) CreateTokenInfoWithES(info *model.TokenInfo) error {
	// 1. 创建数据库记录
	if err := model.CreateTokenInfo(info); err != nil {
		return fmt.Errorf("failed to create token info in DB: %w", err)
	}

	// 2. 转换为 ES 文档
	doc, err := service.TokenInfoToESDoc(info)
	if err != nil {
		return fmt.Errorf("failed to convert to ES doc: %w", err)
	}

	// 3. 索引到 ES
	_, err = es.ESClient.Index().
		Index(conf.ES_INDEX_TOKEN_INFO).
		Id(fmt.Sprintf("%s_%d", info.TokenAddress, info.ChainType)). // 使用 token_address_chainType 作为文档 ID
		BodyJson(doc).
		Do(context.Background())
	if err != nil {
		return fmt.Errorf("failed to index to ES: %w", err)
	}

	return nil
}

// GetTokenInfoMapByDB 直接从数据库获取token信息map
func (s *TokenInfoService) GetTokenInfoMapByDB(addresses []string, chainType uint8) (map[string]*model.TokenInfo, error) {
	result := make(map[string]*model.TokenInfo)

	// 直接从数据库获取
	tokenInfos, err := model.GetTokenInfoByAddresses(addresses, chainType)
	if err != nil {
		return nil, fmt.Errorf("get token info from db failed: %v", err)
	}

	// 构建map
	for _, info := range tokenInfos {
		result[info.TokenAddress] = &info
	}

	return result, nil
}
