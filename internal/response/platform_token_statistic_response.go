package response

// PointsResponse 定义积分数据结构
type PlatfromTokenStatisticResponse struct {
	TokenAddress  string `json:"token_address"`   // 代币地址
	FeeAmount     string `json:"fee_amount"`      // 手续费收入总数 sol
	BackAmount    string `json:"back_amount"`     // token回购数量
	BackSolAmount string `json:"back_sol_amount"` // 回购花费的 sol数量
	BurnAmount    string `json:"burn_amount"`     // token 销毁数量
	PointsAmount  string `json:"points_amount"`   // 已兑换的积分数量
}
