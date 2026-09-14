package inmemory

import (
	"context"
	"errors"
	"log"
	"route256/cart/internal/pkg/handlers"
	"sync"
)

var ErrUserNotFound = errors.New("user not found")
var ErrSKUNotFound = errors.New("sku not found")

type InMemoryRepo struct {
	mu      sync.RWMutex
	storage map[int64]map[uint32]uint16
}

func NewRepository() (*InMemoryRepo, error) {
	return &InMemoryRepo{
		storage: make(map[int64]map[uint32]uint16),
	}, nil
}

func (r *InMemoryRepo) Add(ctx context.Context, user int64, sku uint32, count uint16) error {
	log.Println("InMemoryRepo.Add")
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.storage[user]; !ok {
		r.storage[user] = make(map[uint32]uint16)
	}
	r.storage[user][sku] += count
	return nil
}

func (r *InMemoryRepo) Delete(ctx context.Context, user int64, sku uint32) error {
	log.Println("InMemoryRepo.Delete")
	r.mu.Lock()
	defer r.mu.Unlock()
	items, ok := r.storage[user]
	if !ok {
		return ErrUserNotFound
	}
	if _, ok = items[sku]; !ok {
		return ErrSKUNotFound
	}
	delete(items, sku)
	return nil
}

func (r *InMemoryRepo) GetList(ctx context.Context, user int64) ([]handlers.FullCartItem, error) {
	log.Println("InMemoryRepo.GetList")
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := r.storage[user]
	result := make([]handlers.FullCartItem, 0, len(items))
	for sku, count := range items {
		result = append(result, handlers.FullCartItem{SKU: sku, Count: count})
	}
	return result, nil

}

func (r *InMemoryRepo) Clear(ctx context.Context, user int64) error {
	log.Println("InMemoryRepo.Clear")
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.storage, user)
	return nil
}
