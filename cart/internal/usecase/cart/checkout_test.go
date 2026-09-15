package cart_test

import (
	"context"
	"route256/cart/internal/domain"
	"route256/cart/internal/mocks"
	"route256/cart/internal/usecase/cart"
	"testing"
)

func TestCheckoutService_Checkout_Success(t *testing.T) {
	cartItems := []domain.CartItem{{SKU: 1, Count: 2}}
	var checkoutItems []domain.CartItem

	repo := &mocks.CartRepository{
		GetListFunc: func(ctx context.Context, user int64) ([]domain.CartItem, error) {
			return cartItems, nil
		},
	}
	loms := &mocks.LomsClient{
		CheckoutFunc: func(ctx context.Context, user int64, items []domain.CartItem) (int64, error) {
			checkoutItems = items
			return 99, nil
		},
	}

	svc := cart.NewCheckoutService(repo, loms)
	orderID, err := svc.Checkout(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if orderID != 99 {
		t.Fatalf("expected order id 99, got %d", orderID)
	}
	if len(checkoutItems) != 1 || checkoutItems[0].SKU != 1 {
		t.Fatalf("expected loms checkout to receive cart items, got %v", checkoutItems)
	}
}
