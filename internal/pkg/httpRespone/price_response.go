package httpRespone

// PriceResponse 定义 Birdeye API 返回的价格信息
type PriceResponse struct {
	Data    PriceData `json:"data"`    // 价格数据
	Success bool      `json:"success"` // 请求是否成功
	Message string    `json:"message"` // 返回消息
}

// PriceData 定义价格数据的具体内容
type PriceData struct {
	Value           float64 `json:"value"`           // 代币价格
	UpdateUnixTime  int64   `json:"updateUnixTime"`  // 更新时间戳（Unix 时间）
	UpdateHumanTime string  `json:"updateHumanTime"` // 更新时间（人类可读格式）
}
