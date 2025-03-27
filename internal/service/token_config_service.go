package service

import (
	"game-fun-be/internal/model"
	"game-fun-be/internal/response"
	"net/http"
	"time"
)

// TokenConfigServiceImpl 代币配置服务实现
type TokenConfigServiceImpl struct {
	tokenInfoRepo *model.TokenInfoRepo
}

// NewTokenConfigServiceImpl 创建代币配置服务实例
func NewTokenConfigServiceImpl(tokenInfoRepo *model.TokenInfoRepo) *TokenConfigServiceImpl {
	return &TokenConfigServiceImpl{
		tokenInfoRepo: tokenInfoRepo,
	}
}

// GetTokenConfigs 获取代币配置列表
func (s *TokenConfigServiceImpl) GetTokenConfigs(page, limit int) response.Response {
	configs, total, err := model.GetTokenConfigList(page, limit)
	if err != nil {
		return response.Response{
			Code:  http.StatusInternalServerError,
			Error: "Failed to get token config list: " + err.Error(),
		}
	}

	return response.Response{
		Code: http.StatusOK,
		Data: map[string]interface{}{
			"list":  configs,
			"total": total,
			"page":  page,
			"limit": limit,
		},
		Msg: "Success",
	}
}

// GetTokenConfig 获取代币配置详情
func (s *TokenConfigServiceImpl) GetTokenConfig(id uint) response.Response {
	config, err := model.GetTokenConfigByID(id)
	if err != nil {
		return response.Response{
			Code:  http.StatusInternalServerError,
			Error: "Failed to get token config: " + err.Error(),
		}
	}

	if config == nil {
		return response.Response{
			Code:  http.StatusNotFound,
			Error: "Token config not found",
		}
	}

	return response.Response{
		Code: http.StatusOK,
		Data: config,
		Msg:  "Success",
	}
}

// CreateTokenConfig 创建代币配置
func (s *TokenConfigServiceImpl) CreateTokenConfig(name, symbol, address string, enableMining bool, miningStartTime, miningEndTime string, isListed bool, description string) response.Response {
	// 检查地址是否已经存在
	existing, err := model.GetTokenConfigByAddress(address)
	if err != nil {
		return response.Response{
			Code:  http.StatusInternalServerError,
			Error: "Failed to check token config existence: " + err.Error(),
		}
	}

	if existing != nil {
		return response.Response{
			Code:  http.StatusBadRequest,
			Error: "Token config with this address already exists",
		}
	}

	//检测代币是否在系统中存在（系统支持该代币）
	tokenInfo, err := s.tokenInfoRepo.GetTokenInfoByAddress(address, model.ChainTypeSolana.Uint8())
	if err != nil {
		return response.Err(http.StatusInternalServerError, "Failed to get token info by address", err)
	}
	if tokenInfo == nil {
		return response.Err(http.StatusNotFound, "Token info not found", nil)
	}

	if tokenInfo.Symbol != symbol {
		return response.Err(http.StatusBadRequest, "Token info not found", nil)
	}

	config := &model.TokenConfig{
		Name:        name,
		Symbol:      symbol,
		Address:     address,
		IsListed:    isListed,
		Description: description,
	}

	// 如果启用挖矿，设置相关时间
	if enableMining {
		config.EnableMining = true

		if miningStartTime == "" || miningEndTime == "" {
			return response.Err(http.StatusBadRequest, "Mining start time and end time are required", nil)
		}

		// 解析挖矿开始时间
		if miningStartTime != "" {
			startTime, err := time.Parse("2006-01-02 15:04:05", miningStartTime)
			if err == nil {
				config.MiningStartTime = &startTime
			}
		}

		// 解析挖矿结束时间
		if miningEndTime != "" {
			endTime, err := time.Parse("2006-01-02 15:04:05", miningEndTime)
			if err == nil {
				config.MiningEndTime = &endTime
			}
		}
	}

	// 创建记录
	if err := model.CreateTokenConfig(config); err != nil {
		return response.Response{
			Code:  http.StatusInternalServerError,
			Error: "Failed to create token config: " + err.Error(),
		}
	}

	return response.Response{
		Code: http.StatusOK,
		Data: config,
		Msg:  "Token config created successfully",
	}
}

// UpdateTokenConfig 更新代币配置
func (s *TokenConfigServiceImpl) UpdateTokenConfig(id uint, name, symbol, address string, enableMining bool, miningStartTime, miningEndTime string, isListed bool, description string) response.Response {
	// 获取当前配置
	config, err := model.GetTokenConfigByID(id)
	if err != nil {
		return response.Response{
			Code:  http.StatusInternalServerError,
			Error: "Failed to get token config: " + err.Error(),
		}
	}

	if config == nil {
		return response.Response{
			Code:  http.StatusNotFound,
			Error: "Token config not found",
		}
	}

	// 更新字段
	if name != "" {
		config.Name = name
	}
	if symbol != "" {
		config.Symbol = symbol
	}
	if address != "" {
		config.Address = address
	}

	config.EnableMining = enableMining
	config.IsListed = isListed
	config.Description = description

	// 如果启用挖矿，更新相关时间
	if enableMining {
		// 解析挖矿开始时间
		if miningStartTime != "" {
			startTime, err := time.Parse("2006-01-02 15:04:05", miningStartTime)
			if err == nil {
				config.MiningStartTime = &startTime
			}
		}

		// 解析挖矿结束时间
		if miningEndTime != "" {
			endTime, err := time.Parse("2006-01-02 15:04:05", miningEndTime)
			if err == nil {
				config.MiningEndTime = &endTime
			}
		}
	}

	// 保存更新
	if err := model.UpdateTokenConfig(config); err != nil {
		return response.Response{
			Code:  http.StatusInternalServerError,
			Error: "Failed to update token config: " + err.Error(),
		}
	}

	return response.Response{
		Code: http.StatusOK,
		Data: config,
		Msg:  "Token config updated successfully",
	}
}

// DeleteTokenConfig 删除代币配置
func (s *TokenConfigServiceImpl) DeleteTokenConfig(id uint) response.Response {
	// 检查配置是否存在
	config, err := model.GetTokenConfigByID(id)
	if err != nil {
		return response.Response{
			Code:  http.StatusInternalServerError,
			Error: "Failed to get token config: " + err.Error(),
		}
	}

	if config == nil {
		return response.Response{
			Code:  http.StatusNotFound,
			Error: "Token config not found",
		}
	}

	// 执行删除
	if err := model.DeleteTokenConfig(id); err != nil {
		return response.Response{
			Code:  http.StatusInternalServerError,
			Error: "Failed to delete token config: " + err.Error(),
		}
	}

	return response.Response{
		Code: http.StatusOK,
		Msg:  "Token config deleted successfully",
	}
}
