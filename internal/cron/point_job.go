package cron

import (
	"os"
	"time"

	"game-fun-be/internal/clickhouse"
	"game-fun-be/internal/model"
	"game-fun-be/internal/pkg/util"
	"game-fun-be/internal/service"
)

func ExecutePointJob() {
	util.Log().Info("执行积分任务")

	tokenAddress := os.Getenv("TOKEN_ADDRESS")
	// vaultAddress := os.Getenv("VAULT_ADDRESS")

	pointRecordsRepo := model.NewPointRecordsRepo()
	userInfoRepo := model.NewUserInfoRepo()
	PlatformTokenStatisticRepo := model.NewPlatformTokenStatisticRepo()
	pointsService := service.NewPointsServiceImpl(userInfoRepo, pointRecordsRepo, PlatformTokenStatisticRepo)

	// globalService := service.NewGlobalServiceImpl()

	// 获取当前时间并向下取整到最近的10分钟
	now := time.Now()
	endTime := now.Truncate(10 * time.Minute)
	// 获取上一个10分钟的时间窗口
	startTime := endTime.Add(-10 * time.Minute)

	transactions, err := clickhouse.QueryProxyTransactionsByTime(startTime, endTime, 1, tokenAddress)
	if err != nil {
		util.Log().Error("Failed to add ExecutePointJob: %v", err)
		return
	}

	// resp := globalService.TickerBalance(vaultAddress, tokenAddress, model.ChainTypeSolana)

	// if resp.Code != 0 {
	// 	util.Log().Error("Failed to get vault balance: %v", resp.Error)
	// 	return
	// }

	vaultAmount := uint64(70206195194)

	// 统计所有用户的交易量总和
	var quotaSum uint64
	userAddressMap := make(map[string][]clickhouse.ProxyTransaction)
	for _, transaction := range transactions {
		quotaSum += transaction.QuoteTokenAmount
		userAddressMap[transaction.UserAddress] = append(userAddressMap[transaction.UserAddress], transaction)
	}

	for userAddress, transactions := range userAddressMap {
		var transactionAmountDetails []service.TransactionAmountDetail
		for _, transaction := range transactions {
			transactionAmountDetails = append(transactionAmountDetails, service.TransactionAmountDetail{
				QuotaAmount:     transaction.QuoteTokenAmount,
				TransactionHash: transaction.TransactionHash,
				TransactionTime: transaction.TransactionTime,
			})
		}

		pointsService.SavePointsEveryTimeBucket(service.TransactionAmountDetailByTime{
			UserAddress:              userAddress,
			QuotaTotalAmount:         quotaSum,
			VaultAmount:              vaultAmount,
			TransactionAmountDetails: transactionAmountDetails,
			StartTime:                startTime,
			EndTime:                  endTime,
		})
	}
}
