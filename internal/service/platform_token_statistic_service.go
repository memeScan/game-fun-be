package service

import (
	"fmt"
	"game-fun-be/internal/model"
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
