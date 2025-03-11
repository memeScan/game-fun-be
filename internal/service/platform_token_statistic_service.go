package service

import (
	"fmt"
	"game-fun-be/internal/model"
	"game-fun-be/internal/response"
)

type PlatformTokenStatisticServiceImpl struct {
	platformTokenStatisticRepo *model.PlatformTokenStatisticRepo
}

func NewPlatformTokenStatisticServiceImpl(platformTokenStatisticRepo *model.PlatformTokenStatisticRepo) *PlatformTokenStatisticServiceImpl {
	return &PlatformTokenStatisticServiceImpl{platformTokenStatisticRepo: platformTokenStatisticRepo}
}

func (s *PlatformTokenStatisticServiceImpl) IncrementStatistics(tokenAddress string, amounts map[model.StatisticType]uint64) error {
	if tokenAddress == "" {
		return fmt.Errorf("token address cannot be empty")
	}

	return s.platformTokenStatisticRepo.IncrementStatisticsAndUpdateTime(
		tokenAddress,
		amounts,
	)
}

func (s *PlatformTokenStatisticServiceImpl) GetTokenPointsStatistic(tokenAddress string, chainType uint8) response.Response {
	if tokenAddress == "" {
		return response.ParamErr("token address cannot be empty", fmt.Errorf("token address cannot be empty"))
	}

	statistics, err := s.platformTokenStatisticRepo.GetTokenPointsStatistic(tokenAddress, chainType)
	if err != nil {
		return response.DBErr("获取token统计数据失败", err)
	}
	platfromTokenStatisticResponse := response.PlatfromTokenStatisticResponse{
		TokenAddress:  statistics.TokenAddress,
		FeeAmount:     formatSol(statistics.FeeAmount),
		BackAmount:    formatPoints(statistics.BackAmount),
		BackSolAmount: formatSol(statistics.BackSolAmount),
		BurnAmount:    formatPoints(statistics.BurnAmount),
		PointsAmount:  formatPoints(statistics.PointsAmount),
	}

	return response.Success(platfromTokenStatisticResponse)
}
