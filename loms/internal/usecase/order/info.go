package order

import (
	"context"
	"fmt"
	"route256/loms/internal/domain"
	"route256/loms/internal/port/out"
)

type InfoService struct {
	name            string
	orderRepository out.OrderRepository
}

func NewInfoService(orderRepository out.OrderRepository) *InfoService {
	return &InfoService{
		name:            "order.InfoService",
		orderRepository: orderRepository,
	}
}

func (s *InfoService) GetOrderInfo(ctx context.Context, orderID int64) (*domain.Order, error) {
	order, err := s.orderRepository.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w: %w", s.name, domain.ErrOrderNotFound, err)
	}

	return order, nil
}
