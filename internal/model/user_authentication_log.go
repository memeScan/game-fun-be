package model

import (
	"time"
)

type UserAuthenticationLogRepo struct{}

func NewUserAuthenticationLogRepo() *UserAuthenticationLogRepo {
	return &UserAuthenticationLogRepo{}
}

// UserAuthenticationLog 用户授权登录记录模型
type UserAuthenticationLog struct {
	ID            uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Address       string    `gorm:"column:address;type:varchar(64);not null" json:"address"`
	Message       string    `gorm:"column:message;type:varchar(8);not null" json:"message"`
	Signature     string    `gorm:"column:signature;type:varchar(88)" json:"signature"`
	Status        int8      `gorm:"column:status;type:tinyint(1);default:0" json:"status"`
	SignatureTime time.Time `gorm:"column:signature_time;type:datetime" json:"auth_time"`
	CreateTime    time.Time `gorm:"column:create_time;type:datetime" json:"create_time"`
	UpdateTime    time.Time `gorm:"column:update_time;type:datetime" json:"update_time"`
}

// TableName 指定表名
func (UserAuthenticationLog) TableName() string {
	return "user_authentication_log"
}

func (l *UserAuthenticationLogRepo) GetUserAuthenticationLogByID(id uint) (*UserAuthenticationLog, error) {
	var log UserAuthenticationLog
	err := DB.Where("id = ?", id).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (l *UserAuthenticationLogRepo) GetUserAuthenticationLogsByAddress(address string) ([]UserAuthenticationLog, error) {
	var logs []UserAuthenticationLog
	err := DB.Where("address = ?", address).Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (l *UserAuthenticationLogRepo) IsSignatureUsed(address, signature string) (bool, error) {
	var count int64
	err := DB.Model(&UserAuthenticationLog{}).Where("address = ? AND signature = ?", address, signature).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (l *UserAuthenticationLogRepo) GetAllUserAuthenticationLogs() ([]UserAuthenticationLog, error) {
	var logs []UserAuthenticationLog
	err := DB.Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (l *UserAuthenticationLogRepo) CreateUserAuthenticationLog(log *UserAuthenticationLog) error {
	return DB.Create(log).Error
}

func (l *UserAuthenticationLogRepo) BatchCreateUserAuthenticationLogs(logs []UserAuthenticationLog) error {
	return DB.Create(&logs).Error
}

func (l *UserAuthenticationLogRepo) UpdateUserAuthenticationLogByID(id uint, updates map[string]interface{}) error {
	return DB.Model(&UserAuthenticationLog{}).Where("id = ?", id).Updates(updates).Error
}

func (l *UserAuthenticationLogRepo) UpdateUserAuthenticationLogByAddress(address string, updates map[string]interface{}) error {
	return DB.Model(&UserAuthenticationLog{}).Where("address = ?", address).Updates(updates).Error
}
