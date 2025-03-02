package model

import (
	"time"
)

// UserAuthenticationLog 用户授权登录记录模型
type UserAuthenticationLog struct {
	ID           uint     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Address      string    `gorm:"column:address;type:varchar(64);not null" json:"address"`
	MessageNonce string    `gorm:"column:message_nonce;type:varchar(8);not null" json:"message_nonce"`
	Signature    string    `gorm:"column:signature;type:varchar(88)" json:"signature"`
	Status       int8      `gorm:"column:status;type:tinyint(1);default:0" json:"status"`
	CreateTime   time.Time `gorm:"column:create_time;type:datetime" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:update_time;type:datetime" json:"update_time"`
}

// TableName 指定表名
func (UserAuthenticationLog) TableName() string {
	return "user_authentication_log"
}

// CreateUserAuthenticationLog 创建用户授权登录记录
func CreateUserAuthenticationLog(log *UserAuthenticationLog) error {
	return DB.Create(log).Error
}

// GetUserAuthenticationLogByAddress 通过地址获取最新的用户授权登录记录
func GetUserAuthenticationLogByAddress(address string, messageNonce string) (*UserAuthenticationLog, error) {
	var log UserAuthenticationLog
	err := DB.Where("address = ? and message_nonce = ?", address, messageNonce).Order("create_time DESC").First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}


// GetUserAuthenticationLogByAddress 通过地址获取最新的用户授权登录记录
func GetUserAuthenticationLogID(id uint) (*UserAuthenticationLog, error) {
	var log UserAuthenticationLog
	err := DB.Where("id = ?", id).First(&log).Error
	if err != nil {
		return nil, err	
	}
	return &log, nil
}

// UpdateUserAuthenticationLog 更新用户授权登录记录
func UpdateUserAuthenticationLog(log *UserAuthenticationLog) error {
	return DB.Save(log).Error
}
