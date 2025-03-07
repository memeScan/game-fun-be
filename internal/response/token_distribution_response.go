package response

// TokenHoldersData 结构体
// @Description 代币持有者数据
type TokenDistributionResponse struct {
	// 代币持有者列表
	TokenHolders []TokenHolder `json:"token_holders"`
}

// TokenHolder 表示单个代币持有者
// @Description 代币持有者详情
type TokenHolder struct {
	// 持有者钱包地址
	Account string `json:"account" example:"669VYcBRq51iQzFiPTQcsW2CsvLfHM9AwVmaoM1mAAR7"`
	// 持有代币的百分比
	Percentage string `json:"percentage" example:"100.000000"`
	// 是否与 Bonding Curve 关联
	IsAssociatedBondingCurve bool `json:"is_associated_bonding_curve" example:"false"`
	// 用户资料（可能为空）
	UserProfile interface{} `json:"user_profile"`
	// 持有代币的数量
	Amount string `json:"amount"`
	// 代币的 UI 显示数量
	UIAmount float64 `json:"ui_amount"`
	// 持有者的管理员信息
	Moderator Moderator `json:"moderator"`
	// 是否为社区金库
	IsCommunityVault bool `json:"is_community_vault" example:"false"`
	// 是否为黑洞地址
	IsBlackHole bool `json:"is_black_hole" example:"false"`
}

// Moderator 表示持有者的管理员信息
// @Description 持有者的管理员信息
type Moderator struct {
	// 被封禁的 Moderator ID
	BannedModID int `json:"banned_mod_id" example:"0"`
	// 当前状态
	Status string `json:"status" example:"NORMAL"`
	// 是否被封禁
	Banned bool `json:"banned" example:"false"`
}
