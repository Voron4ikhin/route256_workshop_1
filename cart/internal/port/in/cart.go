package in

import (
	"context"
	"route256/cart/internal/domain"
)

type ItemAdder interface {
	Add(ctx context.Context, user int64, sku uint32, count uint16) error
}

type ItemDeleter interface {
	Delete(ctx context.Context, user int64, sku uint32) error
}

type CartLister interface {
	GetList(ctx context.Context, user int64) ([]domain.FullCartItem, uint32, error)
}

type CartClearer interface {
	Clear(ctx context.Context, user int64) error
}

type CartCheckouter interface {
	Checkout(ctx context.Context, user int64) (int64, error)
}
