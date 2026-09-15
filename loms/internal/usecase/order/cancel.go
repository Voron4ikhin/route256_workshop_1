package order

import (
	"context"
	"fmt"
	"route256/loms/internal/domain"
	"route256/loms/internal/port/out"
)

type CancelService struct {
	name            string
	orderRepository out.OrderRepository
	stockRepository out.StockRepository
}

func NewCancelService(orderRepository out.OrderRepository, stockRepository out.StockRepository) *CancelService {
	return &CancelService{
		name:            "order.CancelService",
		orderRepository: orderRepository,
		stockRepository: stockRepository,
	}
}

func (s *CancelService) CancelOrder(ctx context.Context, orderID int64) error {
	order, err := s.orderRepository.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("%s: %w: %w", s.name, domain.ErrOrderNotFound, err)
	}
	if order.Status != domain.StatusAwaitingPayment {
		return fmt.Errorf("%s: %w", s.name, domain.ErrStatusToPay)
	}

	if err = s.stockRepository.ReserveCancel(ctx, order.Items); err != nil {
		return fmt.Errorf("%s: %w: %w", s.name, domain.ErrReserveCancelOrder, err)
	}

	if err = s.orderRepository.SetStatus(ctx, orderID, domain.StatusCancelled); err != nil {
		return fmt.Errorf("%s: %w: %w", s.name, domain.ErrStatusSetter, err)
	}

	return nil
}
