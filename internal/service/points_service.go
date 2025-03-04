package service

import (
	"my-token-ai-be/internal/model"
	"my-token-ai-be/internal/response"
)

type PointsServiceImpl struct{}

func NewPointsServiceImpl() *PointsServiceImpl {
	return &PointsServiceImpl{}
}

func (s *PointsServiceImpl) Points(userID string, chainType model.ChainType) response.Response {
	var pointsResponse response.PointsResponse
	return response.Success(pointsResponse)
}

func (s *PointsServiceImpl) PointsDetail(userID, page, limit string, chainType model.ChainType) response.Response {
	var pointsDetailsResponse response.PointsDetailsResponse
	return response.Success(pointsDetailsResponse)
}

func (s *PointsServiceImpl) PointsEstimated(userID string, chainType model.ChainType) response.Response {
	pointsEstimatedResponse := response.PointsEstimatedResponse{
		EstimatedPoints: "12862.90277",
	}
	return response.Success(pointsEstimatedResponse)
}
