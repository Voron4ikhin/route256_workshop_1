package orders

import (
	"context"
	"fmt"
	"route256/loms/internal/pkg/handlers/orders"
)

type StockReserveCanceler interface {
	ReserveCancel(ctx context.Context, items []orders.OrderItem) error
}

type OrderCancelService struct {
	name                 string
	orderInformant       OrderInformant
	stockReserveCanceler StockReserveCanceler
	orderStatusSetter    OrderStatusSetter
}

func NewOrderCancelService(orderInformant OrderInformant, stockReserveCanceler StockReserveCanceler, orderStatusSetter OrderStatusSetter) *OrderCancelService {
	return &OrderCancelService{
		name:                 "OrderCancelService",
		orderInformant:       orderInformant,
		stockReserveCanceler: stockReserveCanceler,
		orderStatusSetter:    orderStatusSetter,
	}
}

func (s *OrderCancelService) CancelOrder(ctx context.Context, orderID int64) error {
	order, err := s.orderInformant.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("%s: %s", s.name, ErrOrderNotFound)
	}
	if order.Status != StatusAwaitingPayment {
		return fmt.Errorf("%s: %s", s.name, ErrStatusToPay)
	}

	if err = s.stockReserveCanceler.ReserveCancel(ctx, order.Items); err != nil {
		return fmt.Errorf("%s: %s", s.name, ErrReserveCancelOrder)
	}

	if err = s.orderStatusSetter.SetStatus(ctx, orderID, StatusCancelled); err != nil {
		return fmt.Errorf("%s: %s", s.name, ErrStatusSetter)
	}

	return nil
}
