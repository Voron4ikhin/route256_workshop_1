package mocks

import (
	"context"

	"route256/cart/internal/domain"
)

type ProductClient struct {
	GetProductInfoFunc func(ctx context.Context, sku uint32) (string, uint32, error)
}

func (m *ProductClient) GetProductInfo(ctx context.Context, sku uint32) (string, uint32, error) {
	return m.GetProductInfoFunc(ctx, sku)
}

type StockClient struct {
	GetStocksFunc func(ctx context.Context, sku uint32) (uint64, error)
}

func (m *StockClient) GetStocks(ctx context.Context, sku uint32) (uint64, error) {
	return m.GetStocksFunc(ctx, sku)
}

type LomsClient struct {
	CheckoutFunc func(ctx context.Context, user int64, items []domain.CartItem) (int64, error)
}

func (m *LomsClient) Checkout(ctx context.Context, user int64, items []domain.CartItem) (int64, error) {
	return m.CheckoutFunc(ctx, user, items)
}
