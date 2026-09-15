package out

import (
	"context"
	"route256/loms/internal/domain"
)

type OrderRepository interface {
	Create(ctx context.Context, user int64, items []domain.OrderItem) (int64, error)
	SetStatus(ctx context.Context, orderID int64, status domain.OrderStatus) error
	GetByID(ctx context.Context, orderID int64) (*domain.Order, error)
	ListByStatus(ctx context.Context, status domain.OrderStatus) ([]*domain.Order, error)
}

type StockRepository interface {
	GetBySKU(ctx context.Context, sku uint32) uint64
	Reserve(ctx context.Context, items []domain.OrderItem) error
	ReserveRemove(ctx context.Context, items []domain.OrderItem) error
	ReserveCancel(ctx context.Context, items []domain.OrderItem) error
}
