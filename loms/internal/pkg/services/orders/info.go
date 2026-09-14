package orders

import (
	"context"
	horders "route256/loms/internal/pkg/handlers/orders"
)

type Order struct {
	UserID int64
	Items  []horders.OrderItem
	Status OrderStatus
}

type OrderInformant interface {
	GetByID(ctx context.Context, orderID int64) (*Order, error)
}

type OrderInfoService struct {
	name           string
	orderInformant OrderInformant
}

func NewOrderInfoService(orderInformant OrderInformant) *OrderInfoService {
	return &OrderInfoService{
		name:           "OrderInfoService",
		orderInformant: orderInformant,
	}
}

func (s *OrderInfoService) GetOrderInfo(ctx context.Context, orderID int64) (*horders.OrderInfoResponse, error) {
	resp := &horders.OrderInfoResponse{}
	order, err := s.orderInformant.GetByID(ctx, orderID)
	if err != nil {
		return resp, err
	}
	resp.Status = string(order.Status)
	resp.User = order.UserID
	resp.Items = make([]horders.OrderItem, len(order.Items))
	for i, item := range order.Items {
		resp.Items[i] = horders.OrderItem{
			SKU:   item.SKU,
			Count: item.Count,
		}
	}

	return resp, nil
}
