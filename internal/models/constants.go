package models

import "time"

const (
	JWT_EXPIRE_TIME = 2 * time.Hour
)

type OrderStatus int

const (
	NEW OrderStatus = iota
	PROCESSING
	INVALID
	PROCESSED
)
