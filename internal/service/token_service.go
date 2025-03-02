package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"my-token-ai-be/internal/clickhouse"
	"my-token-ai-be/internal/constants"
	"my-token-ai-be/internal/es"
	"my-token-ai-be/internal/model"
	"my-token-ai-be/internal/pkg/httpRespone"
	"my-token-ai-be/internal/pkg/httpUtil"
	"my-token-ai-be/internal/pkg/util"
	"my-token-ai-be/internal/redis"
	"my-token-ai-be/internal/response"
	"net/http"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
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

// GetTokenInfo retrieves token information by address and chain type
func (service *TokenInfoService) GetTokenInfo(tokenAddress string, chainType uint8) *response.Response {
	info, err := model.GetTokenInfoByAddress(tokenAddress, chainType)

	if err != nil {
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  err.Error(),
		}
	}
	if info == nil {
		return &response.Response{
			Code: http.StatusNotFound,
			Msg:  "Token not found",
		}
	}

	var extInfo model.ExtInfo

	extInfo.Name = info.TokenName
	extInfo.Symbol = info.Symbol
	extInfo.Image = info.URI

	if info.ExtInfo != "" {
		err = json.Unmarshal([]byte(info.ExtInfo), &extInfo)
		if err != nil {
			return &response.Response{
				Code: http.StatusInternalServerError,
				Msg:  err.Error(),
			}
		}
	}

	progress := decimal.NewFromInt(100)
	marketCap := 0.0
	tokenOutLatestTxIndex, err := model.GetLatestTokenTxIndexByTokenAddress(tokenAddress, uint8(model.ChainTypeSolana))

	if tokenOutLatestTxIndex != nil && err == nil {
		formattedDate := tokenOutLatestTxIndex.TransactionDate.Format("20060102")
		tokenOutLatestTransaction, err := model.GetTokenTransactionByHash(formattedDate, tokenOutLatestTxIndex.TransactionHash, tokenOutLatestTxIndex.TokenAddress)
		if err == nil {
			progress = tokenOutLatestTransaction.Progress
			marketCap = tokenOutLatestTransaction.Price.InexactFloat64() * float64(info.CirculatingSupply) / math.Pow(10, float64(info.Decimals))
			info.Price = tokenOutLatestTransaction.Price
		}
	}

	balance := 0
	Top10Holdings := 0.00
	TotalHolders := false
	if info.Top10Percentage > 0.3 {
		TotalHolders = true
	}
	CtoFlag := info.HasFlag(model.FLAG_CTO)
	DexscrAd := info.HasFlag(model.FLAG_DXSCR_AD)
	MintAuthority := info.HasFlag(model.FLAG_MINT_AUTHORITY)
	FreezeAuthority := info.HasFlag(model.FLAG_FREEZE_AUTHORITY)
	DexscrUpdateLink := info.HasFlag(model.FLAG_DEXSCR_UPDATE)

	var creatorBalanceRate decimal.Decimal
	if info.TotalSupply != 0 {
		creatorBalanceRate = decimal.NewFromFloat(float64(balance)).Div(decimal.NewFromInt(int64(info.TotalSupply)))
	} else {
		creatorBalanceRate = decimal.Zero
	}

	totalSupply := float64(info.TotalSupply)

	top10HolderRate := CalculateTop10HolderRate(Top10Holdings, totalSupply)

	return &response.Response{
		Code: http.StatusOK,
		Data: &response.TokenInfoResponse{
			Address:              info.TokenAddress,
			Symbol:               extInfo.Symbol,
			Name:                 extInfo.Name,
			Logo:                 extInfo.Image,
			Website:              extInfo.Website,
			Twitter:              extInfo.Twitter,
			Telegram:             extInfo.Telegram,
			Decimals:             int(info.Decimals),
			BiggestPoolAddress:   "pool",
			MarketCap:            marketCap,
			Price:                info.Price.InexactFloat64(),
			Creator:              info.Creator,
			TotalSupply:          info.TotalSupply,
			Progress:             progress,
			HolderCount:          info.Holder,
			CirculatingSupply:    info.CirculatingSupply,
			Liquidity:            info.Liquidity,
			CreateTimestamp:      info.TransactionTime.Unix(),
			OpenTimestamp:        info.TransactionTime.Unix(),
			DevNativeTokenAmount: info.DevNativeTokenAmount,
			IsComplete:           info.IsComplete,
			TotalHolders:         TotalHolders,
			MintAuthority:        MintAuthority,
			FreezeAuthority:      FreezeAuthority,
			CtoFlag:              CtoFlag,
			CreatorBalanceRate:   creatorBalanceRate.String(),
			// RatTraderAmountRate:  0,
			// DevTokenBurnAmount:   info.DevBurnPercentage,
			// DevTokenBurnRatio:    info.DevBurnPercentage,
			BurnPercentage:      info.BurnPercentage,
			Top10Holdings:       info.Top10Percentage,
			Top10HolderRate:     top10HolderRate,
			DexscrAd:            DexscrAd,
			DexscrUpdateLink:    DexscrUpdateLink,
			CreatedPlatformType: info.CreatedPlatformType,
		},
		Msg: "Token info queried successfully",
	}
}

