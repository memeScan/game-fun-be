package model

// TradeMessage represents the structure of a trade message in Kafka
type TokenTradeMessage struct {
	Mint                 string  `json:"mint"`
	SolAmount            string  `json:"solAmount"`
	TokenAmount          string  `json:"tokenAmount"`
	IsBuy                bool    `json:"isBuy"`
	User                 string  `json:"user"`
	Timestamp            int64   `json:"timestamp"`
	VirtualSolReserves   string  `json:"virtualSolReserves"`
	VirtualTokenReserves string  `json:"virtualTokenReserves"`
	RealSolReserves      string  `json:"realSolReserves"`
	RealTokenReserves    string  `json:"realTokenReserves"`
	Progress             float64 `json:"progress"`
	Signature            string  `json:"signature"`
	Block                uint64  `json:"block"`
	BondingCurve         string  `json:"bondingCurve"`
}

// TokenInfoMessage Kafka消息结构体
type TokenInfoMessage struct {
	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	URI          string `json:"uri"`
	Mint         string `json:"mint"`
	Creator      string `json:"creator"`
	BondingCurve string `json:"bondingCurve"`
	Signature    string `json:"signature"`
	Timestamp    int64  `json:"timestamp"`
	Block        uint64 `json:"block"`
}

type RaydiumSwapMessage struct {
	Timestamp         int64  `json:"timestamp"`         // 时间戳
	Block             uint64 `json:"block"`             // 区块高度
	Signature         string `json:"signature"`         // 签名
	MarketAddress     string `json:"marketAddress"`     // 市场地址
	PoolAddress       string `json:"poolAddress"`       // 交易对池子地址，由 marketId 生成
	User              string `json:"user"`              // 买卖用户地址
	IsBuy             bool   `json:"isBuy"`             // 是否买入
	QuoteToken        string `json:"quoteToken"`        // 询价代币，为 Meme 代币
	BaseToken         string `json:"baseToken"`         // 基础代币，为 SOL
	QuoteAmount       string `json:"quoteAmount"`       // 询价代币数量，不是池子中的数量
	BaseAmount        string `json:"baseAmount"`        // 基础代币数量，不是池子中的数量
	PoolQuoteReserve  string `json:"poolQuoteReserve"`  // 池子中询价代币的当前总量
	PoolBaseReserve   string `json:"poolBaseReserve"`   // 池子中基础代币的当前总量
	Decimals          int    `json:"decimals"`          // 代币精度
	ParentInstAddress string `json:"parentInstAddress"` // 父指令地址，本次新增字段
}

// RaydiumCreateMessage Raydium 创建池子的消息结构体
type RaydiumCreateMessage struct {
	Timestamp             int64  `json:"timestamp"`             // 时间戳
	Block                 uint64 `json:"block"`                 // 区块高度
	Signature             string `json:"signature"`             // 签名
	User                  string `json:"user"`                  // 用户地址
	MarketAddress         string `json:"marketAddress"`         // 市场地址
	PoolAddress           string `json:"poolAddress"`           // 交易对池子地址，由 marketId 生成
	PoolState             int    `json:"poolState"`             // 0: 初始化池子, 1: 往池子添加流动性, 2: 往池子移除流动性
	QuoteToken            string `json:"quoteToken"`            // 询价代币, 为 Meme 代币
	BaseToken             string `json:"baseToken"`             // 基础代币, 为 SOL
	PoolQuoteReserve      string `json:"poolQuoteReserve"`      // 池子中询价代币的当前总量
	PoolBaseReserve       string `json:"poolBaseReserve"`       // 池子中基础代币的当前总量
	ChangePoolQuoteAmount string `json:"changePoolQuoteAmount"` // 询价代币数量的变化值
	ChangePoolBaseAmount  string `json:"changePoolBaseAmount"`  // 基础代币数量的变化值
	Decimals              int    `json:"decimals"`              // 交易代币精度
}

// PumpFuncCompleteMessage represents a pump function completion message
type PumpFuncCompleteMessage struct {
	User         string `json:"user"`
	Mint         string `json:"mint"`
	BondingCurve string `json:"bondingCurve"`
	Timestamp    int64  `json:"timestamp"`
	Signature    string `json:"signature"`
}

