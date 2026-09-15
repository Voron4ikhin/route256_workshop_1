package item_test

import (
	"context"
	"errors"
	"route256/cart/internal/domain"
	"route256/cart/internal/mocks"
	"route256/cart/internal/usecase/item"
	"testing"
)

func TestAddService_Add_Success(t *testing.T) {
	var addedUser int64
	var addedSKU uint32
	var addedCount uint16

	svc := item.NewAddService(
		&mocks.ProductClient{GetProductInfoFunc: func(ctx context.Context, sku uint32) (string, uint32, error) {
			return "widget", 100, nil
		}},
		&mocks.StockClient{GetStocksFunc: func(ctx context.Context, sku uint32) (uint64, error) {
			return 10, nil
		}},
		&mocks.CartRepository{AddFunc: func(ctx context.Context, user int64, sku uint32, count uint16) error {
			addedUser, addedSKU, addedCount = user, sku, count
			return nil
		}},
	)

	if err := svc.Add(context.Background(), 1, 42, 3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addedUser != 1 || addedSKU != 42 || addedCount != 3 {
		t.Fatalf("unexpected Add call: user=%d sku=%d count=%d", addedUser, addedSKU, addedCount)
	}
}

func TestAddService_Add_InsufficientStock(t *testing.T) {
	svc := item.NewAddService(
		&mocks.ProductClient{GetProductInfoFunc: func(ctx context.Context, sku uint32) (string, uint32, error) {
			return "widget", 100, nil
		}},
		&mocks.StockClient{GetStocksFunc: func(ctx context.Context, sku uint32) (uint64, error) {
			return 1, nil
		}},
		&mocks.CartRepository{},
	)

	err := svc.Add(context.Background(), 1, 42, 3)
	if !errors.Is(err, domain.ErrInsufficientStocks) {
		t.Fatalf("expected ErrInsufficientStocks, got %v", err)
	}
}

func TestAddService_Add_UnknownProduct(t *testing.T) {
	wantErr := errors.New("product not found")
	svc := item.NewAddService(
		&mocks.ProductClient{GetProductInfoFunc: func(ctx context.Context, sku uint32) (string, uint32, error) {
			return "", 0, wantErr
		}},
		&mocks.StockClient{},
		&mocks.CartRepository{},
	)

	err := svc.Add(context.Background(), 1, 42, 3)
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected wrapped %v, got %v", wantErr, err)
	}
}
