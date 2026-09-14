package services

import (
	"context"
)

type CartClearer interface {
	Clear(ctx context.Context, user int64) error
}

type ClearService struct {
	name        string
	cartClearer CartClearer
}

func NewClearService(cartClearer CartClearer) *ClearService {
	return &ClearService{
		name:        "cart delete service",
		cartClearer: cartClearer,
	}
}

func (s ClearService) Clear(ctx context.Context, user int64) error {
	return s.cartClearer.Clear(ctx, user)
}