func (service *TokenInfoService) ConvertDefiTokenInfoToTokenInfo(defiTokenInfo *httpRespone.Token) *model.TokenInfo {
	tokenInfo := &model.TokenInfo{}
	tokenInfo.TokenAddress = defiTokenInfo.Mint
	tokenInfo.Creator = defiTokenInfo.Creator
	tokenInfo.Symbol = defiTokenInfo.Symbol
	tokenInfo.TokenName = defiTokenInfo.Name
	tokenInfo.URI = defiTokenInfo.URI
	resp, err := httpUtil.GetHTTPClient().Get(defiTokenInfo.URI)
	if err != nil {
		util.Log().Error("Failed to get URI: %v", err)
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	tokenInfo.ExtInfo = string(bodyBytes)
	timestamp, err := strconv.ParseInt(defiTokenInfo.Timestamp, 10, 64)
	if err != nil {
		util.Log().Error("Failed to parse timestamp: %v", err)
	}
	tokenInfo.TransactionTime = time.Unix(timestamp, 0)
	return tokenInfo
}

func (service *TokenInfoService) GetTokenPrices(tokenAddresses []string, chainType model.ChainType) response.Response {
	infos, err := model.GetTokenInfoByAddresses(tokenAddresses, uint8(chainType))

	if err != nil {
		return response.Err(response.CodeDBError, "Failed to get token info", err)
	}
	return response.Response{
		Code: http.StatusOK,
		Data: response.BuildTokenPriceResponse(infos),
		Msg:  "Token prices queried successfully",
	}
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

// BatchUpdateTokenInfo 批量更新币种信息记录
func (service *TokenInfoService) BatchUpdateTokenInfoV2(tokenInfos []*model.TokenInfo) response.Response {
	err := model.BatchUpdateTokenInfo(tokenInfos)
	if err != nil {
		return response.Err(response.CodeDBError, "批量更新代币信息失败", err)
	}

	// 返回成功修改的代币列表
	return response.Response{
		Code: 0,
		Msg:  "批量更新代币信息成功",
		Data: tokenInfos, // 将成功修改的代币列表返回
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

func (service *TokenInfoService) GetTokenLaunchpadInfo(tokenAddress string, chainType uint8) *response.Response {

	tokenOutLatestTxIndex, err := model.GetLatestTokenTxIndexByTokenAddress(tokenAddress, chainType)
	if tokenOutLatestTxIndex == nil {
		return &response.Response{
			Code: http.StatusNotFound,
			Msg:  "Token not found",
		}
	}
	if err != nil {
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "failed to get token latest transaction index",
		}
	}
	formattedDate := tokenOutLatestTxIndex.TransactionDate.Format("20060102")

	tokenOutLatestTransaction, err := model.GetTokenTransactionByHash(formattedDate, tokenOutLatestTxIndex.TransactionHash, tokenOutLatestTxIndex.TokenAddress)
	if err != nil {
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "failed to get token latest transaction",
		}
	}

	info, err := model.GetTokenInfoByAddress(tokenAddress, chainType)
	if err != nil {
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "Failed to get token info",
		}
	}

	type ExtInfo struct {
		Name        string `json:"name"`
		Symbol      string `json:"symbol"`
		Description string `json:"description"`
		Image       string `json:"image"`
		ShowName    bool   `json:"showName"`
		CreatedOn   string `json:"createdOn"`
		Creator     string `json:"creator"`
	}

	var extInfo ExtInfo
	err = json.Unmarshal([]byte(info.ExtInfo), &extInfo)
	if err != nil {
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "Failed to get token info",
		}
	}

	LaunchpadStatus := 0
	if tokenOutLatestTransaction.Progress.Cmp(decimal.NewFromInt(1)) > 0 {
		LaunchpadStatus = 1
	}

	return &response.Response{
		Code: http.StatusOK,
		Msg:  "Token launchpad info queried successfully",
		Data: response.TokenLaunchpadInfo{
			Address:           tokenAddress,
			Launchpad:         model.CreatedPlatformType(info.CreatedPlatformType).String(),
			LaunchpadStatus:   LaunchpadStatus,
			LaunchpadProgress: tokenOutLatestTransaction.Progress.String(),
			Description:       extInfo.Description,
		},
	}
}

func CalculateTop10HolderRate(top10Holdings float64, totalSupply float64) float64 {
	if totalSupply == 0 {
		// Handle the case where totalSupply is zero
		// Return 0.0 or another appropriate default value
		return 0.0
	}
	return top10Holdings / totalSupply
}

