package in

import (
	"context"
	"route256/loms/internal/domain"
)

type OrderCreator interface {
	CreateOrder(ctx context.Context, user int64, items []domain.OrderItem) (int64, error)
}

type OrderPayer interface {
	PayOrder(ctx context.Context, orderID int64) error
}

type OrderCanceler interface {
	CancelOrder(ctx context.Context, orderID int64) error
}

type OrderInformant interface {
	GetOrderInfo(ctx context.Context, orderID int64) (*domain.Order, error)
}

type OrderAutoCanceler interface {
	SweepUnpaidOrders(ctx context.Context) (cancelled int, err error)
}
