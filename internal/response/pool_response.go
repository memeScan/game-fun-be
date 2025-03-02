package response

type PoolResponse struct {
	Data     []PoolData     `json:"data"`
	Included []IncludedData `json:"included"`
}

type PoolData struct {
	ID           string       `json:"id"`
	Type         string       `json:"type"`
	Attributes   Attributes   `json:"attributes"`
	Relationships Relationships `json:"relationships"`
}

type Attributes struct {
	BaseTokenPriceUSD                string            `json:"base_token_price_usd"`
	BaseTokenPriceNativeCurrency     string            `json:"base_token_price_native_currency"`
	QuoteTokenPriceUSD               string            `json:"quote_token_price_usd"`
	QuoteTokenPriceNativeCurrency    string            `json:"quote_token_price_native_currency"`
	BaseTokenPriceQuoteToken         string            `json:"base_token_price_quote_token"`
	QuoteTokenPriceBaseToken         string            `json:"quote_token_price_base_token"`
	Address                          string            `json:"address"`
	Name                             string            `json:"name"`
	PoolCreatedAt                    string            `json:"pool_created_at"`
	TokenPriceUSD                    string            `json:"token_price_usd"`
	FDVUSD                           string            `json:"fdv_usd"`
	MarketCapUSD                     interface{}       `json:"market_cap_usd"` // Use interface{} to handle null values
	PriceChangePercentage            map[string]string `json:"price_change_percentage"`
	Transactions                     map[string]TransactionData `json:"transactions"`
	VolumeUSD                        map[string]string `json:"volume_usd"`
	ReserveInUSD                     string            `json:"reserve_in_usd"`
	LpBurnPercentage                 float64           `json:"lp_burn_percentage"`
}

type TransactionData struct {
	Buys    int `json:"buys"`
	Sells   int `json:"sells"`
	Buyers  int `json:"buyers"`
	Sellers int `json:"sellers"`
}

type Relationships struct {
	BaseToken  TokenRelationship `json:"base_token"`
	QuoteToken TokenRelationship `json:"quote_token"`
	Dex        DexRelationship   `json:"dex"`
}

type TokenRelationship struct {
	Data TokenData `json:"data"`
}

type DexRelationship struct {
	Data DexData `json:"data"`
}

type TokenData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type DexData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type IncludedData struct {
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Attributes TokenAttributes `json:"attributes"`
}

type TokenAttributes struct {
	Address          string      `json:"address"`
	Name             string      `json:"name"`
	Symbol           string      `json:"symbol"`
	Decimals         int         `json:"decimals"`
	ImageURL         string      `json:"image_url"`
	CoingeckoCoinID  interface{} `json:"coingecko_coin_id"` // Use interface{} to handle null values
} 