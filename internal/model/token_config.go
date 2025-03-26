package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// TokenConfig 代币配置模型
type TokenConfig struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Name            string    `gorm:"column:token_name;type:varchar(255)" json:"name"`
	Symbol          string    `gorm:"column:symbol;type:varchar(64)" json:"symbol"`
	Address         string    `gorm:"column:token_address;type:varchar(64);uniqueIndex" json:"address"`
	EnableMining    bool      `gorm:"column:enable_mining;type:boolean;default:false" json:"enable_mining"`
	MiningStartTime time.Time `gorm:"column:mining_start_time;type:datetime" json:"mining_start_time"`
	MiningEndTime   time.Time `gorm:"column:mining_end_time;type:datetime" json:"mining_end_time"`
	IsListed        bool      `gorm:"column:is_listed;type:boolean;default:false" json:"is_listed"`
	Description     string    `gorm:"column:description;type:text" json:"description"`
	CreateTime      time.Time `gorm:"column:create_time;type:datetime;autoCreateTime" json:"create_time"`
	UpdateTime      time.Time `gorm:"column:update_time;type:datetime;autoUpdateTime" json:"update_time"`
}

// TableName 返回表名
func (TokenConfig) TableName() string {
	return "token_configs"
}

// BeforeCreate GORM 的钩子，在创建记录前自动设置时间
func (t *TokenConfig) BeforeCreate(tx *gorm.DB) error {
	t.CreateTime = time.Now()
	t.UpdateTime = time.Now()
	return nil
}

// BeforeUpdate GORM 的钩子，在更新记录前自动更新时间
func (t *TokenConfig) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = time.Now()
	return nil
}

// CreateTokenConfig 创建代币配置记录
func CreateTokenConfig(config *TokenConfig) error {
	return DB.Create(config).Error
}

// GetTokenConfigByID 通过ID获取代币配置
func GetTokenConfigByID(id uint) (*TokenConfig, error) {
	var config TokenConfig
	result := DB.Where("id = ?", id).First(&config)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &config, nil
}

// GetTokenConfigByAddress 通过地址获取代币配置
func GetTokenConfigByAddress(address string) (*TokenConfig, error) {
	var config TokenConfig
	result := DB.Where("token_address = ?", address).First(&config)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &config, nil
}

// UpdateTokenConfig 更新代币配置记录
func UpdateTokenConfig(config *TokenConfig) error {
	return DB.Save(config).Error
}

// DeleteTokenConfig 删除代币配置记录
func DeleteTokenConfig(id uint) error {
	return DB.Delete(&TokenConfig{}, id).Error
}

// GetTokenConfigList 获取代币配置列表，支持分页
func GetTokenConfigList(page, limit int) ([]TokenConfig, int64, error) {
	var configs []TokenConfig
	var total int64

	offset := (page - 1) * limit

	// 获取总数
	if err := DB.Model(&TokenConfig{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取当前页数据
	if err := DB.Order("id DESC").Offset(offset).Limit(limit).Find(&configs).Error; err != nil {
		return nil, 0, err
	}

	return configs, total, nil
}
