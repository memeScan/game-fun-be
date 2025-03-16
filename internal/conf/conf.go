package conf

import (
	"os"
)

// 环境变量配置
const (
	// 环境类型
	ENV_PROD  = "release"
	ENV_TEST  = "test"
	ENV_DEBUG = "debug"
)

var (
	// 缓存当前环境
	currentEnv string

	// ES_INDEX_TOKEN_TRANSACTIONS_ALIAS 代表代币交易的别名索引
	// 使用别名可以在不影响应用程序的情况下切换实际索引
	ES_INDEX_TOKEN_TRANSACTIONS_ALIAS string

	// ES_INDEX_TOKEN_TRANSACTIONS 代表代币交易的实际索引名称
	// 通常用于存储具体的代币交易数据
	ES_INDEX_TOKEN_TRANSACTIONS string

	// ES_INDEX_TOKEN_INFO 代表代币信息的索引名称
	// 用于存储代币的基本信息数据
	ES_INDEX_TOKEN_INFO string
)

// 常量定义
const (
	UNIQUE_TOKENS = "unique_tokens"
)

// 获取当前环境
func GetEnv() string {
	// 直接返回已设置的环境值
	if currentEnv != "" {
		return currentEnv
	}

	// 如果没有设置（异常情况），使用默认值
	return ENV_DEBUG
}

// 是否是生产环境
func IsProd() bool {
	return GetEnv() == ENV_PROD
}

// 是否是测试环境
func IsTest() bool {
	return GetEnv() == ENV_TEST
}

// 是否是开发环境
func IsDebug() bool {
	return GetEnv() == ENV_DEBUG
}

// SetEnv 设置环境变量
func SetEnv(env string) {
	currentEnv = env
}

// InitGameConfig 初始化相关配置
func InitGameConfig() {

	ES_INDEX_TOKEN_TRANSACTIONS_ALIAS = os.Getenv("ES_INDEX_TOKEN_TRANSACTIONS_ALIAS")
	ES_INDEX_TOKEN_TRANSACTIONS = os.Getenv("ES_INDEX_TOKEN_TRANSACTIONS")
	ES_INDEX_TOKEN_INFO = os.Getenv("ES_INDEX_TOKEN_INFO")

}
