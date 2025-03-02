package request

// Message 表示需要签名的消息结构
type Message struct {
	Domain         string `json:"domain"`
	Statement      string `json:"statement"`
	URI            string `json:"uri"`
	Version        string `json:"version"`
	ChainID        int    `json:"chain_id"`
	Nonce          string `json:"nonce"`
	IssuedAt       string `json:"issued_at"`
	ExpirationTime string `json:"expiration_time"`
}
