package es

const (
	UNIQUE_TOKENS = "unique_tokens"
	// ES_INDEX_TOKEN_TRANSACTIONS_ALIAS 代表代币交易的别名索引
	// 使用别名可以在不影响应用程序的情况下切换实际索引
	ES_INDEX_TOKEN_TRANSACTIONS_ALIAS = "token_transactions_alias"
	// ES_INDEX_TOKEN_TRANSACTIONS 代表代币交易的实际索引名称
	// 通常用于存储具体的代币交易数据
	ES_INDEX_TOKEN_TRANSACTIONS = "token_transactions_new_v1"
	// ES_INDEX_TOKEN_INFO 代表代币信息的索引名称
	// 用于存储代币的基本信息数据
	ES_INDEX_TOKEN_INFO = "token_info_v1"
)