func (service *TokenInfoService) GetTokenOrderBook(tokenAddress string, chainType uint8) response.Response {
	transactions, err := clickhouse.GetTokenTransactions(tokenAddress, 100)
	if err != nil {
		return response.Err(response.CodeDBError, "failed to get token order book", err)
	}

	convertedTransactions := make([]response.TokenOrderBookItem, len(transactions))
	for i, tx := range transactions {
		// 计算基础代币数量
		baseAmount := decimal.NewFromInt(int64(tx.BaseTokenAmount))
		if tx.Decimals > 0 {
			baseAmount = baseAmount.Shift(-9)
		}

		// 计算报价代币数量
		quoteAmount := decimal.NewFromInt(int64(tx.QuoteTokenAmount))
		if tx.Decimals > 0 {
			quoteAmount = quoteAmount.Shift(-int32(tx.Decimals))
		}

		// 计算 USD 金额
		var usdAmount decimal.Decimal
		if tx.TransactionType == 1 { // 1 是买入
			usdAmount = quoteAmount.Mul(tx.QuoteTokenPrice)
		} else {
			usdAmount = baseAmount.Mul(tx.BaseTokenPrice)
		}

		// 创建一个新的 TokenOrderBookItem
		item := response.TokenOrderBookItem{}

		// 逐个赋值
		item.TransactionHash = tx.TransactionHash
		item.ChainType = tx.ChainType
		item.UserAddress = tx.UserAddress
		item.TokenAddress = tx.TokenAddress
		item.PoolAddress = tx.PoolAddress
		item.BaseTokenAmount = baseAmount
		item.QuoteTokenAmount = quoteAmount
		item.BaseTokenPrice = tx.BaseTokenPrice
		item.QuoteTokenPrice = tx.QuoteTokenPrice
		item.TransactionType = tx.TransactionType
		item.PlatformType = tx.PlatformType
		item.TransactionTime = tx.TransactionTime.Unix()
		item.UsdAmount = usdAmount

		// 将item赋值给切片
		convertedTransactions[i] = item
	}

	return response.BuildTokenOrderBookResponse(convertedTransactions)
}

func (service *TokenInfoService) SearchToken(searchName string, chainType uint8) *response.Response {

	isTokenAddress := false
	if len(searchName) >= 32 && isBase58(searchName) {
		isTokenAddress = true
	}

	queryJSON, err := es.SearchTokenBySymbol(searchName, chainType, isTokenAddress)
	if err != nil {
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "Failed to get pump rank",
		}
	}

	result, err := es.SearchDocuments(es.ES_INDEX_TOKEN_INFO, queryJSON)
	if err != nil {
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "Failed to get pump rank",
		}
	}
	jsonResult, err := json.Marshal(result)
	if err != nil {
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "Failed to get pump rank",
		}
	}

	// 吧jsonResult转换为string
	jsonResultString := string(jsonResult)

	// 把jsonResult转换为 es.SearchTokenResult
	var tokens []es.TokenInfo
	err = json.Unmarshal([]byte(jsonResultString), &tokens)
	if err != nil {
		log.Println("err", err)
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "Failed to get pump rank",
		}
	}

	// 拿到代币的token_address列表
	tokenAddresses := make([]string, len(tokens))
	for i, tokenInfo := range tokens {
		tokenAddresses[i] = tokenInfo.TokenAddress
	}

	// 根据token_address列表获取代币信息
	tokenInfos := service.GetExistingTokenInfos(tokenAddresses, chainType)
	if len(tokenInfos) == 0 {
		return &response.Response{
			Code: http.StatusNotFound,
			Msg:  "Token not found",
		}
	}

	var tokenInfoResponses []*response.TokenInfoResponse

	for _, token := range tokenInfos {
		if token.ExtInfo == "" && token.Liquidity.IsZero() {
			delete(tokenInfos, token.TokenAddress)
			continue
		}

		var extInfo model.ExtInfo
		err = json.Unmarshal([]byte(token.ExtInfo), &extInfo)
		if err != nil {
			// 从tokenInfos删除这个token
			delete(tokenInfos, token.TokenAddress)
			continue
		}

		tokenInfoResponses = append(tokenInfoResponses, &response.TokenInfoResponse{
			Name:                extInfo.Name,
			Logo:                extInfo.Image,
			Price:               token.Price.InexactFloat64(),
			Symbol:              extInfo.Symbol,
			Address:             token.TokenAddress,
			Liquidity:           token.Liquidity,
			Decimals:            int(token.Decimals),
			MarketCap:           token.MarketCap.InexactFloat64(),
			CreatedPlatformType: uint8(token.CreatedPlatformType),
		})

	}

	// tokenInfos 大于0 则按流动性排序
	if len(tokenInfoResponses) > 0 {
		sort.Slice(tokenInfoResponses, func(i, j int) bool {
			return tokenInfoResponses[i].Liquidity.GreaterThan(tokenInfoResponses[j].Liquidity)
		})
	}

	// 如果 tokenInfos 大于10, 则只返回10个
	if len(tokenInfoResponses) > 20 {
		tokenInfoResponses = tokenInfoResponses[:20]
	}

	tokenAddressList := make([]string, len(tokenInfoResponses))
	for i, tokenInfoResponse := range tokenInfoResponses {
		tokenAddressList[i] = tokenInfoResponse.Address
	}

	// 用token_address列表，去es查询对应代币的24小时交易数据
	queryJSON, err = es.SearchToken(tokenAddressList, chainType)
	if err != nil {
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "Failed to get token info",
		}
	}

	tokenTransactionsResult, err := es.SearchTokenTransactionsWithAggs(es.ES_INDEX_TOKEN_TRANSACTIONS_ALIAS, queryJSON, es.UNIQUE_TOKENS)
	if err != nil {
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "Failed to get token info",
		}
	}

	aggregationResult, err := es.UnmarshalAggregationResult(tokenTransactionsResult)
	if err != nil {
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "Failed to get pump rank",
		}
	}
	if len(aggregationResult.Buckets) == 0 {
		// 如果没查出来数据，表示懂没有交易 24小时交易了都为0
		for _, tokenInfoResponse := range tokenInfoResponses {
			tokenInfoResponse.Volume = 0
		}
	} else {
		for _, bucket := range aggregationResult.Buckets {
			if len(bucket.LatestTransaction.Hits.Hits) > 0 {
				var tokenTransaction response.TokenTransaction
				if err := json.Unmarshal(bucket.LatestTransaction.Hits.Hits[0].Source, &tokenTransaction); err != nil {
					return &response.Response{
						Code: http.StatusInternalServerError,
						Msg:  "Failed to get pump rank",
					}
				}
				for _, tokenInfoResponse := range tokenInfoResponses {
					if tokenInfoResponse.Address == tokenTransaction.TokenAddress {
						// 除以decimals次方
						tokenInfoResponse.Volume = bucket.Volume.Value * tokenTransaction.Price / math.Pow(10, float64(tokenInfoResponse.Decimals))
					}
				}
			}
		}
	}

	return &response.Response{
		Code: http.StatusOK,
		Msg:  "Token search queried successfully",
		Data: tokenInfoResponses,
	}
}

