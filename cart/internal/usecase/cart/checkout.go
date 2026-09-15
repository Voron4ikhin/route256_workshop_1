package cart

import (
	"context"
	"fmt"
	"route256/cart/internal/port/out"
)

type CheckoutService struct {
	name           string
	cartRepository out.CartRepository
	lomsClient     out.LomsClient
}

func NewCheckoutService(cartRepository out.CartRepository, lomsClient out.LomsClient) *CheckoutService {
	return &CheckoutService{
		name:           "cart.CheckoutService",
		cartRepository: cartRepository,
		lomsClient:     lomsClient,
	}
}

func (s *CheckoutService) Checkout(ctx context.Context, user int64) (int64, error) {
	items, err := s.cartRepository.GetList(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("%s: get cart list: %w", s.name, err)
	}

	orderID, err := s.lomsClient.Checkout(ctx, user, items)
	if err != nil {
		return 0, fmt.Errorf("%s: checkout in loms: %w", s.name, err)
	}

	return orderID, nil
}
