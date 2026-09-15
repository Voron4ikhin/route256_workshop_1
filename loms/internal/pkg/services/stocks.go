package services

import "context"

type StocksProvider interface {
	GetBySKU(ctx context.Context, sku uint32) uint64
}

type StockService struct {
	stocksProvider StocksProvider
}

func NewStocksService(stocksProvider StocksProvider) *StockService {
	return &StockService{
		stocksProvider: stocksProvider,
	}
}

func (s *StockService) GetStocks(ctx context.Context, sku uint32) uint64 {
	return s.stocksProvider.GetBySKU(ctx, sku)
}
