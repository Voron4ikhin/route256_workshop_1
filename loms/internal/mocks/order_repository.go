package mocks

import (
	"context"

	"route256/loms/internal/domain"
)

type OrderRepository struct {
	CreateFunc       func(ctx context.Context, user int64, items []domain.OrderItem) (int64, error)
	SetStatusFunc    func(ctx context.Context, orderID int64, status domain.OrderStatus) error
	GetByIDFunc      func(ctx context.Context, orderID int64) (*domain.Order, error)
	ListByStatusFunc func(ctx context.Context, status domain.OrderStatus) ([]*domain.Order, error)
}

func (m *OrderRepository) Create(ctx context.Context, user int64, items []domain.OrderItem) (int64, error) {
	return m.CreateFunc(ctx, user, items)
}

func (m *OrderRepository) SetStatus(ctx context.Context, orderID int64, status domain.OrderStatus) error {
	return m.SetStatusFunc(ctx, orderID, status)
}

func (m *OrderRepository) GetByID(ctx context.Context, orderID int64) (*domain.Order, error) {
	return m.GetByIDFunc(ctx, orderID)
}

func (m *OrderRepository) ListByStatus(ctx context.Context, status domain.OrderStatus) ([]*domain.Order, error) {
	return m.ListByStatusFunc(ctx, status)
}
