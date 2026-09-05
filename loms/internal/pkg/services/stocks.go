package services

type StocksProvider interface {
	GetStocks(sku uint32) uint64
}

type StockService struct {
	stocksProvider StocksProvider
}

func NewStocksService(stocksProvider StocksProvider) *StockService {
	return &StockService{
		stocksProvider: stocksProvider,
	}
}

func (s *StockService) GetStocks(sku uint32) uint64 {
	return s.stocksProvider.GetStocks(sku)
}
