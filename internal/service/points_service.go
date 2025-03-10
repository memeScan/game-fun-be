package service

import (
	"fmt"
	"game-fun-be/internal/model"
	"game-fun-be/internal/response"
	"time"

	"gorm.io/gorm"
)

type PointsServiceImpl struct {
	db                 *gorm.DB
	userInfoRepo       *model.UserInfoRepo
	pointRecordsRecord *model.PointRecordsRepo
}

func NewPointsServiceImpl(userInfoRepo *model.UserInfoRepo, pointRecordsRecord *model.PointRecordsRepo) *PointsServiceImpl {
	return &PointsServiceImpl{userInfoRepo: userInfoRepo, pointRecordsRecord: pointRecordsRecord}
}

/**
* get points by user id
* @param userID uint64`json:"user_id"`
* @param chainType model.ChainType`json:"chain_type"`
 */
func (s *PointsServiceImpl) Points(userID uint, chainType model.ChainType) response.Response {
	userInfo, err := s.userInfoRepo.GetUserByUserID(userID)

	if err != nil {
		return response.DBErr("", err)
	}

	var pointsResponse response.PointsResponse = response.PointsResponse{
		TradingPoints:   formatPoints(userInfo.TradingPoints),
		InvitePoints:    formatPoints(userInfo.InvitePoints),
		AvailablePoints: formatPoints(userInfo.AvailablePoints),
	}
	return response.Success(pointsResponse)
}

/**
* get points detail by user id
 */
func (s *PointsServiceImpl) PointsDetail(userID uint64, cursor *uint, limit int, chainType model.ChainType) response.Response {
	records, new_curs, has_more, err := s.pointRecordsRecord.GetPointRecordsByUserIDWithCursor(userID, cursor, limit)
	if err != nil {
		return response.DBErr("", err)
	}

	details := make([]response.PointsDetail, len(records))
	for i, record := range records {
		typeName := ""
		switch record.RecordType {
		case 1:
			typeName = "trading"
		case 2:
			typeName = "invite"
		case 3:
			typeName = "activity"
		case 4:
			typeName = "buy_g"
		}

		details[i] = response.PointsDetail{
			Points:    formatPoints(record.PointsChange),
			Timestamp: record.CreateTime.Unix(),
			Type:      typeName, // You might want to map RecordType to appropriate string value
		}
	}

	pointsDetailsResponse := response.PointsDetailsResponse{
		Details: details,
		Cursor:  new_curs,
		HasMore: has_more,
	}
	return response.Success(pointsDetailsResponse)
}

func (s *PointsServiceImpl) InvitedPointsDetail(userID uint64, cursor *uint, limit int, chainType model.ChainType) response.Response {

	// 获取用户信息
	users, new_curs, has_more, err := s.userInfoRepo.GetUsersByInviterId(userID, cursor, limit)
	if err != nil {
		return response.DBErr("获取用户信息失败", err)
	}

	fmt.Println(users)

	userIDs := make([]uint, len(users))

	userInfoMap := make(map[uint]*model.UserInfo)
	for _, user := range users {
		userInfoMap[user.ID] = user
		userIDs = append(userIDs, user.ID)
	}

	// 获取邀请积分记录
	records, err := s.pointRecordsRecord.InvitedPointsDetail(userIDs)
	if err != nil {
		return response.DBErr("获取邀请积分记录失败", err)
	}

	fmt.Println(records)

	// 构建响应数据
	details := make([]response.InvitedPointsDetail, len(records))

	for i, record := range records {
		details[i] = response.InvitedPointsDetail{
			Invitee:       userInfoMap[record.UserID].Address,
			InviteTime:    userInfoMap[record.UserID].CreateTime.Unix(),
			TradingPoints: formatPoints(record.Points),
		}
	}

	pointsTotalResponse := response.InvitedPointsTotalResponse{

		Details: details,
		HasMore: has_more,
		Cursor:  new_curs,
	}

	return response.Success(pointsTotalResponse)
}

