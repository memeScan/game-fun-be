package model

import (
	"time"

	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TokenLiquidityPool 代币流动性池模型
type TokenLiquidityPool struct {
	ID                 uint64    `gorm:"column:id;primaryKey;autoIncrement;type:bigint unsigned;not null" json:"id"`
	ChainType          uint8     `gorm:"column:chain_type;type:tinyint;not null" json:"chain_type"`
	PlatformType       uint8     `gorm:"column:platform_type;type:tinyint unsigned" json:"platform_type"`
	MarketAddress      string    `gorm:"column:market_address;type:varchar(64)" json:"market_address"`
	PoolAddress        string    `gorm:"column:pool_address;type:varchar(64);not null" json:"pool_address"`
	PairHash           uint64    `gorm:"column:pair_hash;type:bigint unsigned" json:"pair_hash"`
	PcAddress          string    `gorm:"column:pc_address;type:varchar(64)" json:"pc_address"`
	CoinAddress        string    `gorm:"column:coin_address;type:varchar(64)" json:"coin_address"`
	PoolPcReserve      uint64    `gorm:"column:pool_pc_reserve;type:bigint unsigned;default:0" json:"pool_pc_reserve"`
	PoolCoinReserve    uint64    `gorm:"column:pool_coin_reserve;type:bigint unsigned;default:0" json:"pool_coin_reserve"`
	RealNativeReserves uint64    `gorm:"column:real_native_reserves;type:bigint unsigned;not null;default:0" json:"real_native_reserves"`
	RealTokenReserves  uint64    `gorm:"column:real_token_reserves;type:bigint unsigned;not null;default:0" json:"real_token_reserves"`
	InitialPcReserve   uint64    `gorm:"column:initial_pc_reserve;type:bigint unsigned;default:0" json:"initial_pc_reserve"`
	InitialCoinReserve uint64    `gorm:"column:initial_coin_reserve;type:bigint unsigned;default:0" json:"initial_coin_reserve"`
	UserAddress        string    `gorm:"column:user_address;type:varchar(64)" json:"user_address"`
	Block              uint64    `gorm:"column:block;type:bigint unsigned" json:"block"`
	BlockTime          time.Time `gorm:"column:block_time;type:datetime" json:"block_time"`
	IsStandardOrder    bool      `gorm:"column:is_standard_order;type:boolean;not null;default:0" json:"is_standard_order"`
	CreateTime         time.Time `gorm:"column:create_time;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime         time.Time `gorm:"column:update_time;type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"update_time"`
}

// TableName 返回表名
func (TokenLiquidityPool) TableName() string {
	return "token_liquidity_pool"
}

// CreateTokenLiquidityPool 创建代币流动性池记录
func CreateTokenLiquidityPool(pool *TokenLiquidityPool) error {
	return DB.Create(pool).Error
}

// GetTokenLiquidityPoolByAddress 通过池子地址和链类型获取代币流动性池记录
func GetTokenLiquidityPoolByAddress(poolAddress string, chainType uint8) (*TokenLiquidityPool, error) {
	var pool TokenLiquidityPool
	err := DB.Where("pool_address = ? AND chain_type = ?", poolAddress, chainType).First(&pool).Error

	// 区分处理 "记录未找到" 和其他错误
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 记录未找到返回 nil, nil
		}
		return nil, err // 其他错误正常返回
	}

	return &pool, nil
}

// UpdateTokenLiquidityPool 更新代币流动性池记录
func UpdateTokenLiquidityPool(pool *TokenLiquidityPool) error {
	return DB.Model(&TokenLiquidityPool{}).Where("pool_address = ? AND chain_type = ?", pool.PoolAddress, pool.ChainType).Updates(pool).Error
}

// DeleteTokenLiquidityPool 删除代币流动性池记录
func DeleteTokenLiquidityPool(id uint64) error {
	return DB.Delete(&TokenLiquidityPool{}, id).Error
}

// ListTokenLiquidityPools 列出代币流动性池记录
func ListTokenLiquidityPools(limit, offset int) ([]TokenLiquidityPool, error) {
	var pools []TokenLiquidityPool
	err := DB.Limit(limit).Offset(offset).Find(&pools).Error
	return pools, err
}

// BeforeCreate GORM 的钩子,在创建记录前自动设置时间
func (t *TokenLiquidityPool) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	t.CreateTime = now
	t.UpdateTime = now
	return nil
}

// BeforeUpdate GORM 的钩子,在更新记���前自动更新时间
func (t *TokenLiquidityPool) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = time.Now()
	return nil
}

// GetTokenLiquidityPoolsByTokenAddresses 通过代币地址和交易平台类型获取代币流动性池记录
func GetTokenLiquidityPoolsByTokenAddresses(tokenAddresses []string, platformType uint8) ([]TokenLiquidityPool, error) {
	var pools []TokenLiquidityPool
	err := DB.Where("coin_address IN (?) AND platform_type = ?", tokenAddresses, platformType).Find(&pools).Error
	return pools, err
}

// UpsertTokenLiquidityPool 创建或更新代币流动性池记录
func UpsertTokenLiquidityPool(pool *TokenLiquidityPool) error {
	return DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "pool_address"}, {Name: "chain_type"}},
		UpdateAll: true,
	}).Create(pool).Error
}

// 根据pool_address获取pool信息
func GetPoolInfoByAddress(poolAddress string) (*TokenLiquidityPool, error) {
	var pool TokenLiquidityPool
	err := DB.Where("pool_address = ?", poolAddress).First(&pool).Error
	return &pool, err
}

// 根据pool_address获取pool信息
func GetPoolInfoByAddressOrderByPoolPcReserve(tokenAddress string, chainType uint8, platformType uint8) (*TokenLiquidityPool, error) {
	var pool TokenLiquidityPool
	err := DB.Where("coin_address = ? AND chain_type = ? AND platform_type = ?", tokenAddress, chainType, platformType).Order("pool_pc_reserve DESC").First(&pool).Error
	// 区分处理 "记录未找到" 和其他错误
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 记录未找到返回 nil, nil
		}
		return nil, err // 其他错误正常返回
	}

	return &pool, nil
}

// GetTokenLiquidityPoolsByAddresses 通过池子地址列表批量查询池子信息
func GetTokenLiquidityPoolsByAddresses(poolAddresses []string, chainType uint8) ([]*TokenLiquidityPool, error) {
	var pools []*TokenLiquidityPool
	err := DB.Where("chain_type = ? AND pool_address IN ?", chainType, poolAddresses).Find(&pools).Error
	return pools, err
}

// BatchUpdateTokenLiquidityPools 批量更新池子信息
func BatchUpdateTokenLiquidityPools(pools []*TokenLiquidityPool) error {
	if len(pools) == 0 {
		return nil
	}

	for _, pool := range pools {
		DB.Model(&TokenLiquidityPool{}).
			Where("pool_address = ? AND chain_type = ?", pool.PoolAddress, pool.ChainType).
			Updates(pool)
	}

	return nil
}
