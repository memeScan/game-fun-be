package model

import (
	"time"

	"gorm.io/gorm"
)

type PointRecordsRepo struct {
	db *gorm.DB
}

func NewPointRecordsRepo() *PointRecordsRepo {
	return &PointRecordsRepo{}
}

// PointRecords 积分记录表
type PointRecords struct {
	ID                uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID            uint      `gorm:"column:user_id;not null" json:"user_id"`
	PointsChange      uint64    `gorm:"column:points_change;type:bigint unsigned;default:0" json:"points_change"`
	PointsBalance     uint64    `gorm:"column:points_balance;type:bigint unsigned;default:0" json:"points_balance"`
	RecordType        int8      `gorm:"column:record_type;type:tinyint;not null" json:"record_type"`
	InviteeID         uint      `gorm:"column:invitee_id;default null" json:"invitee_id,omitempty"`
	NativeTokenAmount uint64    `gorm:"column:native_token_amount;type:bigint unsigned;default:0" json:"native_token_amount"`
	TokenAmount       uint64    `gorm:"column:token_amount;type:bigint unsigned;default:0" json:"token_amount"`
	TransactionHash   string    `gorm:"column:transaction_hash;type:varchar(88)" json:"transaction_hash,omitempty"`
	TransactionDetail string    `gorm:"column:transaction_detail;type:text" json:"transaction_detail,omitempty"`
	CreateTime        time.Time `gorm:"column:create_time;type:datetime" json:"create_time"`
	UpdateTime        time.Time `gorm:"column:update_time;type:datetime" json:"update_time"`
}

type InvitedPointsDetail struct {
	UserID uint   `json:"user_id"`
	Points uint64 `json:"points"`
}

// TableName 返回表名
func (PointRecords) TableName() string {
	return "point_records"
}

// CreatePointRecord 创建积分记录
func (r *PointRecordsRepo) CreatePointRecord(record *PointRecords) error {
	return DB.Create(record).Error
}

// CreatePointRecords 创建积分记录
func (r *PointRecordsRepo) CreatePointRecords(records []*PointRecords) error {
	return DB.Create(records).Error
}

// GetPointRecordsByUserIDWithCursor retrieves point records for a user with cursor-based pagination.
// Returns records, next cursor (if any), hasMore flag, and error
func (r *PointRecordsRepo) GetPointRecordsByUserIDWithCursor(userID uint64, cursor *uint, limit int) ([]*PointRecords, *uint, bool, error) {
	var records []*PointRecords
	query := DB.Where("user_id = ?", userID).Order("id desc").Limit(limit + 1) // Request one extra record

	if cursor != nil && *cursor != 0 {
		query = query.Where("id > ?", *cursor) // Changed from > to < for desc order
	}

	if err := query.Find(&records).Error; err != nil {
		return nil, nil, false, err
	}

	hasMore := false
	if len(records) == limit {
		hasMore = true
		records = records[:limit] // Remove the extra record
	}

	var nextCursor *uint
	if len(records) > 0 {
		nextCursor = &records[len(records)-1].ID
	}

	return records, nextCursor, hasMore, nil
}

func (t *PointRecordsRepo) InvitedPointsDetail(userIDs []uint) ([]*InvitedPointsDetail, error) {
	query := DB.Table("point_records").Where("user_id in ? and record_type = 1", userIDs).Group("user_id").Select("user_id, sum(points_change) as points")
	var records []*InvitedPointsDetail
	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

func (r *PointRecordsRepo) WithTx(tx *gorm.DB) *PointRecordsRepo {
	return &PointRecordsRepo{db: tx}
}

// BeforeCreate GORM 的钩子,在创建记录前自动设置时间
func (t *PointRecords) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	t.CreateTime = now
	t.UpdateTime = now
	return nil
}

// BeforeUpdate GORM 的钩子,在更新记录前自动更新时间
func (t *PointRecords) BeforeUpdate(tx *gorm.DB) error {
	t.UpdateTime = time.Now()
	return nil
}
