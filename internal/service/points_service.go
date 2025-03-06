package service

import (
	"fmt"
	"game-fun-be/internal/model"
	"game-fun-be/internal/response"
)

type PointsServiceImpl struct {
	userInfoRepo *model.UserInfoRepo
}

func NewPointsServiceImpl(userInfoRepo *model.UserInfoRepo) *PointsServiceImpl {
	return &PointsServiceImpl{userInfoRepo: userInfoRepo}
}

func (s *PointsServiceImpl) Points(userID uint64, chainType model.ChainType) response.Response {
	s.userInfoRepo.GetUserByUserID(userID)

	var pointsResponse response.PointsResponse
	return response.Success(pointsResponse)
}

func (s *PointsServiceImpl) PointsDetail(userID uint64, cursor *uint, limit int, chainType model.ChainType) (response.Response, error) {
	records, new_curs, has_more, err := model.GetPointRecordsByUserIDWithCursor(userID, cursor, limit)
	if err != nil {
		return response.DBErr("", err), err
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
			Points:    fmt.Sprintf("%.8f", record.PointsChange),
			Timestamp: record.CreateTime.Unix(),
			Type:      typeName, // You might want to map RecordType to appropriate string value
		}
	}

	pointsDetailsResponse := response.PointsDetailsResponse{
		Details: details,
		Cursor:  new_curs,
		HasMore: has_more,
	}
	return response.Success(pointsDetailsResponse), nil
}

func (s *PointsServiceImpl) PointsEstimated(userID string, chainType model.ChainType) response.Response {
	pointsEstimatedResponse := response.PointsEstimatedResponse{
		EstimatedPoints: "12862.90277",
	}
	return response.Success(pointsEstimatedResponse)
}
