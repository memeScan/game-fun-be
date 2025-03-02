package response

import (
	"my-token-ai-be/internal/clickhouse"

	"github.com/shopspring/decimal"
)

// KlineData 序列化结构
type KlineData struct {
	Timestamp int64           `json:"timestamp"`
	Open      decimal.Decimal `json:"open"`
	High      decimal.Decimal `json:"high"`
	Low       decimal.Decimal `json:"low"`
	Close     decimal.Decimal `json:"close"`
	Volume    decimal.Decimal `json:"volume"`
}

// BuildKlineData 构建单个K线数据
func BuildKlineData(k clickhouse.Kline, decimals uint8) KlineData {
	volume := decimal.NewFromInt(int64(k.Volume))
	realVolume := volume.Shift(-int32(decimals))

	return KlineData{
		Timestamp: k.IntervalTimestamp.UnixMilli(),
		Open:      k.OpenPrice,
		High:      k.HighPrice,
		Low:       k.LowPrice,
		Close:     k.ClosePrice,
		Volume:    realVolume,
	}
}

// BuildKlineDataList 构建K线数据列表
func BuildKlineDataList(klines []clickhouse.Kline, decimals uint8) []KlineData {
	result := make([]KlineData, len(klines))
	for i, k := range klines {
		result[i] = BuildKlineData(k, decimals)
	}
	return result
}

// KlineDataResponse K线数据响应
type KlineDataResponse struct {
	Data  []KlineData `json:"data"`
	Total int64       `json:"total"`
}

// BuildKlineDataResponse 构建K线数据响应
func BuildKlineDataResponse(klineData []KlineData) Response {
	return Response{
		Code: CodeSuccess,
		Msg:  "Success",
		Data: klineData,
	}
}
