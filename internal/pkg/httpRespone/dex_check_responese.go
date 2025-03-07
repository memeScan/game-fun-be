package httpRespone

// Website 结构体
type Website struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

// Social 结构体
type Social struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

// Boost 结构体
type Boost struct {
	Active int `json:"active"` // 假设 active 是整数
}

// Data 结构体
type DexCheckData struct {
	Address  string    `json:"address"`
	Websites []Website `json:"websites"`
	Socials  []Social  `json:"socials"`
	Boosts   *Boost    `json:"boosts"` // 使用指针以处理可能为 null 的情况
}

// DexCheckResponse 结构体
type DexCheckResponse struct {
	Code    int            `json:"code"`
	Data    []DexCheckData `json:"data"`
	Message string         `json:"message"`
}
