package mocks

import (
	"context"

	"route256/loms/internal/domain"
)

type StockRepository struct {
	GetBySKUFunc      func(ctx context.Context, sku uint32) uint64
	ReserveFunc       func(ctx context.Context, items []domain.OrderItem) error
	ReserveRemoveFunc func(ctx context.Context, items []domain.OrderItem) error
	ReserveCancelFunc func(ctx context.Context, items []domain.OrderItem) error
}

func (m *StockRepository) GetBySKU(ctx context.Context, sku uint32) uint64 {
	return m.GetBySKUFunc(ctx, sku)
}

func (m *StockRepository) Reserve(ctx context.Context, items []domain.OrderItem) error {
	return m.ReserveFunc(ctx, items)
}

func (m *StockRepository) ReserveRemove(ctx context.Context, items []domain.OrderItem) error {
	return m.ReserveRemoveFunc(ctx, items)
}

func (m *StockRepository) ReserveCancel(ctx context.Context, items []domain.OrderItem) error {
	return m.ReserveCancelFunc(ctx, items)
}
