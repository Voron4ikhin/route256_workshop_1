package order_test

import (
	"context"
	"errors"
	"route256/loms/internal/domain"
	"route256/loms/internal/mocks"
	"route256/loms/internal/usecase/order"
	"testing"
)

func TestCancelService_CancelOrder_Success(t *testing.T) {
	items := []domain.OrderItem{{SKU: 1, Count: 2}}
	var gotStatus domain.OrderStatus
	var cancelledCalled bool

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
		ReserveCancelFunc: func(ctx context.Context, items []domain.OrderItem) error {
			cancelledCalled = true
			return nil
		},
	}

	svc := order.NewCancelService(repo, stocks)
	if err := svc.CancelOrder(context.Background(), 42); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cancelledCalled {
		t.Fatal("expected ReserveCancel to be called")
	}
	if gotStatus != domain.StatusCancelled {
		t.Fatalf("expected status %q, got %q", domain.StatusCancelled, gotStatus)
	}
}

func TestCancelService_CancelOrder_WrongStatus(t *testing.T) {
	repo := &mocks.OrderRepository{
		GetByIDFunc: func(ctx context.Context, orderID int64) (*domain.Order, error) {
			return &domain.Order{ID: orderID, Status: domain.StatusPayed}, nil
		},
	}
	stocks := &mocks.StockRepository{}

	svc := order.NewCancelService(repo, stocks)
	err := svc.CancelOrder(context.Background(), 42)
	if !errors.Is(err, domain.ErrStatusToPay) {
		t.Fatalf("expected ErrStatusToPay, got %v", err)
	}
}
