package order_test

import (
	"context"
	"errors"
	"route256/loms/internal/domain"
	"route256/loms/internal/mocks"
	"route256/loms/internal/usecase/order"
	"testing"
)

func TestPayService_PayOrder_Success(t *testing.T) {
	items := []domain.OrderItem{{SKU: 1, Count: 2}}
	var gotStatus domain.OrderStatus
	var removedItems []domain.OrderItem

	repo := &mocks.OrderRepository{
		GetByIDFunc: func(ctx context.Context, orderID int64) (*domain.Order, error) {
			return &domain.Order{ID: orderID, Status: domain.StatusAwaitingPayment, Items: items}, nil
		},
		SetStatusFunc: func(ctx context.Context, orderID int64, status domain.OrderStatus) error {
			gotStatus = status
			return nil
		},
	}
	stocks := &mocks.StockRepository{
		ReserveRemoveFunc: func(ctx context.Context, items []domain.OrderItem) error {
			removedItems = items
			return nil
		},
	}

	svc := order.NewPayService(repo, stocks)
	if err := svc.PayOrder(context.Background(), 42); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotStatus != domain.StatusPayed {
		t.Fatalf("expected status %q, got %q", domain.StatusPayed, gotStatus)
	}
	if len(removedItems) != 1 || removedItems[0].SKU != 1 {
		t.Fatalf("expected reserve-remove to be called with order items, got %v", removedItems)
	}
}

func TestPayService_PayOrder_WrongStatus(t *testing.T) {
	repo := &mocks.OrderRepository{
		GetByIDFunc: func(ctx context.Context, orderID int64) (*domain.Order, error) {
			return &domain.Order{ID: orderID, Status: domain.StatusNew}, nil
		},
	}
	stocks := &mocks.StockRepository{}

	svc := order.NewPayService(repo, stocks)
	err := svc.PayOrder(context.Background(), 42)
	if !errors.Is(err, domain.ErrStatusToPay) {
		t.Fatalf("expected ErrStatusToPay, got %v", err)
	}
}

func TestPayService_PayOrder_OrderNotFound(t *testing.T) {
	repo := &mocks.OrderRepository{
		GetByIDFunc: func(ctx context.Context, orderID int64) (*domain.Order, error) {
			return nil, errors.New("not found")
		},
	}
	stocks := &mocks.StockRepository{}

	svc := order.NewPayService(repo, stocks)
	err := svc.PayOrder(context.Background(), 42)
	if !errors.Is(err, domain.ErrOrderNotFound) {
		t.Fatalf("expected ErrOrderNotFound, got %v", err)
	}
}
