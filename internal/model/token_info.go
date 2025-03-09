package model

import (
	"time"

	"errors"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TokenInfoRepo struct{}

func NewTokenInfoRepo() *TokenInfoRepo {
	return &TokenInfoRepo{}
}

const (
	FLAG_MINT_AUTHORITY    = 1 << 0  // 第 1 位 - mint_authority
	FLAG_FREEZE_AUTHORITY  = 1 << 1  // 第 2 位 - freeze_authority
	FLAG_DXSCR_AD          = 1 << 2  // 第 3 位 - dxscr_ad
	FLAG_TWITTER_CHANGE    = 1 << 3  // 第 4 位 - twitter_change_flag
	FLAG_DEXSCR_UPDATE     = 1 << 4  // 第 5 位 - dexscr_update_link
	FLAG_CTO               = 1 << 5  // 第 6 位 - cto_flag
	FLAG_IS_MEDIA          = 1 << 6  // 第 7 位 - is_media
	FLAG_IS_COMPLETE       = 1 << 7  // 第 8 位 - is_complete
	FLAG_DEV_IS_HOLDING    = 1 << 8  // 第 9 位 - dev_is_holding_flag
	FLAG_DEV_BURN_PERCENT  = 1 << 9  // 第 10 位 - dev_burn_percentage_flag
	FLAG_DEV_IS_RUG        = 1 << 10 // 第 11 位 - dev_is_rug_flag
	FLAG_AUTHORITY_CHECKED = 1 << 11 // 第 12 位 - 是否已检查过权限
	FLAG_BURNED_LP         = 1 << 12 // 第 13 位 - 是否已烧池子
)

// TokenInfo 币种信息模型
type TokenInfo struct {
	ID                   int64           `gorm:"column:id;primaryKey;autoIncrement;type:bigint" json:"id"`
	TokenName            string          `gorm:"column:token_name;type:varchar(255)" json:"token_name"`
	Symbol               string          `gorm:"column:symbol;type:varchar(64)" json:"symbol"`
	Creator              string          `gorm:"column:creator;type:varchar(64)" json:"creator"`
	TokenAddress         string          `gorm:"column:token_address;type:varchar(64)" json:"token_address"`
	ChainType            uint8           `gorm:"column:chain_type;type:tinyint" json:"chain_type"`
	CreatedPlatformType  uint8           `gorm:"column:created_platform_type;type:tinyint unsigned" json:"created_platform_type"`
	Decimals             uint8           `gorm:"column:decimals;type:tinyint unsigned" json:"decimals"`
	TotalSupply          uint64          `gorm:"column:total_supply;type:bigint unsigned;default:0" json:"total_supply"`
	CirculatingSupply    uint64          `gorm:"column:circulating_supply;type:bigint unsigned;default:0" json:"circulating_supply"`
	Block                uint64          `gorm:"column:block;type:bigint unsigned;not null" json:"block"`
	TransactionHash      string          `gorm:"column:transaction_hash;type:varchar(88);not null" json:"transaction_hash"`
	TransactionTime      time.Time       `gorm:"column:transaction_time;type:datetime" json:"transaction_time"`
	URI                  string          `gorm:"column:uri;type:varchar(255)" json:"uri"`
	DevNativeTokenAmount uint64          `gorm:"column:dev_native_token_amount;type:bigint unsigned;default:0" json:"dev_native_token_amount"`
	DevTokenAmount       uint64          `gorm:"column:dev_token_amount;type:bigint unsigned;default:0" json:"dev_token_amount"`
	Holder               int             `gorm:"column:holder;type:int;not null" json:"holder"`
	CommentCount         int             `gorm:"column:comment_count;type:int;not null" json:"comment_count"`
	MarketCap            decimal.Decimal `gorm:"column:market_cap;type:decimal(30,10);not null" json:"market_cap"`
	CirculatingMarketCap decimal.Decimal `gorm:"column:circulating_market_cap;type:decimal(30,10);not null" json:"circulating_market_cap"`
	CrownDuration        int64           `gorm:"column:crown_duration;type:bigint;not null" json:"crown_duration"`
	RocketDuration       int64           `gorm:"column:rocket_duration;type:bigint;not null" json:"rocket_duration"`
	DevStatus            uint8           `gorm:"column:dev_status;type:tinyint;not null" json:"dev_status"`
	IsMedia              bool            `gorm:"column:is_media;type:boolean;not null" json:"is_media"`
	IsComplete           bool            `gorm:"column:is_complete;type:boolean;not null" json:"is_complete"`
	Price                decimal.Decimal `gorm:"column:price;type:decimal(30,18)" json:"price"`
	NativePrice          decimal.Decimal `gorm:"column:native_price;type:decimal(30,18)" json:"native_price"`
	Liquidity            decimal.Decimal `gorm:"column:liquidity;type:decimal(38,18)" json:"liquidity"`
	ExtInfo              string          `gorm:"column:ext_info;type:text" json:"ext_info"`
	CreateTime           time.Time       `gorm:"column:create_time;type:datetime" json:"create_time"`
	UpdateTime           time.Time       `gorm:"column:update_time;type:datetime;autoUpdateTime" json:"update_time"`
	DevPercentage        float64         `gorm:"column:dev_percentage;type:decimal(5,2);not null;default:0.00" json:"dev_percentage"`
	Top10Percentage      float64         `gorm:"column:top10_percentage;type:decimal(5,2);not null;default:0.00" json:"top10_percentage"`
	BurnPercentage       float64         `gorm:"column:burn_percentage;type:decimal(5,2);not null;default:0.00" json:"burn_percentage"`
	DevBurnPercentage    float64         `gorm:"column:dev_burn_percentage;type:decimal(5,2);not null;default:0.00" json:"dev_burn_percentage"`
	TokenFlags           int             `gorm:"column:token_flags;type:int;not null;default:0" json:"token_flags"`
	Progress             decimal.Decimal `gorm:"column:progress;type:decimal(5,2)" json:"progress"`
	PoolAddress          string          `gorm:"column:pool_address;type:varchar(64)" json:"pool_address"`
}

type ExtInfo struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
	Image       string `json:"image"`
	CreatedOn   string `json:"createdOn"`
	Twitter     string `json:"twitter"`
	Website     string `json:"website"`
	Telegram    string `json:"telegram"`
	Banner      string `json:"banner"`
	Rules       string `json:"rules"`
	Sort        uint   `json:"sort"`
	ShowName    bool   `json:"showName"`
}

