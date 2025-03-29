package model

import (
	"time"
)

const (
	PUMP_INITIAL_REAL_TOKEN_RESERVES    = "793100000000000"
	PUMP_INITIAL_VIRTUAL_SOL_RESERVES   = "30000000000"
	PUMP_INITIAL_VIRTUAL_TOKEN_RESERVES = "1073000000000000"
	PUMP_INITIAL_REALSOL_TOKEN_RESERVES = "0"
)

const (
	TokenExpireDuration = 24 * time.Hour
)

const (
	PointsDecimal = 6
)

const (
	LoginMessageTemplate = "Sign in to the game.exchange\r\n\r\nTimestamp: %s"
)

const (
	G_TOKEN_ADDRESS  = "8iFREvVdmLKxVeibpC5VLRr1S6X5dm7gYR3VCU1wpump"
	G_TOKEN_DECIMALS = 6
)

// PlatformType 定义交易平台类型
type PlatformType uint8

const (
	PlatformTypePump     PlatformType = iota + 1 // 1
	PlatformTypeRaydium                          // 2
	PlatformTypePumpSwap                         // 3
)

// String 方法用于返回 PlatformType 的字符串表示
func (p PlatformType) String() string {
	switch p {
	case PlatformTypePump:
		return "Pump"
	case PlatformTypeRaydium:
		return "Raydium"
	case PlatformTypePumpSwap:
		return "PumpSwap"
	default:
		return "Unknown"
	}
}

// ChainType 定义链类型
type ChainType uint8

const (
	ChainTypeUnknown  ChainType = iota // 0
	ChainTypeSolana                    // 1
	ChainTypeEthereum                  // 2
	ChainTypeBSC                       // 3
)

// String 方法用于返回 ChainType 的字符串表示
func (c ChainType) String() string {
	switch c {
	case ChainTypeSolana:
		return "Solana"
	case ChainTypeEthereum:
		return "Ethereum"
	case ChainTypeBSC:
		return "Bsc"
	default:
		return "Unknown"
	}
}

// String 方法用于返回 ChainType 的字符串表示
func (c ChainType) Uint8() uint8 {
	return uint8(c)
}

// String 方法用于返回 ChainType 的字符串表示
func (c ChainType) ToString() string {
	switch c {
	case ChainTypeSolana:
		return "solana"
	case ChainTypeEthereum:
		return "ethereum"
	case ChainTypeBSC:
		return "bsc"
	default:
		return "Unknown"
	}
}

// FromString 根据字符串参数（如 sol/eth/bsc）返回 ChainType
func ChainTypeFromString(chainStr string) ChainType {
	switch chainStr {
	case "sol", "Solana":
		return ChainTypeSolana
	case "eth", "Ethereum":
		return ChainTypeEthereum
	case "bsc", "BSC":
		return ChainTypeBSC
	default:
		return ChainTypeUnknown
	}
}

// 新增 CreatedPlatformType 定义
type CreatedPlatformType uint8

const (
	CreatedPlatformTypeUnknown  CreatedPlatformType = iota // Unknown/其他 = 0
	CreatedPlatformTypePump                                // Pump 平台 = 1
	CreatedPlatformTypeMoonshot                            // Moonshot 平台 = 2
	CreatedPlatformTypeGamPump                             // GamPump 平台 = 3
)

const (
	SOL_DECIMALS uint8 = 9
)

// 平台对应的代币数量精度
var createdPlatformDecimals = map[CreatedPlatformType]uint8{
	CreatedPlatformTypePump:     6, // Pump平台代币数量精度
	CreatedPlatformTypeMoonshot: 6, // Moonshot平台代币数量精度
}

// String 方法用于返回 CreatedPlatformType 的字符串表示
func (p CreatedPlatformType) String() string {
	switch p {
	case CreatedPlatformTypePump:
		return "Pump"
	case CreatedPlatformTypeMoonshot:
		return "Moonshot"
	default:
		return "Unknown"
	}
}

// GetDecimals 获取平台对应的代币精度
func (p CreatedPlatformType) GetDecimals() uint8 {
	if decimals, ok := createdPlatformDecimals[p]; ok {
		return decimals
	}
	return 6 // 默认返回6
}

// 链相关常量
const (
	// WSOL 是 Solana 上 SOL 的包装代币地址
	SolanaWrappedSOLAddress = "So11111111111111111111111111111111111111112"
)

// GetNativeTokenAddress 根据链类型获取生代币的包装地址
func (c ChainType) GetNativeTokenAddress() string {
	switch c {
	case ChainTypeSolana:
		return SolanaWrappedSOLAddress
	default:
		return ""
	}
}

// DevStatus 代币开发者状态
type DevStatus uint8

const (
	DevStatusInit     DevStatus = iota // 0-初始
	DevStatusHold                      // 1-持有
	DevStatusSell                      // 2-卖出
	DevStatusClear                     // 3-清仓
	DevStatusIncrease                  // 4-加仓
	DevStatusAddLP                     // 5-加池子
	DevStatusRemoveLP                  // 6-减池子
	DevStatusBurn                      // 7-烧币
)

// TransactionType 交易类型
type TransactionType uint8

const (
	TransactionTypeBuy             TransactionType = iota + 1 // 1-买
	TransactionTypeSell                                       // 2-卖
	TransactionTypeAddLiquidity                               // 3-加池子
	TransactionTypeRemoveLiquidity                            // 4-减池子
	TransactionTypeBurn                                       // 5-烧币
)

// ProxyType 代理类型
type ProxyType uint8

const (
	ProxyTypeUnknown ProxyType = iota // 0-未知
	ProxyTypeGame                     // 1-game  内盘交易
	ProxyTypeGameIn                   // 2-gameIn 积分购买
)

// 回购地址
const TreasuryAddress = "EfeuLDmQbockFee9vYWyPijw1eZoihhHp6ufYGDb7CJT"

// GameProxyAddress game代理合约地址
const GameProxyAddress = "GameUoe5CTRtakryJjFWdzixbgqPCycwdFHhdEiwvHsE"

type PointType string

const (
	AvailablePoints    PointType = "available_points"
	TradingPoints      PointType = "trading_points"
	InvitePoints       PointType = "invite_points"
	InviteRebate       PointType = "invite_rebate"
	WithdrawableRebate PointType = "withdrawable_rebate"
)

type RecordType int8

const (
	Trading  RecordType = 1
	Invite   RecordType = 2
	Activity RecordType = 3
	BuyG     RecordType = 4
)

type StatisticType string

const (
	FeeAmount     StatisticType = "fee_amount"
	BackAmount    StatisticType = "back_amount"
	BackSolAmount StatisticType = "back_sol_amount"
	BurnAmount    StatisticType = "burn_amount"
	PointsAmount  StatisticType = "points_amount"
)
