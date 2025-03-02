package response

// TokensSearchResponse 代币搜索响应结构体
// 包含代币的基本信息、价格信息、交易数据、安全指标等
type TokensSearchResponse struct {
    // 基本信息
    Chain    string `json:"chain"`    // 区块链网络类型(如: ETH, BSC, SOL等)
    Symbol   string `json:"symbol"`   // 代币符号(如: BTC, ETH等)
    Name     string `json:"name"`     // 代币名称
    Decimals int    `json:"decimals"` // 代币精度(小数位数)
    Logo     string `json:"logo"`     // 代币logo图片URL
    Address  string `json:"address"`  // 代币合约地址

    // 价格相关
    Price    float64 `json:"price"`     // 当前价格
    Price1h  float64 `json:"price_1h"`  // 1小时前价格
    Price24h float64 `json:"price_24h"` // 24小时前价格

    // 交易数据
    Swaps5m   int     `json:"swaps_5m"`    // 5分钟内交易次数
    Swaps1h   int     `json:"swaps_1h"`    // 1小时内交易次数
    Swaps6h   int     `json:"swaps_6h"`    // 6小时内交易次数
    Swaps24h  int     `json:"swaps_24h"`   // 24小时内交易次数
    Volume24h float64 `json:"volume_24h"`   // 24小时交易量
    Liquidity float64 `json:"liquidity"`    // 流动性(池子大小)

    // 供应信息
    TotalSupply int `json:"total_supply"` // 代币总供应量

    // 代币特征
    SymbolLen     int  `json:"symbol_len"`      // 代币符号长度
    NameLen       int  `json:"name_len"`        // 代币名称长度
    IsInTokenList bool `json:"is_in_token_list"` // 是否在官方代币列表中
    HotLevel      int  `json:"hot_level"`       // 热度等级
    IsShowAlert   bool `json:"is_show_alert"`    // 是否显示风险警告

    // 安全指标
    BuyTax     float64 `json:"buy_tax"`     // 买入税率
    SellTax    float64 `json:"sell_tax"`    // 卖出税率
    IsHoneypot bool    `json:"is_honeypot"` // 是否为蜜罐合约(无法卖出的骗局代币)
    Renounced  bool    `json:"renounced"`   // 是否已放弃合约所有权

    // 持仓分布
    Top10HolderRate float64 `json:"top_10_holder_rate"` // 前10大持有者占比

    // 合约权限
    RenouncedMint         int `json:"renounced_mint"`           // 是否放弃铸币权(0:未放弃 1:已放弃)
    RenouncedFreezeAccount int `json:"renounced_freeze_account"` // 是否放弃冻结账户权限(0:未放弃 1:已放弃)

    // 燃烧相关
    BurnRatio   string `json:"burn_ratio"`   // 代币燃烧比例
    BurnStatus  string `json:"burn_status"`  // 燃烧状态(如:正常燃烧、异常燃烧等)
}