func (t *TokenInfo) SetFlag(flag int) {
	t.TokenFlags |= flag
}

func (t *TokenInfo) ClearFlag(flag int) {
	t.TokenFlags &^= flag
}

func (t *TokenInfo) HasFlag(flag int) bool {
	return t.TokenFlags&flag != 0
}

func (t *TokenInfo) ToggleFlag(flag int) {
	t.TokenFlags = t.TokenFlags ^ flag
}

// TableName 返回表名
func (TokenInfo) TableName() string {
	return "token_info"
}

// CreateTokenInfo 创建币种信息记录
func CreateTokenInfo(info *TokenInfo) error {
	return DB.Create(info).Error
}

// BatchInsertTokenInfo 批量插入币种信息记录
func BatchInsertTokenInfo(infos []TokenInfo) error {
	return DB.Create(&infos).Error
}

// GetTokenInfoByAddress 通过代币地址和链类型获取币种信息
func GetTokenInfoByAddress(tokenAddress string, chainType uint8) (*TokenInfo, error) {
	var info TokenInfo
	err := DB.Where("token_address = ? AND chain_type = ?", tokenAddress, chainType).
		First(&info).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil // 未找到返回 nil, nil
	}
	return &info, nil
}

// GetTokenInfoByAddressCreateType 通过代币地址和链类型获取币种信息
func GetTokenInfoByAddressCreateType(tokenAddress string, chainType uint8, createdPlatformType uint8) (*TokenInfo, error) {
	var info TokenInfo
	err := DB.Where("token_address = ? AND chain_type = ? AND created_platform_type = ?", tokenAddress, chainType, createdPlatformType).First(&info).Error
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// GetTokenInfoByAddresses 通过代币地址和链类型获取币种信息
func GetTokenInfoByAddresses(tokenAddresses []string, chainType uint8) ([]TokenInfo, error) {
	var infos []TokenInfo
	err := DB.Where("token_address IN (?) AND chain_type = ?", tokenAddresses, chainType).Find(&infos).Error
	if err != nil {
		return nil, err
	}
	return infos, nil
}

// UpdateTokenInfo 更新币种信息记录
func UpdateTokenInfo(info *TokenInfo) error {
	return DB.Save(info).Error
}

// DeleteTokenInfo 删除币种信息记录
func DeleteTokenInfo(id int64) error {
	return DB.Delete(&TokenInfo{}, id).Error
}

// ListTokenInfos 列出币种信息记录
func ListTokenInfos(limit, offset int) ([]TokenInfo, error) {
	var infos []TokenInfo
	err := DB.Limit(limit).Offset(offset).Find(&infos).Error
	return infos, err
}

// BeforeCreate GORM 的钩子，在创建记录前自动设置时间
func (t *TokenInfo) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	t.CreateTime = now
	t.UpdateTime = now
	return nil
}

