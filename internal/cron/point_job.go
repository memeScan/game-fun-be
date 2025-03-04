package cron

import (
	"game-fun-be/internal/pkg/util"
	"game-fun-be/internal/service"
)

func ExecutePointJob() {
	util.Log().Info("执行积分任务")

	pointService := service.PointServiceImpl{}
	pointService.CalculateVolumeStatistics()
}
