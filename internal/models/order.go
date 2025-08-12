package models

type Order struct {
	ID            int64
	Number        string
	CurrentStatus OrderStatus
}
