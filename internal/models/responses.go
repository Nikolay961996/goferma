package models

type OrdersResponse struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}

type BalanceResponse struct {
	Accrual   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
