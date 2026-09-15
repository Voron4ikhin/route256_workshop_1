package inmemory

import (
	"context"
	"fmt"
	"route256/loms/internal/pkg/handlers/orders"
	"sync"
)

type stockCount struct {
	total    uint64
	reserved uint64
}

type StocksStorage struct {
	stocks map[uint32]*stockCount
	mu     sync.RWMutex
}

func NewStocksStorage() (*StocksStorage, error) {
	seed := seedStocks()
	stocks := make(map[uint32]*stockCount, len(seed))
	for sku, count := range seed {
		stocks[sku] = &stockCount{total: count}
	}
	return &StocksStorage{stocks: stocks}, nil
}

func (s *StocksStorage) GetBySKU(sku uint32) uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count, ok := s.stocks[sku]
	if !ok {
		return 0
	}

	return count.total - count.reserved
}

func (s *StocksStorage) Reserve(ctx context.Context, items []orders.OrderItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, v := range items {
		count, ok := s.stocks[v.SKU]
		if !ok {
			return fmt.Errorf("stock not found item with %d", v.SKU)
		}
		if count.total-count.reserved < uint64(v.Count) {
			return fmt.Errorf("stock item with %d less then need to reserve", v.SKU)
		}
	}

	for _, v := range items {
		s.stocks[v.SKU].reserved += uint64(v.Count)
	}

	return nil
}

func (s *StocksStorage) ReserveRemove(ctx context.Context, items []orders.OrderItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, v := range items {
		count, ok := s.stocks[v.SKU]
		if !ok {
			return fmt.Errorf("stock not found item with %d", v.SKU)
		}
		if count.reserved < uint64(v.Count) {
			return fmt.Errorf("stock item with %d has less reserved than need to remove", v.SKU)
		}
	}

	for _, v := range items {
		count := s.stocks[v.SKU]
		count.reserved -= uint64(v.Count)
		count.total -= uint64(v.Count)
	}

	return nil
}

func (s *StocksStorage) ReserveCancel(ctx context.Context, items []orders.OrderItem) error {
	return nil
}
