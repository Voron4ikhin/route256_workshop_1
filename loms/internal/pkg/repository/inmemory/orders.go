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
	ID     uint64
	UserID int64
	Items  []OrderItem
	Status sorders.OrderStatus
}

type OrdersStorage struct {
	mu     sync.RWMutex
	orders map[uint64]*Order
	nextID uint64
}

func NewOrdersStorage() (*OrdersStorage, error) {
	return &OrdersStorage{
		orders: make(map[uint64]*Order),
		nextID: 1,
	}, nil
}

func (s *OrdersStorage) Create(ctx context.Context, user int64, items []horders.OrderItem) (uint64, error) {
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

func (s *OrdersStorage) SetStatus(ctx context.Context, orderID uint64, status sorders.OrderStatus) error {
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
