package inmemory

import (
	"context"
	"fmt"
	horders "route256/loms/internal/pkg/handlers/orders"
	sorders "route256/loms/internal/pkg/services/orders"
	"sync"
)

type OrderItem struct {
	SKU   uint32
	Count uint16
}

type Order struct {
	ID     int64
	UserID int64
	Items  []OrderItem
	Status sorders.OrderStatus
}

type OrdersStorage struct {
	mu     sync.RWMutex
	orders map[int64]*Order
	nextID int64
}

func NewOrdersStorage() (*OrdersStorage, error) {
	return &OrdersStorage{
		orders: make(map[int64]*Order),
		nextID: 1,
	}, nil
}

func (s *OrdersStorage) Create(ctx context.Context, user int64, items []horders.OrderItem) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	orderID := s.nextID

	itemsCopy := make([]OrderItem, len(items))
	for k, item := range items {
		itemsCopy[k] = OrderItem{
			SKU:   item.SKU,
			Count: item.Count,
		}
	}

	order := &Order{
		ID:     orderID,
		UserID: user,
		Items:  itemsCopy,
		Status: sorders.StatusNew,
	}

	s.orders[orderID] = order
	s.nextID++

	return orderID, nil
}

func (s *OrdersStorage) SetStatus(ctx context.Context, orderID int64, status sorders.OrderStatus) error {
	if !status.IsValid() {
		return fmt.Errorf("not valid status: %s", status)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	order, ok := s.orders[orderID]
	if !ok {
		return fmt.Errorf("order %d not found", orderID)
	}

	order.Status = status

	return nil
}

func (s *OrdersStorage) GetByID(ctx context.Context, orderID int64) (*sorders.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := &sorders.Order{}
	order, ok := s.orders[orderID]
	if !ok {
		return res, fmt.Errorf("order %d not found", orderID)
	}

	res.UserID = order.UserID
	res.Status = order.Status
	res.Items = make([]horders.OrderItem, len(order.Items))
	for i, item := range order.Items {
		res.Items[i] = horders.OrderItem{
			SKU:   item.SKU,
			Count: item.Count,
		}
	}

	return res, nil
}
