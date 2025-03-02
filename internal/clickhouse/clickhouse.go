package clickhouse

import (
	"context"
	"my-token-ai-be/internal/pkg/util"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

// ClickHouseClient ClickHouse客户端单例
var ClickHouseClient driver.Conn

// ClickHouse 初始化ClickHouse连接
func ClickHouse() {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{os.Getenv("CLICKHOUSE_ADDR")},
		Auth: clickhouse.Auth{
			Database: os.Getenv("CLICKHOUSE_DB"),
			Username: os.Getenv("CLICKHOUSE_USER"),
			Password: os.Getenv("CLICKHOUSE_PASSWORD"),
		},

		// 连接池优化 - 适度调整
		MaxOpenConns:    util.GetEnvAsInt("CLICKHOUSE_MAX_OPEN_CONNS", 16),                       // 考虑到8核CPU，设置为2倍CPU核心数
		MaxIdleConns:    util.GetEnvAsInt("CLICKHOUSE_MAX_IDLE_CONNS", 8),                        // 设置为MaxOpenConns的一半
		ConnMaxLifetime: util.GetEnvAsDuration("CLICKHOUSE_CONN_MAX_LIFETIME", 1800*time.Second), // 降低到30分钟，避免连接过久不释放

		// 超时设置
		DialTimeout: util.GetEnvAsDuration("CLICKHOUSE_DIAL_TIMEOUT", 10*time.Second),

		// 压缩设置
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionZSTD,
		},

		// 批处理和缓冲设置 - 降低批量大小，增加处理频率
		Settings: clickhouse.Settings{
			"max_block_size":                   util.GetEnvAsInt("CLICKHOUSE_MAX_BLOCK_SIZE", 50000), // 降低单个块大小
			"max_insert_block_size":            util.GetEnvAsInt("CLICKHOUSE_MAX_INSERT_BLOCK_SIZE", 50000),
			"min_insert_block_size_rows":       5000,                                                  // 降低最小批量，更频繁写入
			"min_insert_block_size_bytes":      5000000,                                               // 降低到5MB
			"max_execution_time":               util.GetEnvAsInt("CLICKHOUSE_MAX_EXECUTION_TIME", 60), // 降低执行超时时间
			"max_memory_usage":                 "8000000000",                                          // 限制单个查询内存使用为8GB
			"max_memory_usage_for_all_queries": "24000000000",                                         // 限制总内存使用为24GB
		},

		// 关闭调试
		Debug: false,

		// 预分配缓冲区
		BlockBufferSize: 15, // 降低缓冲区大小
	})

	if err != nil {
		util.Log().Panic("连接ClickHouse不成功: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 测试连接
	if err := conn.Ping(ctx); err != nil {
		util.Log().Panic("ClickHouse连接测试失败: %v", err)
	}

	// 获取版本号
	var version string
	if err := conn.QueryRow(ctx, "SELECT version()").Scan(&version); err == nil {
		util.Log().Info("ClickHouse版本: %s", version)
	}

	ClickHouseClient = conn
}

// 添加关闭函数
func CloseClickHouse() error {
	if ClickHouseClient != nil {
		return ClickHouseClient.Close()
	}
	return nil
}