// BeforeUpdate GORM 的钩子，在更新记录前自动更新时间
func (t *TokenInfo) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = time.Now()
	return nil
}

// BatchUpdateTokenInfo 批量更新币种信息记录
func BatchUpdateTokenInfo(infos []*TokenInfo) error {
	if len(infos) == 0 {
		return nil
	}

	for _, info := range infos {
		DB.Model(&TokenInfo{}).
			Where("chain_type = ? AND token_address = ?", info.ChainType, info.TokenAddress).
			Updates(info)
	}

	return nil
}

// SearchToken 通过mysql模糊搜索币种信息，��回token_address列表
func SearchToken(name string, chainType uint8) ([]string, error) {
	var infos []TokenInfo
	err := DB.Where("token_address LIKE ? OR token_name LIKE ? AND chain_type = ?",
		"%"+name+"%",
		"%"+name+"%",
		chainType).
		Order("market_cap DESC"). // 首先按市值降序
		Limit(1).
		Find(&infos).Error
	if err != nil {
		return nil, err
	}

	var tokenAddresses []string
	for _, info := range infos {
		tokenAddresses = append(tokenAddresses, info.TokenAddress)
	}

	return tokenAddresses, nil

}

// UpdateTokenComplete 只更新代币的完成状态
func UpdateTokenComplete(tokenAddress string, chainType uint8) error {
	return DB.Model(&TokenInfo{}).
		Where("token_address = ? AND chain_type = ?", tokenAddress, chainType).
		Updates(map[string]interface{}{
			"is_complete": true,
			"update_time": time.Now(),
		}).Error
}

// ListTokenInfosByCursor 使用游标分页获取代币信息
func ListTokenInfosByCursor(lastID int64, chainType uint8, createdPlatformType uint8, isComplete bool, limit int) ([]TokenInfo, error) {
	var infos []TokenInfo
	query := DB.Where("chain_type = ? AND created_platform_type = ? AND is_complete = ?",
		chainType,
		createdPlatformType,
		isComplete)

	if lastID > 0 {
		query = query.Where("id > ?", lastID)
	}

	err := query.Order("id ASC").
		Limit(limit).
		Find(&infos).Error

	return infos, err
}

func (t *TokenInfoRepo) GetTokenInfoByAddress(tokenAddress string, chainType uint8) (*TokenInfo, error) {
	var tokenInfo TokenInfo
	err := DB.Model(&TokenInfo{}).
		Where("token_address = ? AND chain_type = ?", tokenAddress, chainType).
		First(&tokenInfo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 未找到记录，返回 nil
		}
		return nil, err // 其他错误，返回错误
	}
	return &tokenInfo, nil
}

// CreateTokenInfo 创建新的代币记录
func (t *TokenInfoRepo) CreateTokenInfo(tokenInfo *TokenInfo) error {
	// 先检查数据库中是否已存在该 token
	var existing TokenInfo
	err := DB.Model(&TokenInfo{}).
		Where("token_address = ? AND chain_type = ?", tokenInfo.TokenAddress, tokenInfo.ChainType).
		First(&existing).Error

	if err == nil {
		// 记录已存在，返回错误
		return errors.New("token already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 其他数据库错误
		return err
	}

	// 插入新记录
	return DB.Create(tokenInfo).Error
}
