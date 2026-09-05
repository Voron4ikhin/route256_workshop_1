package item

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type StocksProvider interface {
	GetStocks(ctx context.Context, sku uint32) (uint64, error)
}

type ProductProvider interface {
	GetProductInfo(sku uint32) (string, uint32, error)
}

type AddService struct {
	name            string
	stocksProvider  StocksProvider
	productProvider ProductProvider
}

func NewAddService(stocksProvider StocksProvider, productProvider ProductProvider) *AddService {
	return &AddService{
		name:            "item and service",
		stocksProvider:  stocksProvider,
		productProvider: productProvider,
	}
}

var ErrInsufficientStocks = errors.New("insufficient stocks")

func (s AddService) Add(ctx context.Context, user int64, sku uint32, count uint16) error {
	if _, _, err := s.productProvider.GetProductInfo(sku); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	stocksCount, err := s.stocksProvider.GetStocks(ctx, sku)
	if err != nil {
		return err
	}
	if uint64(count) > stocksCount {
		return fmt.Errorf("%s: %s", s.name, ErrInsufficientStocks)
	}
	return nil
}
