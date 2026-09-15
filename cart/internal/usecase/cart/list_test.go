package cart_test

import (
	"context"
	"route256/cart/internal/domain"
	"route256/cart/internal/mocks"
	"route256/cart/internal/usecase/cart"
	"testing"
)

func TestListService_GetList_ComputesTotalPrice(t *testing.T) {
	repo := &mocks.CartRepository{
		GetListFunc: func(ctx context.Context, user int64) ([]domain.CartItem, error) {
			return []domain.CartItem{{SKU: 1, Count: 2}, {SKU: 2, Count: 1}}, nil
		},
	}
	product := &mocks.ProductClient{
		GetProductInfoFunc: func(ctx context.Context, sku uint32) (string, uint32, error) {
			if sku == 1 {
				return "widget", 100, nil
			}
			return "gadget", 50, nil
		},
	}

	svc := cart.NewListService(product, repo)
	items, total, err := svc.GetList(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if total != 250 {
		t.Fatalf("expected total price 250, got %d", total)
	}
}
