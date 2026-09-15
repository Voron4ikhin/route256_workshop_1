package inmemory

import (
	"context"
	"fmt"
	"route256/loms/internal/domain"
	"sync"
)

type stockCount struct {
	total    uint64
	reserved uint64
}

type StockRepository struct {
	mu     sync.RWMutex
	stocks map[uint32]*stockCount
}

func NewStockRepository() *StockRepository {
	seed := seedStocks()
	stocks := make(map[uint32]*stockCount, len(seed))
	for sku, count := range seed {
		stocks[sku] = &stockCount{total: count}
	}
	return &StockRepository{stocks: stocks}
}

func (r *StockRepository) GetBySKU(ctx context.Context, sku uint32) uint64 {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count, ok := r.stocks[sku]
	if !ok {
		return 0
	}

	return count.total - count.reserved
}

func (r *StockRepository) Reserve(ctx context.Context, items []domain.OrderItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, v := range items {
		count, ok := r.stocks[v.SKU]
		if !ok {
			return fmt.Errorf("stock not found item with %d", v.SKU)
		}
		if count.total-count.reserved < uint64(v.Count) {
			return fmt.Errorf("stock item with %d less then need to reserve", v.SKU)
		}
	}

	for _, v := range items {
		r.stocks[v.SKU].reserved += uint64(v.Count)
	}

	return nil
}

func (r *StockRepository) ReserveRemove(ctx context.Context, items []domain.OrderItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, v := range items {
		count, ok := r.stocks[v.SKU]
		if !ok {
			return fmt.Errorf("stock not found item with %d", v.SKU)
		}
		if count.reserved < uint64(v.Count) {
			return fmt.Errorf("stock item with %d has less reserved than need to remove", v.SKU)
		}
	}

	for _, v := range items {
		count := r.stocks[v.SKU]
		count.reserved -= uint64(v.Count)
		count.total -= uint64(v.Count)
	}

	return nil
}

func (r *StockRepository) ReserveCancel(ctx context.Context, items []domain.OrderItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, v := range items {
		count, ok := r.stocks[v.SKU]
		if !ok {
			return fmt.Errorf("stock not found item with %d", v.SKU)
		}
		if count.reserved < uint64(v.Count) {
			return fmt.Errorf("stock item with %d has less reserved than need to cancel", v.SKU)
		}
	}

	for _, v := range items {
		r.stocks[v.SKU].reserved -= uint64(v.Count)
	}

	return nil
}
