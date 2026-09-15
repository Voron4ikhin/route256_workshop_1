package orders

import (
	"context"
	"errors"
	"fmt"
	"route256/loms/internal/pkg/handlers/orders"
)

type OrderCreator interface {
	Create(ctx context.Context, user int64, items []orders.OrderItem) (int64, error)
}

type OrderStatusSetter interface {
	SetStatus(ctx context.Context, orderID int64, status OrderStatus) error
}

type StockReserver interface {
	Reserve(ctx context.Context, item []orders.OrderItem) error
}

type OrderStatus string

func (s OrderStatus) IsValid() bool {
	switch s {
	case StatusNew, StatusAwaitingPayment, StatusFailed, StatusCancelled, StatusPayed:
		return true
	default:
		return false
	}
}

const (
	StatusNew             OrderStatus = "new"
	StatusAwaitingPayment OrderStatus = "awaiting_payment"
	StatusFailed          OrderStatus = "failed"
	StatusCancelled       OrderStatus = "cancelled"
	StatusPayed           OrderStatus = "payed"
)

type OrderCreateService struct {
	name              string
	orderCreator      OrderCreator
	orderStatusSetter OrderStatusSetter
	stockReserver     StockReserver
}

var ErrCreateOrder = errors.New("cannot create order in OrderStorage")
var ErrOrderNotFound = errors.New("cannot find order in OrderStorage")
var ErrReserveOrder = errors.New("cannot reserve items in StockStorage")
var ErrReserveRemoveOrder = errors.New("cannot remove reserved items in StockStorage")
var ErrStatusSetter = errors.New("cannot set status in OrderStorage")
var ErrStatusToPay = errors.New("cannot set status to pay because status is not ready for pay")
var ErrReserveCancelOrder = errors.New("cannot cancel reserve order")

func NewOrderCreateService(orderCreator OrderCreator, orderStatusSetter OrderStatusSetter, stockReserver StockReserver) *OrderCreateService {
	return &OrderCreateService{
		name:              "OrderCreateService",
		orderCreator:      orderCreator,
		orderStatusSetter: orderStatusSetter,
		stockReserver:     stockReserver,
	}
}

func (s *OrderCreateService) CreateOrder(ctx context.Context, user int64, items []orders.OrderItem) (*orders.OrderCreateResponse, error) {
	resp := &orders.OrderCreateResponse{}
	orderID, err := s.orderCreator.Create(ctx, user, items)
	if err != nil {
		return resp, fmt.Errorf("%s: %s", s.name, ErrCreateOrder)
	}
	if err = s.stockReserver.Reserve(ctx, items); err != nil {
		if err = s.orderStatusSetter.SetStatus(ctx, orderID, StatusFailed); err != nil {
			return resp, fmt.Errorf("%s: %s", s.name, ErrStatusSetter)
		}
		return resp, fmt.Errorf("%s: %s", s.name, ErrReserveOrder)
	}
	if err = s.orderStatusSetter.SetStatus(ctx, orderID, StatusAwaitingPayment); err != nil {
		return resp, fmt.Errorf("%s: %s", s.name, ErrStatusSetter)
	}

	resp.OrderId = orderID
	return resp, nil
}
