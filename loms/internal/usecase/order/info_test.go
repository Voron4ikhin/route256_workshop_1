package order_test

import (
	"context"
	"errors"
	"route256/loms/internal/domain"
	"route256/loms/internal/mocks"
	"route256/loms/internal/usecase/order"
	"testing"
)

func TestInfoService_GetOrderInfo_Success(t *testing.T) {
	want := &domain.Order{ID: 42, UserID: 7, Status: domain.StatusNew}
	repo := &mocks.OrderRepository{
		GetByIDFunc: func(ctx context.Context, orderID int64) (*domain.Order, error) {
			return want, nil
		},
	}

	svc := order.NewInfoService(repo)
	got, err := svc.GetOrderInfo(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestInfoService_GetOrderInfo_NotFound(t *testing.T) {
	repo := &mocks.OrderRepository{
		GetByIDFunc: func(ctx context.Context, orderID int64) (*domain.Order, error) {
			return nil, errors.New("not found")
		},
	}

	svc := order.NewInfoService(repo)
	_, err := svc.GetOrderInfo(context.Background(), 42)
	if !errors.Is(err, domain.ErrOrderNotFound) {
		t.Fatalf("expected ErrOrderNotFound, got %v", err)
	}
}