func (service *TokenInfoService) GetTokenMarketAnalytics(tokenAddress string, chainType uint8) *response.Response {

	queryJSON, err := es.TokenMarketAnalyticsQuery(tokenAddress, chainType)
	if queryJSON == "" {
		return &response.Response{
			Code: http.StatusNotFound,
			Msg:  "token not found",
		}
	}

	if err != nil {
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "failed to search token",
		}
	}

	result, err := es.SearchTokenTransactionsWithAggs(es.ES_INDEX_TOKEN_TRANSACTIONS_ALIAS, queryJSON, es.UNIQUE_TOKENS)
	if result == nil {
		return &response.Response{
			Code: http.StatusNotFound,
			Msg:  "token not found",
		}
	}

	aggregationResult, err := es.UnmarshalAggregationResult(result)

	if err != nil {
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "Failed to get pump rank",
		}
	}
	if len(aggregationResult.Buckets) == 0 {
		return &response.Response{
			Code: http.StatusOK,
			Msg:  "No data found",
		}
	}

	buyVolume1m := decimal.NewFromInt(0)
	sellVolume1m := decimal.NewFromInt(0)
	buyVolume5m := decimal.NewFromInt(0)
	sellVolume5m := decimal.NewFromInt(0)
	buyVolume1h := decimal.NewFromInt(0)
	sellVolume1h := decimal.NewFromInt(0)
	buyVolume24h := decimal.NewFromInt(0)
	sellVolume24h := decimal.NewFromInt(0)

	buyCount1m := decimal.NewFromInt(0)
	sellCount1m := decimal.NewFromInt(0)
	buyCount5m := decimal.NewFromInt(0)
	sellCount5m := decimal.NewFromInt(0)
	buyCount1h := decimal.NewFromInt(0)
	sellCount1h := decimal.NewFromInt(0)
	buyCount24h := decimal.NewFromInt(0)
	sellCount24h := decimal.NewFromInt(0)

	price := float64(0)
	price1m := float64(0)
	price5m := float64(0)
	price1h := float64(0)
	price24h := float64(0)

	nativePrice := float64(0)
	solPrice := float64(0)
	priceChange1m := float64(0)
	priceChange5m := float64(0)
	priceChange1h := float64(0)
	priceChange24h := float64(0)
	var decimals int

	for _, bucket := range aggregationResult.Buckets {

		if len(bucket.LastTransactionPrice.Latest.Hits.Hits) > 0 {
			price = bucket.LastTransactionPrice.Latest.Hits.Hits[0].Source.Price
			decimals = bucket.LastTransactionPrice.Latest.Hits.Hits[0].Source.Decimals
			nativePrice = bucket.LastTransactionPrice.Latest.Hits.Hits[0].Source.NativePrice

			// sol 的价格
			solPrice = price / nativePrice

			decimals = response.SolDecimals

		} else {
			decimals = 0
			price = 0
			continue
		}
		if len(bucket.LastTransaction1mPrice.Latest.Hits.Hits) > 0 {
			price1m = bucket.LastTransaction1mPrice.Latest.Hits.Hits[0].Source.Price
		} else {
			price1m = 0
		}
		if len(bucket.LastTransaction5mPrice.Latest.Hits.Hits) > 0 {
			price5m = bucket.LastTransaction5mPrice.Latest.Hits.Hits[0].Source.Price
		}
		if len(bucket.LastTransaction1hPrice.Latest.Hits.Hits) > 0 {
			price1h = bucket.LastTransaction1hPrice.Latest.Hits.Hits[0].Source.Price
		}
		if len(bucket.LastTransaction24hPrice.Latest.Hits.Hits) > 0 {
			price24h = bucket.LastTransaction24hPrice.Latest.Hits.Hits[0].Source.Price
		}

		// 计算价格变化 也要先判断价格是否为0
		if price != 0 && price1m != 0 {
			priceChange1m = (price - price1m) / price1m
		}
		if price != 0 && price5m != 0 {
			priceChange5m = (price - price5m) / price5m
		}
		if price != 0 && price1h != 0 {
			priceChange1h = (price - price1h) / price1h
		}
		if price != 0 && price24h != 0 {
			priceChange24h = (price - price24h) / price24h
		}

		if bucket.BuyVolume1m.TotalVolume.Value > 0 {
			buyVolume1m = decimal.NewFromFloat(bucket.BuyVolume1m.TotalVolume.Value).Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(decimals)))).Mul(decimal.NewFromFloat(solPrice))
		}
		if bucket.SellVolume1m.TotalVolume.Value > 0 {
			sellVolume1m = decimal.NewFromFloat(bucket.SellVolume1m.TotalVolume.Value).Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(decimals)))).Mul(decimal.NewFromFloat(solPrice))
		}
		if bucket.BuyVolume5m.TotalVolume.Value > 0 {
			buyVolume5m = decimal.NewFromFloat(bucket.BuyVolume5m.TotalVolume.Value).Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(decimals)))).Mul(decimal.NewFromFloat(solPrice))
		}
		if bucket.SellVolume5m.TotalVolume.Value > 0 {
			sellVolume5m = decimal.NewFromFloat(bucket.SellVolume5m.TotalVolume.Value).Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(decimals)))).Mul(decimal.NewFromFloat(solPrice))
		}
		if bucket.BuyVolume1h.TotalVolume.Value > 0 {
			buyVolume1h = decimal.NewFromFloat(bucket.BuyVolume1h.TotalVolume.Value).Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(decimals)))).Mul(decimal.NewFromFloat(solPrice))
		}
		if bucket.SellVolume1h.TotalVolume.Value > 0 {
			sellVolume1h = decimal.NewFromFloat(bucket.SellVolume1h.TotalVolume.Value).Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(decimals)))).Mul(decimal.NewFromFloat(solPrice))
		}
		if bucket.BuyVolume24h.TotalVolume.Value > 0 {
			buyVolume24h = decimal.NewFromFloat(bucket.BuyVolume24h.TotalVolume.Value).Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(decimals)))).Mul(decimal.NewFromFloat(solPrice))
		}
		if bucket.SellVolume24h.TotalVolume.Value > 0 {
			sellVolume24h = decimal.NewFromFloat(bucket.SellVolume24h.TotalVolume.Value).Div(decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(decimals)))).Mul(decimal.NewFromFloat(solPrice))
		}

		if bucket.BuyCount1m.BuyVolume.Value > 0 {
			buyCount1m = decimal.NewFromInt(bucket.BuyCount1m.BuyVolume.Value)
		}
		if bucket.SellCount1m.SellVolume.Value > 0 {
			sellCount1m = decimal.NewFromInt(bucket.SellCount1m.SellVolume.Value)
		}
		if bucket.BuyCount5m.BuyVolume.Value > 0 {
			buyCount5m = decimal.NewFromInt(bucket.BuyCount5m.BuyVolume.Value)
		}
		if bucket.SellCount5m.SellVolume.Value > 0 {
			sellCount5m = decimal.NewFromInt(bucket.SellCount5m.SellVolume.Value)
		}
		if bucket.BuyCount1h.BuyVolume.Value > 0 {
			buyCount1h = decimal.NewFromInt(bucket.BuyCount1h.BuyVolume.Value)
		}
		if bucket.SellCount1h.SellVolume.Value > 0 {
			sellCount1h = decimal.NewFromInt(bucket.SellCount1h.SellVolume.Value)
		}
		if bucket.BuyCount24h.BuyVolume.Value > 0 {
			buyCount24h = decimal.NewFromInt(bucket.BuyCount24h.BuyVolume.Value)
		}
		if bucket.SellCount24h.SellVolume.Value > 0 {
			sellCount24h = decimal.NewFromInt(bucket.SellCount24h.SellVolume.Value)
		}
	}

	var analytics response.TokenMarketAnalyticsResponse

	analytics.TokenAddress = tokenAddress

	analytics.BuyVolume1m = buyVolume1m
	analytics.SellVolume1m = sellVolume1m
	analytics.BuyVolume5m = buyVolume5m
	analytics.SellVolume5m = sellVolume5m
	analytics.BuyVolume1h = buyVolume1h
	analytics.SellVolume1h = sellVolume1h
	analytics.BuyVolume24h = buyVolume24h
	analytics.SellVolume24h = sellVolume24h

	analytics.Volume1m = analytics.BuyVolume1m.Add(analytics.SellVolume1m)
	analytics.Volume5m = analytics.BuyVolume5m.Add(analytics.SellVolume5m)
	analytics.Volume1h = analytics.BuyVolume1h.Add(analytics.SellVolume1h)
	analytics.Volume24h = analytics.BuyVolume24h.Add(analytics.SellVolume24h)

	analytics.TotalCount1m = analytics.BuyCount1m.Add(analytics.SellCount1m)
	analytics.TotalCount5m = analytics.BuyCount5m.Add(analytics.SellCount5m)
	analytics.TotalCount1h = analytics.BuyCount1h.Add(analytics.SellCount1h)
	analytics.TotalCount24h = analytics.BuyCount24h.Add(analytics.SellCount24h)

	analytics.BuyCount1m = buyCount1m
	analytics.BuyCount5m = buyCount5m
	analytics.BuyCount1h = buyCount1h
	analytics.BuyCount24h = buyCount24h

	analytics.SellCount1m = sellCount1m
	analytics.SellCount5m = sellCount5m
	analytics.SellCount1h = sellCount1h
	analytics.SellCount24h = sellCount24h

	analytics.PriceChange1m = priceChange1m
	analytics.PriceChange5m = priceChange5m
	analytics.PriceChange1h = priceChange1h
	analytics.PriceChange24h = priceChange24h

	analytics.CurrentPrice = price

	return &response.Response{
		Code: http.StatusOK,
		Msg:  "Token market analytics queried successfully",
		Data: analytics,
	}
}

