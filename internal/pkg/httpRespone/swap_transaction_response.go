package httpRespone

// SwapTransactionResponse represents the structure of the JSON response
type SwapTransactionResponse struct {
	Code int `json:"code"`
	Data struct {
		LastValidBlockHeight int `json:"lastValidBlockHeight"`
		SwapTransaction      struct {
			Signatures []Signature        `json:"signatures"`
			Message    TransactionMessage `json:"message"`
		} `json:"swapTransaction"`
		Base64SwapTransaction string `json:"base64SwapTransaction"`
	} `json:"data"`
	Message string `json:"message"`
}

// Signature represents the structure of a signature object
type Signature struct {
	Signature map[string]string `json:"signature"` // Use a map to capture numeric keys
}

// Message represents the structure of the message object
type TransactionMessage struct {
	Header struct {
		NumRequiredSignatures       int `json:"numRequiredSignatures"`
		NumReadonlySignedAccounts   int `json:"numReadonlySignedAccounts"`
		NumReadonlyUnsignedAccounts int `json:"numReadonlyUnsignedAccounts"`
	} `json:"header"`
	AccountKeys       []string               `json:"accountKeys"`
	RecentBlockhash   string                 `json:"recentBlockhash"`
	Instructions      []Instruction          `json:"instructions"`
	IndexToProgramIds map[string]interface{} `json:"indexToProgramIds"`
}

// Instruction represents the structure of an instruction object
type Instruction struct {
	ProgramIdIndex int    `json:"programIdIndex"`
	Accounts       []int  `json:"accounts"`
	Data           string `json:"data"`
}
