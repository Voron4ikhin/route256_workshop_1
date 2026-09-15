package order_test

import (
	"context"
	"errors"
	"route256/loms/internal/domain"
	"route256/loms/internal/mocks"
	"route256/loms/internal/usecase/order"
	"testing"
)

func TestCreateService_CreateOrder_Success(t *testing.T) {
	items := []domain.OrderItem{{SKU: 1, Count: 2}}
	var gotStatus domain.OrderStatus

	repo := &mocks.OrderRepository{
		CreateFunc: func(ctx context.Context, user int64, items []domain.OrderItem) (int64, error) {
			return 42, nil
		},
		SetStatusFunc: func(ctx context.Context, orderID int64, status domain.OrderStatus) error {
			gotStatus = status
			return nil
		},
	}
	stocks := &mocks.StockRepository{
		ReserveFunc: func(ctx context.Context, items []domain.OrderItem) error {
			return nil
		},
	}

	svc := order.NewCreateService(repo, stocks)
	orderID, err := svc.CreateOrder(context.Background(), 1, items)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if orderID != 42 {
		t.Fatalf("expected order id 42, got %d", orderID)
	}
	if gotStatus != domain.StatusAwaitingPayment {
		t.Fatalf("expected status %q, got %q", domain.StatusAwaitingPayment, gotStatus)
	}
}

func TestCreateService_CreateOrder_ReserveFails(t *testing.T) {
	items := []domain.OrderItem{{SKU: 1, Count: 2}}
	var gotStatus domain.OrderStatus

	repo := &mocks.OrderRepository{
		CreateFunc: func(ctx context.Context, user int64, items []domain.OrderItem) (int64, error) {
			return 42, nil
		},
		SetStatusFunc: func(ctx context.Context, orderID int64, status domain.OrderStatus) error {
			gotStatus = status
			return nil
		},
	}
	stocks := &mocks.StockRepository{
		ReserveFunc: func(ctx context.Context, items []domain.OrderItem) error {
			return errors.New("not enough stock")
		},
	}

	svc := order.NewCreateService(repo, stocks)
	_, err := svc.CreateOrder(context.Background(), 1, items)
	if !errors.Is(err, domain.ErrReserveOrder) {
		t.Fatalf("expected ErrReserveOrder, got %v", err)
	}
	if gotStatus != domain.StatusFailed {
		t.Fatalf("expected status %q, got %q", domain.StatusFailed, gotStatus)
	}
}

func TestCreateService_CreateOrder_CreateFails(t *testing.T) {
	repo := &mocks.OrderRepository{
		CreateFunc: func(ctx context.Context, user int64, items []domain.OrderItem) (int64, error) {
			return 0, errors.New("storage down")
		},
	}
	stocks := &mocks.StockRepository{}

	svc := order.NewCreateService(repo, stocks)
	_, err := svc.CreateOrder(context.Background(), 1, nil)
	if !errors.Is(err, domain.ErrCreateOrder) {
		t.Fatalf("expected ErrCreateOrder, got %v", err)
	}
}
