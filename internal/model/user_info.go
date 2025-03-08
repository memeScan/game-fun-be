package model

import (
	"errors"
	"fmt"
	"game-fun-be/internal/pkg/util"
	"time"

	"gorm.io/gorm"
)

type UserInfoRepo struct{}

func NewUserInfoRepo() *UserInfoRepo {
	return &UserInfoRepo{}
}

type UserInfo struct {
	ID               uint       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Address          string     `gorm:"column:address;type:varchar(64);uniqueIndex;not null" json:"address"`
	TwitterID        string     `gorm:"column:twitter_id;type:varchar(64);omitempty" json:"twitter_id"`
	TwitterUsername  string     `gorm:"column:twitter_username;type:varchar(64);omitempty" json:"twitter_username"`
	InviterID        uint       `gorm:"column:inviter_id;type:bigint unsigned;omitempty" json:"inviter_id"`
	ParentInviteId   uint       `gorm:"column:parent_inviter_id;type:bigint unsigned;omitempty" json:"parent_inviter_id"`
	InvitationCode   string     `gorm:"column:invitation_code;type:varchar(32);uniqueIndex;not null" json:"invitation_code"`
	TradingPoints    float64    `gorm:"column:trading_points;type:decimal(20,8);not null;default:0" json:"trading_points"`
	InvitePoints     float64    `gorm:"column:invite_points;type:decimal(20,8);not null;default:0" json:"invite_points"`
	AvailablePoints  float64    `gorm:"column:available_points;type:decimal(20,8);not null;default:0" json:"avaliable_points"`
	Status           uint8      `gorm:"column:status;type:tinyint(4);not null" json:"status"`
	FirstTradingTime *time.Time `gorm:"column:first_trading_time;type:datetime;omitempty" json:"first_trading_time"`
	ChainType        uint8      `gorm:"column:chain_type;type:tinyint;omitempty" json:"chain_type"`
	CreateTime       time.Time  `gorm:"column:create_time;type:datetime;omitempty" json:"create_time"`
	UpdateTime       time.Time  `gorm:"column:update_time;type:datetime;omitempty" json:"update_time"`
}

func (UserInfo) TableName() string {
	return "user_info"
}

func (r *UserInfoRepo) getInviteUser(inviteCode string, chainType uint8) (*UserInfo, error) {
	if inviteCode == "" {
		return nil, nil
	}
	inviteUser, err := r.GetUserByInvitationCode(inviteCode, chainType)
	if err != nil {
		util.Log().Error("Failed to get user by invitation code: %v, inviteCode: %s, chainType: %d", err, inviteCode, chainType)
		return nil, err
	}
	return inviteUser, nil
}

func (r *UserInfoRepo) setInviterInfo(user *UserInfo, inviteUser *UserInfo) {
	if inviteUser != nil {
		user.InviterID = inviteUser.ID
		if inviteUser.InviterID != 0 {
			user.ParentInviteId = inviteUser.InviterID
		}
	}
}

func (r *UserInfoRepo) GetOrCreateUserByAddress(address string, chainType uint8, inviteCode string) (*UserInfo, error) {
	var user UserInfo
	result := DB.Where("address = ? AND chain_type = ?", address, chainType).First(&user)
	if result.Error == nil {
		if user.Status == 0 {
			user.Status = 1
			if user.InviterID == 0 && inviteCode != "" {
				inviteUser, err := r.getInviteUser(inviteCode, chainType)
				if err != nil {
					return nil, fmt.Errorf("failed to get user by invitation code: %v", err)
				}
				r.setInviterInfo(&user, inviteUser)
			}
			user.UpdateTime = time.Now()
			if err := DB.Save(&user).Error; err != nil {
				return nil, err
			}
		}
		return &user, nil
	}

	inviteUser, err := r.getInviteUser(inviteCode, chainType)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by invitation code: %v", err)
	}

	now := time.Now()
	user = UserInfo{
		Address:        address,
		Status:         1,
		ChainType:      chainType,
		CreateTime:     now,
		UpdateTime:     now,
		InvitationCode: util.GenerateInviteCode(address),
	}

	r.setInviterInfo(&user, inviteUser)

	result = DB.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserInfoRepo) GetInviteCodeAndCount(address string, chainType uint8) (UserInfo, int, error) {
	var user UserInfo
	result := DB.Where("address = ? AND chain_type = ?", address, chainType).First(&user)
	if result.Error != nil {
		return user, 0, result.Error
	}

	var inviteCount int64
	result = DB.Model(&UserInfo{}).Where("inviter_id = ?", user.ID).Count(&inviteCount)
	if result.Error != nil {
		return user, 0, result.Error
	}

	return user, int(inviteCount), nil
}

func (r *UserInfoRepo) GetUserByInvitationCode(inviteCode string, chainType uint8) (*UserInfo, error) {
	if inviteCode == "" {
		return nil, fmt.Errorf("invite code is empty")
	}

	var user UserInfo
	result := DB.Where("invitation_code = ? AND chain_type = ?", inviteCode, chainType).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserInfoRepo) GetUserByAddress(address string, chainType uint8) (*UserInfo, error) {
	var user UserInfo

	result := DB.Where("address = ? AND chain_type = ?", address, chainType).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found with address %s and chain type %d", address, chainType)
		}
		return nil, result.Error
	}

	return &user, nil
}

func (r *UserInfoRepo) GetUserByUserID(userID uint64) (*UserInfo, error) {
	var user UserInfo

	// 查询用户信息
	result := DB.Where("id = ? ", userID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found with ID %d", userID)
		}
		return nil, result.Error
	}

	return &user, nil
}
