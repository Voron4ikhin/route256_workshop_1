package order

import (
	"context"
	"fmt"
	"route256/loms/internal/domain"
	"route256/loms/internal/port/out"
)

type CreateService struct {
	name            string
	orderRepository out.OrderRepository
	stockRepository out.StockRepository
}

func NewCreateService(orderRepository out.OrderRepository, stockRepository out.StockRepository) *CreateService {
	return &CreateService{
		name:            "order.CreateService",
		orderRepository: orderRepository,
		stockRepository: stockRepository,
	}
}

func (s *CreateService) CreateOrder(ctx context.Context, user int64, items []domain.OrderItem) (int64, error) {
	orderID, err := s.orderRepository.Create(ctx, user, items)
	if err != nil {
		return 0, fmt.Errorf("%s: %w: %w", s.name, domain.ErrCreateOrder, err)
	}

	if err = s.stockRepository.Reserve(ctx, items); err != nil {
		if statusErr := s.orderRepository.SetStatus(ctx, orderID, domain.StatusFailed); statusErr != nil {
			return orderID, fmt.Errorf("%s: %w: %w", s.name, domain.ErrStatusSetter, statusErr)
		}
		return orderID, fmt.Errorf("%s: %w: %w", s.name, domain.ErrReserveOrder, err)
	}

	if err = s.orderRepository.SetStatus(ctx, orderID, domain.StatusAwaitingPayment); err != nil {
		return orderID, fmt.Errorf("%s: %w: %w", s.name, domain.ErrStatusSetter, err)
	}

	return orderID, nil
}
