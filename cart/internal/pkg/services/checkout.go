package services

import (
	"context"
	"route256/cart/internal/pkg/handlers"
)

type CartCheckouter interface {
	GetList(ctx context.Context, user int64) ([]handlers.CartItem, error)
}

type LomsCheckouter interface {
	Checkout(ctx context.Context, user int64, items []handlers.CartItem) (int64, error)
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
	itemList, err := s.cartCheckouter.GetList(ctx, user)
	if err != nil {
		return 0, err
	}
	order, err := s.lomsCheckouter.Checkout(ctx, user, itemList)
	if err != nil {
		return 0, err
	}

	return order, nil
}
