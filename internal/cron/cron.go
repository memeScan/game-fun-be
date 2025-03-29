package cron

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"

	"game-fun-be/internal/model"
	"game-fun-be/internal/pkg/util"
	"game-fun-be/internal/service"

	"github.com/robfig/cron/v3"
)

var cronJob *cron.Cron

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
	// 每小时执行一次的任务
	_, err := cronJob.AddFunc("0 0 * * * *", hourlyTask)
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

	// 每5分钟去更新redis代币配置
	_, err = cronJob.AddFunc("0 */5 * * * *", syncTokenConfigToRedis)
	if err != nil {
		util.Log().Error("Failed to add syncTokenConfigToRedis: %v", err)
	}

	// 每5分执行一次积分任务
	// _, err = cronJob.AddFunc("0 */10 * * * *", ExecutePointJob)
	// if err != nil {
	// 	util.Log().Error("Failed to add ExecutePointJob: %v", err)
	// }

	// 添加每分钟获取 SOL 价格的任务
	_, err = cronJob.AddFunc("0 * * * * *", fetchSolPrice)
	if err != nil {
		util.Log().Error("Failed to add SOL price fetching task: %v", err)
	}

	_, err = cronJob.AddFunc("0 */10 * * * *", refreshHolderTaskJobs)
	if err != nil {
		util.Log().Error("Failed to add ExecutePointJob: %v", err)
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

// 每5分钟去刷新代币列表的holder
func refreshHolderTaskJobs() {
	RefreshHotTokensHolderJob()
}

// 每5分钟去更新redis代币配置
func syncTokenConfigToRedis() {
	SyncTokenConfigToRedis()
}
