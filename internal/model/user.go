package model

import (
	"time"
)

// User 用户模型
type User struct {
    ID         uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
    Address    string    `gorm:"column:address;type:varchar(64);uniqueIndex;not null" json:"address"`
    CreateTime time.Time `gorm:"column:create_time;type:datetime" json:"create_time"`
    UpdateTime time.Time `gorm:"column:update_time;type:datetime" json:"update_time"`
}


// TableName 指定表名
func (User) TableName() string {
	return "user_info"
}

// GetUserByAddress 通过地址获取用户
func GetUserByAddress(address string) (*User, error) {
	var user User
	result := DB.Where("address = ?", address).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetOrCreateUserByAddress 根据钱包地址获取或创建用户
func GetOrCreateUserByAddress(address string) (*User, error) {
	var user User
	result := DB.Where("address = ?", address).First(&user)
	if result.Error == nil {
		return &user, nil // 用户已存在
	}
	
	// 如果没有找到用户，创建新用户
	now := time.Now()
	user = User{
		Address:    address,
		CreateTime: now,
		UpdateTime: now,
	}
	result = DB.Create(&user) // 创建新用户
	if result.Error != nil {
		return nil, result.Error // 返回创建错误
	}
	return &user, nil	
}
