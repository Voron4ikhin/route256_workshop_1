package inmemory

import (
	"context"
	"route256/cart/internal/domain"
	"sync"
)

type CartRepository struct {
	mu      sync.RWMutex
	storage map[int64]map[uint32]uint16
}

func NewCartRepository() *CartRepository {
	return &CartRepository{
		storage: make(map[int64]map[uint32]uint16),
	}
}

func (r *CartRepository) Add(ctx context.Context, user int64, sku uint32, count uint16) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.storage[user]; !ok {
		r.storage[user] = make(map[uint32]uint16)
	}
	r.storage[user][sku] += count
	return nil
}

func (r *CartRepository) Delete(ctx context.Context, user int64, sku uint32) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	items, ok := r.storage[user]
	if !ok {
		return domain.ErrUserNotFound
	}
	if _, ok = items[sku]; !ok {
		return domain.ErrSKUNotFound
	}
	delete(items, sku)
	return nil
}

func (r *CartRepository) GetList(ctx context.Context, user int64) ([]domain.CartItem, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := r.storage[user]
	result := make([]domain.CartItem, 0, len(items))
	for sku, count := range items {
		result = append(result, domain.CartItem{SKU: sku, Count: count})
	}
	return result, nil
}

func (r *CartRepository) Clear(ctx context.Context, user int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.storage, user)
	return nil
}
