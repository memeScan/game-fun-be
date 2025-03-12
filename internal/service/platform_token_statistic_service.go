package service

import (
	"fmt"
	"net/http"

	"game-fun-be/internal/model"
	"game-fun-be/internal/response"

	"github.com/shopspring/decimal"
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

	solUsdPrice, err := getSolPrice()
	if err != nil {
		return response.Err(http.StatusInternalServerError, "failed to get sol price", err)
	}

	platfromTokenStatisticResponse := response.PlatfromTokenStatisticResponse{
		TokenAddress:  statistics.TokenAddress,
		FeeAmount:     formatSolUsd(decimal.NewFromInt(int64(statistics.FeeAmount)).Mul(solUsdPrice)),
		BackAmount:    formatPoints(statistics.BackAmount),
		BackSolAmount: formatSolUsd(decimal.NewFromInt(int64(statistics.BackSolAmount)).Mul(solUsdPrice)),
		BurnAmount:    formatPoints(statistics.BurnAmount),
		PointsAmount:  formatPoints(statistics.PointsAmount),
	}

	return response.Success(platfromTokenStatisticResponse)
}