func (service *TokenInfoService) GetTokenBaseInfo(tokenAddress string, chainType uint8) *response.Response {
	tokenInfo, err := model.GetTokenInfoByAddress(tokenAddress, chainType)
	if err != nil {
		return &response.Response{
			Code: http.StatusNotFound,
			Msg:  "Token not found",
		}
	}
	if tokenInfo == nil {
		return &response.Response{
			Code: http.StatusNotFound,
			Data: nil,
			Msg:  "Token not supported",
		}
	}
	var tokenBaseInfo model.TokenBaseInfo
	tokenBaseInfo.Address = tokenInfo.TokenAddress
	tokenBaseInfo.Symbol = tokenInfo.Symbol
	tokenBaseInfo.URI = tokenInfo.ExtInfo
	tokenBaseInfo.Name = tokenInfo.TokenName
	tokenBaseInfo.Creator = tokenInfo.Creator
	tokenBaseInfo.Decimals = int(tokenInfo.Decimals)
	tokenBaseInfo.ChainType = tokenInfo.ChainType
	tokenBaseInfo.CreatedPlatformType = tokenInfo.CreatedPlatformType
	tokenBaseInfo.IsComplete = tokenInfo.IsComplete
	platformType := uint8(1)
	if tokenInfo.IsComplete {
		platformType = 2
	}
	pool, err := QueryAndCheckPool(tokenAddress, chainType, platformType)
	if err != nil {
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "Failed to get pool info",
		}
	}
	tokenBaseInfo.PoolAddress = pool.PoolAddress

	// 先检查 ExtInfo 是否为空
	if tokenInfo.ExtInfo != "" {
		var extInfo model.ExtInfo
		err = json.Unmarshal([]byte(tokenInfo.ExtInfo), &extInfo)
		if err != nil {
			return &response.Response{
				Code: http.StatusInternalServerError,
				Msg:  "Failed to parse token ext info",
			}
		}

		// 赋值前检查每个字段
		if extInfo.Website != "" {
			tokenBaseInfo.Website = extInfo.Website
		}
		if extInfo.Twitter != "" {
			tokenBaseInfo.Twitter = extInfo.Twitter
		}
		if extInfo.Telegram != "" {
			tokenBaseInfo.Telegram = extInfo.Telegram
		}
		if extInfo.Image != "" {
			tokenBaseInfo.URI = extInfo.Image
		}
	}

	return &response.Response{
		Code: http.StatusOK,
		Msg:  "Token base info queried successfully",
		Data: tokenBaseInfo,
	}
}

