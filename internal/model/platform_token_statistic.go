package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// PlatformTokenStatisticRepo 平台币种信息统计表仓库
type PlatformTokenStatisticRepo struct {
	db *gorm.DB
}

func NewPlatformTokenStatisticRepo() *PlatformTokenStatisticRepo {
	return &PlatformTokenStatisticRepo{}
}

// PlatformTokenStatistic 平台币种信息统计表
type PlatformTokenStatistic struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name          string    `gorm:"column:name;type:varchar(255);not null" json:"name"`
	TokenName     string    `gorm:"column:token_name;type:varchar(255)" json:"token_name"`
	Symbol        string    `gorm:"column:symbol;type:varchar(64)" json:"symbol"`
	TokenAddress  string    `gorm:"column:token_address;type:varchar(64)" json:"token_address"`
	FeeAmount     uint64    `gorm:"column:fee_amount;type:bigint unsigned;not null" json:"fee_amount"`
	BackAmount    uint64    `gorm:"column:back_amount;type:bigint unsigned;not null" json:"back_amount"`
	BackSolAmount uint64    `gorm:"column:back_sol_amount;type:bigint unsigned;not null" json:"back_sol_amount"`
	BurnAmount    uint64    `gorm:"column:burn_amount;type:bigint unsigned;not null" json:"burn_amount"`
	PointsAmount  uint64    `gorm:"column:points_amount;type:bigint unsigned;not null" json:"points_amount"`
	CreateTime    time.Time `gorm:"column:create_time;type:datetime" json:"create_time"`
	UpdateTime    time.Time `gorm:"column:update_time;type:datetime" json:"update_time"`
}

// TableName 返回表名
func (PlatformTokenStatistic) TableName() string {
	return "platform_token_statistics"
}

func (r *PlatformTokenStatisticRepo) CreatePlatformTokenStatistic(record *PlatformTokenStatistic) error {
	return DB.Create(record).Error
}

func (r *PlatformTokenStatisticRepo) IncrementStatisticsAndUpdateTime(address string, amounts map[StatisticType]uint64) error {
	updates := make(map[string]any)
	for pt, val := range amounts {
		updates[string(pt)] = gorm.Expr(string(pt)+" + ?", val)
	}
	updates["update_time"] = time.Now()
	return DB.Model(&PlatformTokenStatistic{}).
		Where("token_address = ?", address).
		Updates(updates).Error
}

func (s *PlatformTokenStatisticRepo) GetTokenPointsStatistic(tokenAddress string, chainType uint8) (*PlatformTokenStatistic, error) {
	if tokenAddress == "" {
		return nil, fmt.Errorf("token address cannot be empty")
	}

	var statistics PlatformTokenStatistic
	result := DB.Where("token_address = ?", tokenAddress).First(&statistics)
	if result.Error != nil {
		return nil, result.Error
	}

	return &statistics, nil
}

// BeforeCreate GORM 的钩子,在创建记录前自动设置时间
func (p *PlatformTokenStatistic) BeforeCreate(tx *gorm.DB) error {
	p.CreateTime = time.Now()
	p.UpdateTime = time.Now()
	return nil
}

// BeforeUpdate GORM 的钩子,在更新记录前自动更新时间
func (p *PlatformTokenStatistic) BeforeUpdate(tx *gorm.DB) error {
	p.UpdateTime = time.Now()
	return nil
}

func (r *PlatformTokenStatisticRepo) WithTx(tx *gorm.DB) *PlatformTokenStatisticRepo {
	return &PlatformTokenStatisticRepo{db: tx}
}
