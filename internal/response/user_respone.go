package response

import (
	"my-token-ai-be/internal/model"
	"time"
)

// User 用户序列化器
type User struct {
	ID         uint      `json:"id"`
	Address    string    `json:"address"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

// BuildUser 序列化用户
func BuildUser(user model.User) User {
	return User{
		ID:         user.ID,
		Address:    user.Address,
		CreateTime: user.CreateTime,
		UpdateTime: user.UpdateTime,
	}
}
