package request

// CreateTokenConfigRequest 创建代币配置请求
type CreateTokenConfigRequest struct {
	Name            string `json:"name" binding:"required"`
	Symbol          string `json:"symbol" binding:"required"`
	Address         string `json:"address" binding:"required"`
	EnableMining    bool   `json:"enable_mining"`
	MiningStartTime string `json:"mining_start_time"`
	MiningEndTime   string `json:"mining_end_time"`
	IsListed        bool   `json:"isListed"`
	Description     string `json:"description"`
}

// UpdateTokenConfigRequest 更新代币配置请求
type UpdateTokenConfigRequest struct {
	Name            string `json:"name"`
	Symbol          string `json:"symbol"`
	Address         string `json:"address"`
	EnableMining    bool   `json:"enable_mining"`
	MiningStartTime string `json:"mining_start_time"`
	MiningEndTime   string `json:"mining_end_time"`
	IsListed        bool   `json:"isListed"`
	Description     string `json:"description"`
}