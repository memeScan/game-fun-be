package request

// TickersRequest 获取市场行情请求参数
// @Description 获取市场行情时提交的请求参数
type TickersRequest struct {
	// 排序字段，支持以下值：
	// - MARKET_CAP：市值
	// - PRICE_CHANGE_5M：5 分钟价格变化
	// - PRICE_CHANGE_1H：1 小时价格变化
	// - PRICE_CHANGE_24H：24 小时价格变化
	// - NATIVE_VOLUME_1H：1 小时原生交易量
	// - NATIVE_VOLUME_24H：24 小时原生交易量
	// - TX_COUNT_24H：24 小时交易次数
	// - HOLDERS：持有者数量
	// - INITIALIZE_AT：初始化时间
	// - Links：链接
	SortedBy string `form:"sorted_by" example:"INITIALIZE_AT" validate:"oneof=MARKET_CAP PRICE_CHANGE_5M PRICE_CHANGE_1H PRICE_CHANGE_24H NATIVE_VOLUME_1H NATIVE_VOLUME_24H TX_COUNT_24H HOLDERS INITIALIZE_AT Links"`
	// 排序方向，支持以下值：
	// - DESC：降序
	// - ASC：升序
	SortDirection string `form:"sort_direction" example:"DESC" validate:"oneof=DESC ASC"`
	// 分页游标，用于分页查询
	PageCursor string `form:"page_cursor" example:""`
	// 每页返回的数据条数
	Limit int `form:"limit" example:"50"`
	// 搜索关键字，用于筛选数据
	Search string `form:"search" example:""`
	// 新交易对的时间分辨率，例如 1D（1 天）
	NewPairsResolution string `form:"new_pairs_resolution" example:"1D"`
}
