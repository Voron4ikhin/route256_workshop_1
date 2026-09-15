package stock

import (
	"context"
	"route256/loms/internal/port/out"
)

type QueryService struct {
	stockRepository out.StockRepository
}

func NewQueryService(stockRepository out.StockRepository) *QueryService {
	return &QueryService{
		stockRepository: stockRepository,
	}
}

func (s *QueryService) GetStocks(ctx context.Context, sku uint32) uint64 {
	return s.stockRepository.GetBySKU(ctx, sku)
}
