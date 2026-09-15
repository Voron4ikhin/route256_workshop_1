package domain

import "time"

type OrderItem struct {
	SKU   uint32
	Count uint16
}

type OrderStatus string

const (
	StatusNew             OrderStatus = "new"
	StatusAwaitingPayment OrderStatus = "awaiting_payment"
	StatusFailed          OrderStatus = "failed"
	StatusCancelled       OrderStatus = "cancelled"
	StatusPayed           OrderStatus = "payed"
)

func (s OrderStatus) IsValid() bool {
	switch s {
	case StatusNew, StatusAwaitingPayment, StatusFailed, StatusCancelled, StatusPayed:
		return true
	default:
		return false
	}
}

type Order struct {
	ID        int64
	UserID    int64
	Items     []OrderItem
	Status    OrderStatus
	CreatedAt time.Time
}
