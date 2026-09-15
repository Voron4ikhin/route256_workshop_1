package cart

import (
	"context"
	"route256/cart/internal/port/out"
)

type ClearService struct {
	name           string
	cartRepository out.CartRepository
}

func NewClearService(cartRepository out.CartRepository) *ClearService {
	return &ClearService{
		name:           "cart.ClearService",
		cartRepository: cartRepository,
	}
}

func (s *ClearService) Clear(ctx context.Context, user int64) error {
	return s.cartRepository.Clear(ctx, user)
}