func (service *TokenInfoService) GetTokenCheckInfo(tokenAddress string, chainType uint8, tokenPool string) *response.Response {

	// 判断参数是否正常
	if tokenAddress == "" || chainType == 0 || tokenPool == "" {
		return &response.Response{
			Code: http.StatusBadRequest,
			Msg:  "Invalid parameters",
		}
	}

	redisKey := constants.RedisKeyHotTokens

	// 1. 先检查 key 是否存在
	exists, err := redis.Exists(redisKey)
	if err != nil {
		util.Log().Error("检查 Redis key 是否存在失败: %v", err)
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "Failed to check Redis key",
		}
	}

	// 2. 如果 key 不存在，先创建
	if !exists {
		// key 的 ttl 0 表示永久存储
		keyTTL := 0
		util.Log().Info("Redis key 不存在，创建新的 Sorted Set")
		err = redis.CreateSortedSet(redisKey, int64(keyTTL))
		if err != nil {
			util.Log().Error("创建 Sorted Set 失败: %v", err)
			return &response.Response{
				Code: http.StatusInternalServerError,
				Msg:  "Failed to create Sorted Set",
			}
		}
	}

	// 3. 安全清理过期数据
	err = redis.SafeCleanExpiredTokens(context.Background(), redisKey)
	if err != nil {
		util.Log().Error("清理过期数据失败: %v", err)
		// 不要直接返回，继续执行
	}

	var tokenInfo model.TokenInfo

	// 是否需要调用链端
	var dexscrAd bool
	var dexscrUpdateLink bool
	var top10HolderRate float64

	lockKey := fmt.Sprintf("token_check_info_lock_%s_%s", tokenAddress, tokenPool)
	lockValue := uuid.New().String()

	// 查询redis
	isTokenValid, err := redis.IsTokenValid(context.Background(), redisKey, tokenAddress)

	if err != nil {
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "Failed to check token validity",
		}
	}

	// 查询token_info
	token, err := model.GetTokenInfoByAddress(tokenAddress, chainType)
	if err != nil {
		return &response.Response{
			Code: http.StatusNotFound,
			Msg:  "Token not found",
		}
	}

	if token.CreatedPlatformType == 1 {
		token.PoolAddress = ""
	}

	tokenInfo = *token

	if !isTokenValid {
		// 加分布式锁
		lock, err := redis.Lock(lockKey, lockValue, 10*time.Second, 1*time.Second)
		if err != nil {
			return &response.Response{
				Code: http.StatusInternalServerError,
				Msg:  "Failed to lock",
			}
		}
		if !lock {
			return &response.Response{
				Code: http.StatusInternalServerError,
				Msg:  "Failed to lock",
			}
		}

		// 2. 再次查询redis
		isTokenValid, err = redis.IsTokenValid(context.Background(), redisKey, tokenAddress)

		if err != nil {
			return &response.Response{
				Code: http.StatusInternalServerError,
				Msg:  "Failed to check token validity",
			}
		}

		// 如果查到了，则重新查询数据库
		if isTokenValid {
			token, err := model.GetTokenInfoByAddress(tokenAddress, chainType)
			if err != nil {
				return &response.Response{
					Code: http.StatusNotFound,
					Msg:  "Token not found",
				}
			}
			tokenInfo = *token
			return &response.Response{
				Code: http.StatusInternalServerError,
				Msg:  "Failed to unmarshal token info",
				Data: InitializeTokenCheckPoolResponse(tokenAddress, tokenPool, tokenInfo),
			}
		}

		// 调用接口检测
		var wg sync.WaitGroup

		// 用于存储结果
		var safetyCheckData []httpRespone.SafetyCheckPoolData
		var safetyCheckErr error

		wg.Add(1)
		go func() {
			// 添加 defer recover
			defer func() {
				if r := recover(); r != nil {
					util.Log().Error("Safety check panic: %v\nStack: %s", r, debug.Stack())
					safetyCheckErr = fmt.Errorf("panic in safety check: %v", r)
				}
				wg.Done()
			}()

			var tokens []map[string]string
			tokens = append(tokens, map[string]string{
				"mints":         tokenAddress,
				"poolAddresses": tokenPool,
			})

			// 获取安全检查数据
			data, err := httpUtil.GetSafetyCheckPool(tokens)
			if err != nil {
				util.Log().Error("GetSafetyCheckPool failed: %v", err)
				safetyCheckErr = err
				return
			}

			// 空值检查
			if data == nil {
				util.Log().Error("GetSafetyCheckPool returned nil")
				safetyCheckErr = errors.New("safety check data is nil")
				return
			}

			// 安全地访问数据
			safetyCheckData = *data
			if len(safetyCheckData) > 0 {
				// 使用互斥锁保护 tokenInfo 的更新
				if safetyCheckData[0].Holders > 0 {
					tokenInfo.Holder = safetyCheckData[0].Holders
				}

				if safetyCheckData[0].LpBurnedPercentage > 0 {
					tokenInfo.BurnPercentage = safetyCheckData[0].LpBurnedPercentage
				}

				// 安全地计算持有比例
				circulatingSupply := float64(tokenInfo.CirculatingSupply)
				if tokenInfo.Decimals > 0 {
					circulatingSupply = circulatingSupply / math.Pow(10, float64(tokenInfo.Decimals))
				}

				if circulatingSupply > 0 {
					percentage := float64(safetyCheckData[0].Top10Holdings) / circulatingSupply
					top10HolderRate = percentage
					tokenInfo.Top10Percentage = percentage
				}
			}
		}()

		// 等待所有 goroutine 完成
		wg.Wait()

		// 检查是否发生错误
		if safetyCheckErr != nil {
			util.Log().Error("Safety check failed: %v", safetyCheckErr)
			// 可以选择继续执行或返回错误
			return &response.Response{
				Code: http.StatusInternalServerError,
				Msg:  fmt.Sprintf("Safety check failed: %v", safetyCheckErr),
			}
		}

		if dexscrUpdateLink {
			tokenInfo.SetFlag(model.FLAG_DEXSCR_UPDATE)
		}
		if dexscrAd {
			tokenInfo.SetFlag(model.FLAG_DXSCR_AD)
		}

		tokenInfo.Top10Percentage = top10HolderRate
		model.UpdateTokenInfo(&tokenInfo)
		redis.AddToken(context.Background(), redisKey, tokenAddress, time.Now().Add(2*time.Hour).Unix())
		defer redis.Unlock(lockKey, lockValue)
	}

	return &response.Response{
		Code: http.StatusOK,
		Msg:  "Token check info queried successfully",
		Data: InitializeTokenCheckPoolResponse(tokenAddress, tokenPool, tokenInfo),
	}
}

