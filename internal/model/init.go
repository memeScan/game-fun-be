package model

import (
	"time"

	"my-token-ai-be/internal/pkg/util"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 数据库链接单例
var DB *gorm.DB

// Database 在中间件中初始化mysql链接
func Database(connString string) {
	// 初始化GORM日志配置
	newLogger := logger.New(
		&customWriter{}, // 使用自定义 writer
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Error, // Log level(这里记得根据需求改一下)
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,        // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(connString), &gorm.Config{
		Logger:                 newLogger,
		PrepareStmt:            true, // 可以改为 false，全局禁用预编译
		SkipDefaultTransaction: true,
		CreateBatchSize:        300,
		AllowGlobalUpdate:      false, // 添加：禁止全局更新
		QueryFields:            true,  // 添加：显式指定查询字段
		DisableAutomaticPing:   false, // 启用自动 ping，及时发现连接问题
		ConnPool: &gorm.PreparedStmtDB{
			Stmts: make(map[string]*gorm.Stmt, 200),
		},
	})
	// Error
	if connString == "" || err != nil {
		util.Log().Error("mysql lost: %v", err)
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		util.Log().Error("mysql lost: %v", err)
		panic(err)
	}

	//设置连接池
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Minute * 15)
	sqlDB.SetConnMaxIdleTime(time.Minute * 5)

	DB = db

	// 打印数据库版本
	var version string
	DB.Raw("SELECT VERSION()").Scan(&version)
	util.Log().Info("Database connected successfully. MySQL version: %v", version)

}

// customWriter 实现 logger.Writer 接口
type customWriter struct{}

func (w *customWriter) Printf(format string, args ...interface{}) {
	util.Log().Info(format, args...)
}
