package order

import (
	"context"
	"errors"
	"fmt"
	"route256/loms/internal/domain"
	portin "route256/loms/internal/port/in"
	"route256/loms/internal/port/out"
	"time"
)

type AutoCancelService struct {
	name            string
	orderRepository out.OrderRepository
	canceler        portin.OrderCanceler
	unpaidTTL       time.Duration
}

func NewAutoCancelService(orderRepository out.OrderRepository, canceler portin.OrderCanceler, unpaidTTL time.Duration) *AutoCancelService {
	return &AutoCancelService{
		name:            "order.AutoCancelService",
		orderRepository: orderRepository,
		canceler:        canceler,
		unpaidTTL:       unpaidTTL,
	}
}

func (s *AutoCancelService) SweepUnpaidOrders(ctx context.Context) (int, error) {
	orders, err := s.orderRepository.ListByStatus(ctx, domain.StatusAwaitingPayment)
	if err != nil {
		return 0, fmt.Errorf("%s: list awaiting-payment orders: %w", s.name, err)
	}

	cutoff := time.Now().Add(-s.unpaidTTL)

	var cancelled int
	var errs error
	for _, o := range orders {
		if o.CreatedAt.After(cutoff) {
			continue
		}

		if err := s.canceler.CancelOrder(ctx, o.ID); err != nil {
			if errors.Is(err, domain.ErrStatusToPay) {
				continue
			}
			errs = errors.Join(errs, fmt.Errorf("cancel order %d: %w", o.ID, err))
			continue
		}
		cancelled++
	}

	return cancelled, errs
}
