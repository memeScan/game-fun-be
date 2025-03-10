package test

import (
	"os"

	"game-fun-be/internal/model"
	"game-fun-be/internal/pkg/httpUtil"
	"game-fun-be/internal/pkg/util"
	"game-fun-be/internal/redis"
	"game-fun-be/internal/server"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var s *gin.Engine

func init() {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Load configurations
	confInit()
	// Initialize router
	s = server.NewRouter()
}

// TestSetup initializes the test environment
func TestSetup() {
	// Set test environment
	os.Setenv("APP_ENV", "test")
	os.Setenv("LOG_LEVEL", "debug")

	confInit()
}

func confInit() {
	// Try to load test environment file from multiple possible locations
	envFiles := []string{
		".env.test",
		"../.env.test",
		"../../.env.test",
	}

	envLoaded := false
	for _, envFile := range envFiles {
		if err := godotenv.Load(envFile); err == nil {
			envLoaded = true
			break
		}
	}

	if !envLoaded {
		// Set default values if env file not found
		defaultEnv := map[string]string{
			"MYSQL_DSN":  "user:password@tcp(localhost:3306)/dbname?charset=utf8&parseTime=True&loc=Local",
			"REDIS_ADDR": "localhost:6379",
			"LOG_LEVEL":  "debug",
		}

		for key, value := range defaultEnv {
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}

		util.Log().Warning("No .env.test file found, using default values")
	}

	// Set up logger
	// util.BuildLogger(os.Getenv("LOG_LEVEL"))
	// conf.InitTradingConfig()

	// 设置日志级别

	// Initialize Redis and Database connections
	redis.Redis()
	model.Database(os.Getenv("MYSQL_DSN"))
	util.Log().Info("Test environment setup completed1")

	endpoint := os.Getenv("API_ENDPOINT")
	httpUtil.InitAPI(&endpoint)
	httpUtil.InitMetrics(redis.RedisClient)
}
