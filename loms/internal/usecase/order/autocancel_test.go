package order_test

import (
	"context"
	"errors"
	"route256/loms/internal/domain"
	"route256/loms/internal/mocks"
	"route256/loms/internal/usecase/order"
	"testing"
	"time"
)

func TestAutoCancelService_SweepUnpaidOrders_CancelsOnlyOldEnough(t *testing.T) {
	now := time.Now()
	orders := []*domain.Order{
		{ID: 1, Status: domain.StatusAwaitingPayment, CreatedAt: now.Add(-15 * time.Minute)}, // old enough
		{ID: 2, Status: domain.StatusAwaitingPayment, CreatedAt: now.Add(-1 * time.Minute)},  // too fresh
	}

	var cancelledIDs []int64
	repo := &mocks.OrderRepository{
		ListByStatusFunc: func(ctx context.Context, status domain.OrderStatus) ([]*domain.Order, error) {
			if status != domain.StatusAwaitingPayment {
				t.Fatalf("unexpected status queried: %s", status)
			}
			return orders, nil
		},
	}
	canceler := &mocks.OrderCanceler{
		CancelOrderFunc: func(ctx context.Context, orderID int64) error {
			cancelledIDs = append(cancelledIDs, orderID)
			return nil
		},
	}

	svc := order.NewAutoCancelService(repo, canceler, 10*time.Minute)
	cancelled, err := svc.SweepUnpaidOrders(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cancelled != 1 {
		t.Fatalf("expected 1 cancelled order, got %d", cancelled)
	}
	if len(cancelledIDs) != 1 || cancelledIDs[0] != 1 {
		t.Fatalf("expected only order 1 to be cancelled, got %v", cancelledIDs)
	}
}

func TestAutoCancelService_SweepUnpaidOrders_IgnoresRaceWithPayment(t *testing.T) {
	orders := []*domain.Order{
		{ID: 1, Status: domain.StatusAwaitingPayment, CreatedAt: time.Now().Add(-20 * time.Minute)},
	}

	repo := &mocks.OrderRepository{
		ListByStatusFunc: func(ctx context.Context, status domain.OrderStatus) ([]*domain.Order, error) {
			return orders, nil
		},
	}
	canceler := &mocks.OrderCanceler{
		CancelOrderFunc: func(ctx context.Context, orderID int64) error {
			return domain.ErrStatusToPay
		},
	}

	svc := order.NewAutoCancelService(repo, canceler, 10*time.Minute)
	cancelled, err := svc.SweepUnpaidOrders(context.Background())
	if err != nil {
		t.Fatalf("expected the status race to be swallowed, got error: %v", err)
	}
	if cancelled != 0 {
		t.Fatalf("expected 0 cancelled orders, got %d", cancelled)
	}
}

func TestAutoCancelService_SweepUnpaidOrders_ReportsRealFailures(t *testing.T) {
	orders := []*domain.Order{
		{ID: 1, Status: domain.StatusAwaitingPayment, CreatedAt: time.Now().Add(-20 * time.Minute)},
	}
	wantErr := errors.New("storage unavailable")

	repo := &mocks.OrderRepository{
		ListByStatusFunc: func(ctx context.Context, status domain.OrderStatus) ([]*domain.Order, error) {
			return orders, nil
		},
	}
	canceler := &mocks.OrderCanceler{
		CancelOrderFunc: func(ctx context.Context, orderID int64) error {
			return wantErr
		},
	}

	svc := order.NewAutoCancelService(repo, canceler, 10*time.Minute)
	cancelled, err := svc.SweepUnpaidOrders(context.Background())
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected wrapped %v, got %v", wantErr, err)
	}
	if cancelled != 0 {
		t.Fatalf("expected 0 cancelled orders, got %d", cancelled)
	}
}
