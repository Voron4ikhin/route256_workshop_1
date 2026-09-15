package order

import (
	"context"
	"fmt"
	"route256/loms/internal/domain"
	"route256/loms/internal/port/out"
)

type PayService struct {
	name            string
	orderRepository out.OrderRepository
	stockRepository out.StockRepository
}

func NewPayService(orderRepository out.OrderRepository, stockRepository out.StockRepository) *PayService {
	return &PayService{
		name:            "order.PayService",
		orderRepository: orderRepository,
		stockRepository: stockRepository,
	}
}

func (s *PayService) PayOrder(ctx context.Context, orderID int64) error {
	order, err := s.orderRepository.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", s.name, domain.ErrOrderNotFound, err)
	}
	if order.Status != domain.StatusAwaitingPayment {
		return fmt.Errorf("%s: %w", s.name, domain.ErrStatusToPay)
	}

	if err = s.stockRepository.ReserveRemove(ctx, order.Items); err != nil {
		return fmt.Errorf("%s: %w: %w", s.name, domain.ErrReserveRemoveOrder, err)
	}

	if err = s.orderRepository.SetStatus(ctx, orderID, domain.StatusPayed); err != nil {
		return fmt.Errorf("%s: %w: %w", s.name, domain.ErrStatusSetter, err)
	}

	return nil
}
