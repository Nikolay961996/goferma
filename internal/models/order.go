package models

type Order struct {
	Id            int64
	Number        string
	CurrentStatus OrderStatus
}
