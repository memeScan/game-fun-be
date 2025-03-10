package service

import (
	"fmt"
	"game-fun-be/internal/model"
	"game-fun-be/internal/response"
)

type PointsServiceImpl struct {
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
func (s *PointsServiceImpl) Points(userID uint64, chainType model.ChainType) response.Response {
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

func (s *PointsServiceImpl) PointsSave(address string, point uint64) response.Response {

	// user, err := s.userInfoRepo.GetUserByAddress(address, 1)
	// if _ != nil {
	// 	return response.DBErr("查询用户失败", err)
	// }
	// point_before := user.AvailablePoints
	// inviter_id := user.InviterID
	// parent_inviter_id := user.ParentInviteId

	// 更新用户积分
	// err = s.userInfoRepo.UpdatePoints(address, point)

	// s.userInfoRepo.UpdatePoints(address, point)

	pointsEstimatedResponse := response.PointsEstimatedResponse{
		EstimatedPoints: "12862.90277",
	}
	return response.Success(pointsEstimatedResponse)
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
