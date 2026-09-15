package mocks

import (
	"context"

	"route256/cart/internal/domain"
)

type CartRepository struct {
	AddFunc     func(ctx context.Context, user int64, sku uint32, count uint16) error
	DeleteFunc  func(ctx context.Context, user int64, sku uint32) error
	GetListFunc func(ctx context.Context, user int64) ([]domain.CartItem, error)
	ClearFunc   func(ctx context.Context, user int64) error
}

func (m *CartRepository) Add(ctx context.Context, user int64, sku uint32, count uint16) error {
	return m.AddFunc(ctx, user, sku, count)
}

func (m *CartRepository) Delete(ctx context.Context, user int64, sku uint32) error {
	return m.DeleteFunc(ctx, user, sku)
}

func (m *CartRepository) GetList(ctx context.Context, user int64) ([]domain.CartItem, error) {
	return m.GetListFunc(ctx, user)
}

func (m *CartRepository) Clear(ctx context.Context, user int64) error {
	return m.ClearFunc(ctx, user)
}
