package service

import (
	"my-token-ai-be/internal/response"
)

type PointsServiceImpl struct{}

func NewPointsServiceImpl() *PointsServiceImpl {
	return &PointsServiceImpl{}
}

func (s *PointsServiceImpl) Points(userID string) response.Response {
	var pointsResponse response.PointsResponse
	return response.Success(pointsResponse)
}

func (s *PointsServiceImpl) PointsDetail(userID, page, limit string) response.Response {
	var pointsDetailsResponse response.PointsDetailsResponse
	return response.Success(pointsDetailsResponse)
}

func (s *PointsServiceImpl) PointsEstimated(userID string) response.Response {
	pointsEstimatedResponse := response.PointsEstimatedResponse{
		EstimatedPoints: "12862.90277",
	}
	return response.Success(pointsEstimatedResponse)
}
