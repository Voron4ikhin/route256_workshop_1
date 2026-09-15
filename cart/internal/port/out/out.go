package out

import (
	"context"
	"route256/cart/internal/domain"
)

type CartRepository interface {
	Add(ctx context.Context, user int64, sku uint32, count uint16) error
	Delete(ctx context.Context, user int64, sku uint32) error
	GetList(ctx context.Context, user int64) ([]domain.CartItem, error)
	Clear(ctx context.Context, user int64) error
}

type ProductClient interface {
	GetProductInfo(ctx context.Context, sku uint32) (string, uint32, error)
}

type StockClient interface {
	GetStocks(ctx context.Context, sku uint32) (uint64, error)
}

type LomsClient interface {
	Checkout(ctx context.Context, user int64, items []domain.CartItem) (int64, error)
}