func InitializeTokenCheckPoolResponse(tokenAddress string, tokenPool string, tokenInfo model.TokenInfo) response.TokenCheckPool {
	var tokenCheckPoolResponse response.TokenCheckPool
	tokenCheckPoolResponse.TokenAddress = tokenAddress
	tokenCheckPoolResponse.PoolAddress = tokenPool

	// Extracted logic
	tokenCheckPoolResponse.CtoFlag = tokenInfo.HasFlag(model.FLAG_CTO)
	tokenCheckPoolResponse.DexscrAd = tokenInfo.HasFlag(model.FLAG_DXSCR_AD)
	tokenCheckPoolResponse.DexscrUpdateLink = tokenInfo.HasFlag(model.FLAG_DEXSCR_UPDATE)
	tokenCheckPoolResponse.MintAuthority = tokenInfo.HasFlag(model.FLAG_MINT_AUTHORITY)
	tokenCheckPoolResponse.FreezeAuthority = tokenInfo.HasFlag(model.FLAG_FREEZE_AUTHORITY)
	if tokenInfo.CreatedPlatformType == uint8(model.CreatedPlatformTypePump) {
		tokenCheckPoolResponse.IsBurnedLp = true
	} else {
		tokenCheckPoolResponse.IsBurnedLp = tokenInfo.HasFlag(model.FLAG_BURNED_LP)
	}
	tokenCheckPoolResponse.DevStatus = tokenInfo.DevStatus
	tokenCheckPoolResponse.Top10HolderRate = tokenInfo.Top10Percentage
	tokenCheckPoolResponse.Holders = tokenInfo.Holder
	tokenCheckPoolResponse.LpBurnedPercentage = tokenInfo.BurnPercentage

	return tokenCheckPoolResponse
}

