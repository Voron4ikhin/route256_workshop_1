package orders

import (
	"context"
	"fmt"
	"route256/loms/internal/pkg/handlers/orders"
)

type StockReserveRemover interface {
	ReserveRemove(ctx context.Context, items []orders.OrderItem) error
}

type OrderPayService struct {
	name                string
	orderInformant      OrderInformant
	stockReserveRemover StockReserveRemover
	orderStatusSetter   OrderStatusSetter
}

func NewOrderPayService(orderInformant OrderInformant, stockReserveRemover StockReserveRemover, orderStatusSetter OrderStatusSetter) *OrderPayService {
	return &OrderPayService{
		name:                "OrderPayService",
		orderInformant:      orderInformant,
		stockReserveRemover: stockReserveRemover,
		orderStatusSetter:   orderStatusSetter,
	}
}

func (s *OrderPayService) PayOrder(ctx context.Context, orderID int64) error {
	order, err := s.orderInformant.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("%s: %s", s.name, ErrOrderNotFound)
	}
	if order.Status != StatusAwaitingPayment {
		return fmt.Errorf("%s: %s", s.name, ErrStatusToPay)
	}

	if err = s.stockReserveRemover.ReserveRemove(ctx, order.Items); err != nil {
		return fmt.Errorf("%s: %s", s.name, ErrReserveRemoveOrder)
	}

	if err = s.orderStatusSetter.SetStatus(ctx, orderID, StatusPayed); err != nil {
		return fmt.Errorf("%s: %s", s.name, ErrStatusSetter)
	}

	return nil
}