// GameOutTradeMessage 代理合约外盘买卖事件消息结构体
type GameOutTradeMessage struct {
	Timestamp            int64  `json:"timestamp"`            // 时间戳
	Block                uint64 `json:"block"`                // 区块高度
	Signature            string `json:"signature"`            // 签名
	User                 string `json:"user"`                 // 用户地址
	PoolAddress          string `json:"poolAddress"`          // 池子地址
	IsBuy                bool   `json:"isBuy"`                // 是否买入
	QuoteToken           string `json:"quoteToken"`           // 询价代币，为 Meme 代币
	BaseToken            string `json:"baseToken"`            // 基础代币，为 SOL
	MarketAddress        string `json:"marketAddress"`        // 市场地址
	PoolQuoteReserve     string `json:"poolQuoteReserve"`     // 池子中询价代币的当前总量
	PoolBaseReserve      string `json:"poolBaseReserve"`      // 池子中基础代币的当前总量
	QuoteAmount          string `json:"quoteAmount"`          // 询价代币数量
	BaseAmount           string `json:"baseAmount"`           // 基础代币数量
	Decimals             int    `json:"decimals"`             // 代币精度
	FeeQuoteAmount       string `json:"feeQuoteAmount"`       // 手续费回购的代币数量
	FeeBaseAmount        string `json:"feeBaseAmount"`        // 平台收到的手续费
	BuybackFeeBaseAmount string `json:"buybackFeeBaseAmount"` // 手续费回购花的sol数量
	IsBurn               bool   `json:"isBurn"`               // 是否销毁
}

// GameInTradeMessage 代理合约内盘买事件（积分兑换买）消息结构体
type GameInTradeMessage struct {
	Timestamp     int64  `json:"timestamp"`     // 时间戳
	Block         uint64 `json:"block"`         // 区块高度
	Signature     string `json:"signature"`     // 签名
	User          string `json:"user"`          // 用户地址
	IsBuy         bool   `json:"isBuy"`         // 是否买入
	QuoteToken    string `json:"quoteToken"`    // 询价代币，为 Meme 代币
	BaseToken     string `json:"baseToken"`     // 基础代币，为 SOL
	QuoteAmount   string `json:"quoteAmount"`   // 兑换的代币数量
	BaseAmount    string `json:"baseAmount"`    // 兑换花费的sol数量
	Decimals      int    `json:"decimals"`      // 代币精度
	PointsAmount  string `json:"pointsAmount"`  // 积分数量(带单位,6位精度)
	FeeBaseAmount string `json:"feeBaseAmount"` // 平台收到的手续费
}

// PointTxStatusMessage 积分交易链上状态消息结构体
type PointTxStatusMessage struct {
	Signature string `json:"signature"` // 交易签名
	UserId    uint   `json:"userId"`    // 用户ID
	Points    uint64 `json:"points"`    // 积分数量
	Rebate    uint64 `json:"rebate"`    // 提现金额
	TxType    uint   `json:"txType"`    // 交易类型 1:积分兑换 2:手续费提取
}

// PumpAmmCreateMessage Pump AMM 创建池子的消息结构体
type PumpAmmCreateMessage struct {
	Timestamp             int64  `json:"timestamp"`             // 时间戳
	Block                 uint64 `json:"block"`                 // 区块高度
	Signature             string `json:"signature"`             // 签名
	User                  string `json:"user"`                  // 用户地址
	MarketAddress         string `json:"marketAddress"`         // 市场地址
	PoolAddress           string `json:"poolAddress"`           // 交易对池子地址，由 marketId 生成
	PoolState             int    `json:"poolState"`             // 0: 初始化池子, 1: 往池子添加流动性, 2: 往池子移除流动性
	QuoteToken            string `json:"quoteToken"`            // 询价代币, 为 Meme 代币
	BaseToken             string `json:"baseToken"`             // 基础代币, 为 SOL
	PoolQuoteReserve      string `json:"poolQuoteReserve"`      // 池子中询价代币的当前总量
	PoolBaseReserve       string `json:"poolBaseReserve"`       // 池子中基础代币的当前总量
	ChangePoolQuoteAmount string `json:"changePoolQuoteAmount"` // 询价代币数量的变化值
	ChangePoolBaseAmount  string `json:"changePoolBaseAmount"`  // 基础代币数量的变化值
	Decimals              int    `json:"decimals"`              // 交易代币精度
}

// PumpAmmSwapMessage Pump AMM 交易池子的消息结构体
type PumpAmmSwapMessage struct {
	Type              string `json:"type"`              // 消息类型，例如："PumpAmmSwap"
	Signature         string `json:"signature"`         // 交易签名
	Timestamp         int64  `json:"timestamp"`         // 时间戳
	Block             uint64 `json:"block"`             // 区块高度
	PoolAddress       string `json:"poolAddress"`       // 池子地址
	User              string `json:"user"`              // 用户地址
	IsBuy             bool   `json:"isBuy"`             // 是否为买入操作
	QuoteToken        string `json:"quoteToken"`        // 询价代币地址(Meme代币)
	BaseToken         string `json:"baseToken"`         // 基础代币地址(通常为SOL)
	QuoteAmount       string `json:"quoteAmount"`       // 询价代币数量
	BaseAmount        string `json:"baseAmount"`        // 基础代币数量
	PoolQuoteReserve  string `json:"poolQuoteReserve"`  // 池子中询价代币的当前总量
	PoolBaseReserve   string `json:"poolBaseReserve"`   // 池子中基础代币的当前总量
	Decimals          int    `json:"decimals"`          // 交易代币精度
	IsStandardOrder   bool   `json:"isStandardOrder"`   // 是否为标准订单
	ParentInstAddress string `json:"parentInstAddress"` // 父指令地址
}
