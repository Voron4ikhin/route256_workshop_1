package inmemory

import (
	"context"
	"fmt"
	"route256/loms/internal/domain"
	"sync"
	"time"
)

type OrderRepository struct {
	mu     sync.RWMutex
	orders map[int64]*domain.Order
	nextID int64
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		orders: make(map[int64]*domain.Order),
		nextID: 1,
	}
}

func (r *OrderRepository) Create(ctx context.Context, user int64, items []domain.OrderItem) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	orderID := r.nextID

	itemsCopy := make([]domain.OrderItem, len(items))
	copy(itemsCopy, items)

	r.orders[orderID] = &domain.Order{
		ID:        orderID,
		UserID:    user,
		Items:     itemsCopy,
		Status:    domain.StatusNew,
		CreatedAt: time.Now(),
	}
	r.nextID++

	return orderID, nil
}

func (r *OrderRepository) SetStatus(ctx context.Context, orderID int64, status domain.OrderStatus) error {
	if !status.IsValid() {
		return fmt.Errorf("not valid status: %s", status)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.orders[orderID]
	if !ok {
		return fmt.Errorf("order %d not found", orderID)
	}
	order.Status = status

	return nil
}

func (r *OrderRepository) GetByID(ctx context.Context, orderID int64) (*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[orderID]
	if !ok {
		return nil, fmt.Errorf("order %d not found", orderID)
	}

	itemsCopy := make([]domain.OrderItem, len(order.Items))
	copy(itemsCopy, order.Items)

	return &domain.Order{
		ID:        order.ID,
		UserID:    order.UserID,
		Items:     itemsCopy,
		Status:    order.Status,
		CreatedAt: order.CreatedAt,
	}, nil
}

func (r *OrderRepository) ListByStatus(ctx context.Context, status domain.OrderStatus) ([]*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Order
	for _, order := range r.orders {
		if order.Status != status {
			continue
		}

		itemsCopy := make([]domain.OrderItem, len(order.Items))
		copy(itemsCopy, order.Items)

		result = append(result, &domain.Order{
			ID:        order.ID,
			UserID:    order.UserID,
			Items:     itemsCopy,
			Status:    order.Status,
			CreatedAt: order.CreatedAt,
		})
	}

	return result, nil
}
