package workers

import "github.com/Nikolay961996/goferma/internal/models"

type job struct {
	Order models.Order
}

type loyaltyStatus string

const (
	registered loyaltyStatus = "REGISTERED"
	invalid    loyaltyStatus = "INVALID"
	processing loyaltyStatus = "PROCESSING"
	processed  loyaltyStatus = "PROCESSED"
)

type loyaltyResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}
