package item

import (
	"context"
	"time"
)

type CartDeleter interface {
	Delete(ctx context.Context, user int64, sku uint32) error
}

type DeleteService struct {
	name        string
	cartDeleter CartDeleter
}

func NewDeleteService(cartDeleter CartDeleter) *DeleteService {
	return &DeleteService{
		name:        "item delete service",
		cartDeleter: cartDeleter,
	}
}

func (s DeleteService) Delete(ctx context.Context, user int64, sku uint32) error {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	return s.cartDeleter.Delete(ctx, user, sku)
}
