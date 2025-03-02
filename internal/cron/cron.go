package cron

import (
	"context"
	"fmt"
	"my-token-ai-be/internal/model"
	"my-token-ai-be/internal/pkg/util"
	"my-token-ai-be/internal/service"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	cronJob *cron.Cron
)

// InitCronJobs 初始化并启动定时任务
func InitCronJobs() {
	cronJob = cron.New(cron.WithSeconds())
	cronJob.Start()
	addJobs()

	entries := cronJob.Entries()

	// 打印头部
	util.Log().Info("\n====================================")
	util.Log().Info("        Scheduled Tasks (%d)", len(entries))
	util.Log().Info("====================================\n")

	// 遍历打印每个任务
	for i, entry := range entries {
		fullName := runtime.FuncForPC(reflect.ValueOf(entry.Job).Pointer()).Name()
		funcName := fullName[strings.LastIndex(fullName, ".")+1:]
		nextRun := entry.Next.In(time.Local).Format("2006-01-02 15:04:05")

		// 使用字符串拼接方式，使源码更整齐
		message := fmt.Sprintf(
			"Task %d:\n"+
				"    Name: %s\n"+
				"    Next Run: %s\n"+
				"-------------------\n",
			i+1,
			funcName,
			nextRun,
		)
		util.Log().Info("%s", message)
	}

	// 打印启动确认
	util.Log().Info("All tasks scheduled successfully!")
	util.Log().Info("====================================\n")
}

// addJobs 添加所有定时任务
func addJobs() {

	// 每分钟执行一次
	_, err := cronJob.AddFunc("0 */1 * * * *", completedTokenDataRefreshTaskQuery)
	if err != nil {
		util.Log().Error("Failed to add completedToken refresh task: %v", err)
	}

	// 每分钟执行一次
	_, err = cronJob.AddFunc("0 */1 * * * *", swapToken1mDataRefreshTaskQuery)
	if err != nil {
		util.Log().Error("Failed to add 1m swap token refresh task: %v", err)
	}

	// 每5分钟执行一次
	_, err = cronJob.AddFunc("0 */5 * * * *", swapToken5mDataRefreshTaskQuery)
	if err != nil {
		util.Log().Error("Failed to add 5m swap token refresh task: %v", err)
	}

	// 每5分钟执行一次
	_, err = cronJob.AddFunc("0 */5 * * * *", swapToken1hDataRefreshTaskQuery)
	if err != nil {
		util.Log().Error("Failed to add 1h swap token refresh task: %v", err)
	}

	// 添加每5分钟刷新热门币种的任务
	_, err = cronJob.AddFunc("0 */5 * * * *", swapToken6hDataRefreshTaskQuery)
	if err != nil {
		util.Log().Error("Failed to add refreshHotTokensJob: %v", err)
	}

	// 添加每5分钟刷新热门币种的任务
	_, err = cronJob.AddFunc("0 */6 * * * *", swapToken1dDataRefreshTaskQuery)
	if err != nil {
		util.Log().Error("Failed to add refreshHotTokensJob: %v", err)
	}

	// 添加每5分钟刷新热门币种的任务
	_, err = cronJob.AddFunc("0 */5 * * * *", trading24hJob)
	if err != nil {
		util.Log().Error("Failed to add refreshHotTokensJob: %v", err)
	}

	// 添加每5分钟刷新热门币种的任务
	_, err = cronJob.AddFunc("0 */5 * * * *", trading6hJob)
	if err != nil {
		util.Log().Error("Failed to add refreshHotTokensJob: %v", err)
	}

	// 每小时执行一次的任务
	_, err = cronJob.AddFunc("0 0 * * * *", hourlyTask)
	if err != nil {
		util.Log().Error("Failed to add hourly task: %v", err)
	}

	// 每天凌晨执行的任务
	_, err = cronJob.AddFunc("0 0 0 * * *", dailyTask)
	if err != nil {
		util.Log().Error("Failed to add daily task: %v", err)
	}

	// 每5分执行一次的钟定时任务
	_, err = cronJob.AddFunc("0 */5 * * * *", every5MinuteTask)
	if err != nil {
		util.Log().Error("Failed to add searchDocumentsJob: %v", err)
	}

	// 添加每分钟获取 SOL 价格的任务
	_, err = cronJob.AddFunc("0 * * * * *", fetchSolPrice)
	if err != nil {
		util.Log().Error("Failed to add SOL price fetching task: %v", err)
	}

	// 添加每天早上7点执行重新索引的任务
	_, err = cronJob.AddFunc("0 0 7 * * *", func() {
		if err := ExecuteReindexJob(); err != nil {
			util.Log().Error("Reindex job failed: %v", err)
		}
	})
	if err != nil {
		util.Log().Error("Failed to add reindex task: %v", err)
	}

}

// 每1分钟刷新热门币种的任务
func completedTokenDataRefreshTaskQuery() {
	CompletedTokenDataRefreshTaskQuery()
}

// 每1分钟刷新热门币种的任务
func swapToken1mDataRefreshTaskQuery() {
	SwapToken1mDataRefreshTaskQuery()
}

// 每5分钟刷新热门币种的任务
func swapToken5mDataRefreshTaskQuery() {
	SwapToken5mDataRefreshTaskQuery()
}
func swapToken1hDataRefreshTaskQuery() {
	SwapToken1hDataRefreshTaskQuery()
}
func swapToken6hDataRefreshTaskQuery() {
	SwapToken6hDataRefreshTaskQuery()
}
func swapToken1dDataRefreshTaskQuery() {
	SwapToken1dDataRefreshTaskQuery()
}

// DeleteDocumentsJob 删除文档的任务
// func deleteDocumentsJob() {
// 	DeleteDocumentsJob()
// }

// hourlyTask 每小时执行的任务
func hourlyTask() {
	const task = "[hourlyTask]"
	util.Log().Info("%s Starting task\n", task)
	// 在这里添加每小时需要执行的逻辑
	createNextDayTable()
	util.Log().Info("%s Task completed\n", task)
}

// dailyTask 每天执的任务
func dailyTask() {
	// 在这里添加每天需要执行的逻辑
	// RefreshHotTokensJob()
}

// every5MinuteTask 每5分钟执行的任务
func every5MinuteTask() {
	// 在这里添加每5分钟需要执行的逻辑

}

func createNextDayTable() {
	const task = "[createNextDayTable]"
	util.Log().Info("%s Starting at: %v\n", task, time.Now())
	util.Log().Info("%s Current timezone: %v\n", task, time.Now().Location())

	tomorrow := time.Now().Add(24 * time.Hour)
	util.Log().Info("%s Will create table for date: %v\n", task, tomorrow.Format("2006-01-02"))

	err := model.CreateTableForDate(tomorrow.Format("20060102"))
	if err != nil {
		util.Log().Error("%s Failed to create table: %v\n", task, err)
		return
	}

	util.Log().Info("%s Table creation completed successfully\n", task)
	fmt.Println("Table creation completed successfully")
}

// fetchSolPrice 每分钟获取 SOL 价格的任务
func fetchSolPrice() {
	const task = "[fetchSolPrice]"
	ctx := context.Background()
	err := service.FetchAndStoreSolPrice(ctx)
	if err != nil {
		util.Log().Error("%s Failed to fetch and store SOL price: %v\n", task, err)
	}
}

// StopCronJobs 停止所有定时任务
func StopCronJobs() {
	if cronJob != nil {
		cronJob.Stop()
		util.Log().Info("Cron jobs stopped\n")
	}
}

func trading6hJob() {
	Trading6hJob()
}

func trading24hJob() {
	Trading24hJob()
}