func (s *PointsServiceImpl) PointsSave(address string, point uint64, hash string) response.Response {

	transaction_err := s.db.Transaction(func(tx *gorm.DB) error {

		user, err := s.userInfoRepo.WithTx(tx).GetUserByAddress(address, 1)
		if user == nil || err != nil { // 用户不存在
			return err // 400 Bad Request
		}

		// 创建积分记录
		insertErr := s.pointRecordsRecord.WithTx(tx).CreatePointRecord(&model.PointRecords{
			UserID:          user.ID,
			PointsChange:    point, // 积分变动
			PointsBalance:   user.AvailablePoints + point,
			RecordType:      1, // 积分类型
			TransactionHash: hash,
		})
		if insertErr != nil {
			return insertErr
		}

		userPoints := make(map[model.PointType]uint64)
		userPoints[model.AvailablePoints] = point
		userPoints[model.TradingPoints] = point

		// 更新用户积分
		if err := s.userInfoRepo.WithTx(tx).IncrementMultiplePointsAndUpdateTime(address, userPoints); err != nil {
			return err
		}

		if user.InviterID != 0 {

			inviter, err := s.userInfoRepo.WithTx(tx).GetUserByUserID(user.InviterID)
			if inviter == nil || err != nil { // 用户不存在
				return err // 400 Bad Request
			}

			invitePoints := uint64(float64(point) * 0.15)

			// 创建积分记录
			insertErr := s.pointRecordsRecord.WithTx(tx).CreatePointRecord(&model.PointRecords{
				UserID:        user.ID,
				PointsChange:  invitePoints, // 积分变动
				PointsBalance: user.AvailablePoints + invitePoints,
				RecordType:    2, // 积分类型
				InviteeID:     user.ID,
				UpdateTime:    time.Now(),
			})
			if insertErr != nil {
				return insertErr
			}

			userPoints := make(map[model.PointType]uint64)
			userPoints[model.AvailablePoints] = invitePoints
			userPoints[model.TradingPoints] = invitePoints

			// 更新用户积分
			if err := s.userInfoRepo.WithTx(tx).IncrementMultiplePointsAndUpdateTime(inviter.Address, userPoints); err != nil {
				return err
			}

		}

		if user.ParentInviteId != 0 {

			parentInviter, err := s.userInfoRepo.WithTx(tx).GetUserByUserID(user.ParentInviteId)
			if parentInviter == nil || err != nil { // 用户不存在
				return err // 400 Bad Request
			}

			parentInviterPoints := uint64(float64(point) * 0.05)

			// 创建积分记录
			insertErr := s.pointRecordsRecord.WithTx(tx).CreatePointRecord(&model.PointRecords{
				UserID:        user.ID,
				PointsChange:  parentInviterPoints, // 积分变动
				PointsBalance: user.AvailablePoints + parentInviterPoints,
				RecordType:    2, // 积分类型
				InviteeID:     user.InviterID,
				UpdateTime:    time.Now(),
			})
			if insertErr != nil {
				return insertErr
			}

			userPoints := make(map[model.PointType]uint64)
			userPoints[model.AvailablePoints] = parentInviterPoints
			userPoints[model.InvitePoints] = parentInviterPoints

			// 更新用户积分
			if err := s.userInfoRepo.WithTx(tx).IncrementMultiplePointsAndUpdateTime(parentInviter.Address, userPoints); err != nil {
				return err
			}

		}

		return nil
	})

	if transaction_err != nil {
		return response.DBErr("积分保存失败", transaction_err)

	}
	return response.Success("积分保存成功")
}

func (s *PointsServiceImpl) PointsEstimated(userID string, chainType model.ChainType) response.Response {
	pointsEstimatedResponse := response.PointsEstimatedResponse{
		EstimatedPoints: "12862.90277",
	}
	return response.Success(pointsEstimatedResponse)
}

func formatPoints(points uint64) string {
	return fmt.Sprintf("%.6f", float64(points)/1e6)
}
