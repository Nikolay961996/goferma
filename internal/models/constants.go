package models

import "time"

const (
	JWTExpireTime = 2 * time.Hour
)

type OrderStatus int

const (
	New OrderStatus = iota
	Processing
	Invalid
	Processed
)

func (s OrderStatus) String() string {
	return [...]string{"NEW", "PROCESSING", "INVALID", "PROCESSED"}[s]
}

type contextKey string

const (
	UserIDKey contextKey = "userID"
)
