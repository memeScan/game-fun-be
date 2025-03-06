package model

import (
	"time"

	"gorm.io/gorm"
)

// PointRecords 积分记录表
type PointRecords struct {
	ID              uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID          uint      `gorm:"column:user_id;not null" json:"user_id"`
	PointsChange    float64   `gorm:"column:points_change;type:decimal(20,8);not null" json:"points_change"`
	PointsBalance   float64   `gorm:"column:points_balance;type:decimal(20,8);not null" json:"points_balance"`
	RecordType      int8      `gorm:"column:record_type;type:tinyint;not null" json:"record_type"`
	InviteeID       *uint     `gorm:"column:invitee_id" json:"invitee_id,omitempty"`
	TransactionHash string    `gorm:"column:transaction_hash;type:varchar(88)" json:"transaction_hash,omitempty"`
	Description     string    `gorm:"column:description;type:varchar(255)" json:"description,omitempty"`
	CreateTime      time.Time `gorm:"column:create_time;type:datetime" json:"create_time"`
	UpdateTime      time.Time `gorm:"column:update_time;type:datetime" json:"update_time"`
}

// TableName 返回表名
func (PointRecords) TableName() string {
	return "point_records"
}

// CreatePointRecord 创建积分记录
func CreatePointRecord(record *PointRecords) error {
	return DB.Create(record).Error
}

// GetPointRecordsByUserIDWithCursor retrieves point records for a user with cursor-based pagination.
// Returns records, next cursor (if any), hasMore flag, and error
func GetPointRecordsByUserIDWithCursor(userID uint64, cursor *uint, limit int) ([]*PointRecords, *uint, bool, error) {
	var records []*PointRecords
	query := DB.Where("user_id = ?", userID).Order("id desc").Limit(limit + 1) // Request one extra record

	if cursor != nil {
		query = query.Where("id < ?", *cursor) // Changed from > to < for desc order
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
