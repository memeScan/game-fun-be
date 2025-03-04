// @title game-fun AI API
// @version 1.0
// @description This is the API server for My Token AI.
// @host localhost:3000
// @BasePath /api/v1

package main

import (
	_ "game-fun-be/docs"
	"game-fun-be/internal/initializer"
	"game-fun-be/internal/server"
	"log"
	"os"
	"strings"

	// "net/http"

	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Environment 环境类型
type Environment string

const (
	// Debug 开发环境
	Debug Environment = "debug"
	// Test 测试环境
	Test Environment = "test"
	// Release 生产环境
	Release Environment = "release"
)

// String 实现String接口，用于输出环境文件名
func (e Environment) String() string {
	switch e {
	case Debug, Test, Release:
		return string(e)
	default:
		return string(Debug) // 默认返回开发环境配置文件
	}
}

// ParseEnvironment 从字符串解析环境类型
func ParseEnvironment(env string) Environment {
	switch strings.ToLower(env) {
	case "debug":
		return Debug
	case "test":
		return Test
	case "release":
		return Release
	default:
		return Debug
	}
}

// IsValid 检查环境是否有效
func (e Environment) IsValid() bool {
	switch e {
	case Debug, Test, Release:
		return true
	default:
		return false
	}
}

// @title My API
// @version 1.0
// @description This is a sample API for demonstration purposes.
// @host 8.218.102.112:4998
// @BasePath /api/v1
func main() {
	currentEnv := ParseEnvironment(os.Getenv("APP_ENV"))

	// 载入配置文件
	if err := godotenv.Load(".env." + currentEnv.String()); err != nil {
		log.Printf("Error loading env file: %v", err)
	}

	// 初始化配置
	initializer.Setup(currentEnv.String())

	// 装载路由
	gin.SetMode(os.Getenv("GIN_MODE"))

	r := server.NewRouter()

	// 从环境变量获取端口，如果没有设置，则使用默认值 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "4880" // 默认端口
	}

	// 启动 pprof
	// go func() {
	// 	logger.Info("Starting pprof server on :8081")
	// 	if err := http.ListenAndServe("0.0.0.0:8081", nil); err != nil {
	// 		logger.Error("pprof server failed: %v", err)
	// 	}
	// }()

	r.Run(":" + port)
}
