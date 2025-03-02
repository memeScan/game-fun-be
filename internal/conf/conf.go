package conf

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
