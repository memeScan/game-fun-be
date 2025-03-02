package httpRespone

type PythResponse struct {
	Status string    `json:"s"`
	Time   []int64   `json:"t"`
	Close  []float64 `json:"c"`
}