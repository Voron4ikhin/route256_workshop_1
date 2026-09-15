package item

import (
	"context"
	"route256/cart/internal/port/out"
)

type DeleteService struct {
	name           string
	cartRepository out.CartRepository
}

func NewDeleteService(cartRepository out.CartRepository) *DeleteService {
	return &DeleteService{
		name:           "item.DeleteService",
		cartRepository: cartRepository,
	}
}

func (s *DeleteService) Delete(ctx context.Context, user int64, sku uint32) error {
	return s.cartRepository.Delete(ctx, user, sku)
}