func (service *TokenInfoService) GetTokenMarketQuery(tokenAddress string, chainType uint8) *response.Response {
	queryJSON, err := es.TokenAnalyticsQuery(tokenAddress, chainType)
	if err != nil {
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "failed to get token market query",
		}
	}

	result, err := es.SearchTokenTransactionsWithAggs(es.ES_INDEX_TOKEN_TRANSACTIONS_ALIAS, queryJSON, "unique_tokens")
	if err != nil {
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "Failed to get pump rank",
		}
	}

	aggregationResult, err := es.UnmarshalAggregationResult(result)
	if err != nil {
		return &response.Response{
			Code: http.StatusInternalServerError,
			Msg:  "Failed to get pump rank",
		}
	}

	price1h := 0.0
	price4h := 0.0
	price := 0.0
	isComplete := false
	createdPlatformType := 0
	var tokenTransactions []*response.TokenTransaction

	for _, bucket := range aggregationResult.Buckets {
		if len(bucket.LatestTransaction.Hits.Hits) > 0 {
			var tokenTransaction response.TokenTransaction
			if err := json.Unmarshal(bucket.LatestTransaction.Hits.Hits[0].Source, &tokenTransaction); err != nil {
				return &response.Response{
					Code: http.StatusInternalServerError,
					Msg:  "Failed to get pump rank",
				}
			}
			tokenTransactions = append(tokenTransactions, &tokenTransaction)
		}

		if len(bucket.LastTransaction1hPrice.Latest.Hits.Hits) > 0 {
			price1h = bucket.LastTransaction1hPrice.Latest.Hits.Hits[0].Source.Price

		}

		if len(bucket.LastTransaction4hPrice.Latest.Hits.Hits) > 0 {
			price4h = bucket.LastTransaction4hPrice.Latest.Hits.Hits[0].Source.Price
		}
	}

	if len(tokenTransactions) != 0 {
		tokenTransaction := *tokenTransactions[0]
		price = tokenTransaction.Price
		isComplete = tokenTransaction.IsComplete
		createdPlatformType = tokenTransaction.CreatedPlatformType
	} else {
		// 查询数据库
		tokenInfo, err := model.GetTokenInfoByAddress(tokenAddress, chainType)
		if err != nil {
			return &response.Response{
				Code: http.StatusNotFound,
				Msg:  "Token not found",
			}
		}
		priceFloat64, success := tokenInfo.Price.Float64()
		if !success {
			return &response.Response{
				Code: http.StatusInternalServerError,
				Msg:  "Failed to convert token price",
			}
		}
		price = priceFloat64
		createdPlatformType = int(tokenInfo.CreatedPlatformType)
	}

	var tokenAnalytisResponse response.TokenAnalytisResponse
	tokenAnalytisResponse.TokenAddress = tokenAddress
	tokenAnalytisResponse.Price = price
	tokenAnalytisResponse.CreatedPlatformType = uint8(createdPlatformType)
	tokenAnalytisResponse.IsComplete = isComplete

	if price1h != 0 {
		// 变成百分比
		tokenAnalytisResponse.PriceChange1h = (price - price1h) / price1h
	} else {
		tokenAnalytisResponse.PriceChange1h = 0
	}
	if price4h != 0 {
		tokenAnalytisResponse.PriceChange4h = (price - price4h) / price4h
	} else {
		tokenAnalytisResponse.PriceChange4h = 0
	}

	return &response.Response{
		Code: http.StatusOK,
		Msg:  "Token market query queried successfully",
		Data: tokenAnalytisResponse,
	}
}

func isBase58(s string) bool {
	base58Chars := "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	for _, c := range s {
		if !strings.ContainsRune(base58Chars, c) {
			return false
		}
	}
	return true
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
		Index(es.ES_INDEX_TOKEN_INFO).
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
