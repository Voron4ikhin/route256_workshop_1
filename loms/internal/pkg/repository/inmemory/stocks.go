package inmemory

import (
	"context"
	"fmt"
	"route256/loms/internal/pkg/handlers/orders"
	"sync"
)

type StocksStorage struct {
	stocks map[uint32]uint64
	mu     sync.RWMutex
}

func NewStocksStorage() (*StocksStorage, error) {
	return &StocksStorage{stocks: seedStocks()}, nil
}

func (s *StocksStorage) GetBySKU(sku uint32) uint64 {
	return s.stocks[sku]
}

func (s *StocksStorage) Reserve(ctx context.Context, items []orders.OrderItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, v := range items {
		count, ok := s.stocks[v.SKU]
		if !ok {
			return fmt.Errorf("stock not found item with %d", v.SKU)
		}
		if count < uint64(v.Count) {
			return fmt.Errorf("stock item with %d less then need to reserve", v.SKU)
		}
	}

	for _, v := range items {
		s.stocks[v.SKU] -= uint64(v.Count)
	}

	return nil
}
