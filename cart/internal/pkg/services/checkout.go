package services

import (
	"context"
	"route256/cart/internal/pkg/handlers"
)

type CartCheckouter interface {
	GetList(ctx context.Context, user int64) ([]handlers.FullCartItem, error)
}

type LomsCheckouter interface {
	Checkout(ctx context.Context, user int64, items []handlers.FullCartItem) (int64, error)
}

type CheckoutService struct {
	name           string
	cartCheckouter CartCheckouter
	lomsCheckouter LomsCheckouter
}

func NewCheckoutService(cartCheckouter CartCheckouter, lomsCheckouter LomsCheckouter) *CheckoutService {
	return &CheckoutService{
		name:           "checkout service",
		cartCheckouter: cartCheckouter,
		lomsCheckouter: lomsCheckouter,
	}
}

func (s CheckoutService) Checkout(ctx context.Context, user int64) (int64, error) {
	list, err := s.cartCheckouter.GetList(ctx, user)
	if err != nil {
		return 0, err
	}
	order, err := s.lomsCheckouter.Checkout(ctx, user, list)
	if err != nil {
		return 0, err
	}

	return order, nil
}
