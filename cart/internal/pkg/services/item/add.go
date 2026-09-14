package item

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type StocksProvider interface {
	GetStocks(ctx context.Context, sku uint32) (uint64, error)
}

type ProductProvider interface {
	GetProductInfo(ctx context.Context, sku uint32) (string, uint32, error)
}

type CartAdder interface {
	Add(ctx context.Context, user int64, sku uint32, count uint16) error
}

type AddService struct {
	name            string
	stocksProvider  StocksProvider
	productProvider ProductProvider
	cartAdder       CartAdder
}

var ErrInsufficientStocks = errors.New("insufficient stocks")
var ErrAddItemToCart = errors.New("add item to cart")

func NewAddService(stocksProvider StocksProvider, productProvider ProductProvider, cartAdder CartAdder) *AddService {
	return &AddService{
		name:            "item add service",
		stocksProvider:  stocksProvider,
		productProvider: productProvider,
		cartAdder:       cartAdder,
	}
}

func (s AddService) Add(ctx context.Context, user int64, sku uint32, count uint16) error {
	if _, _, err := s.productProvider.GetProductInfo(ctx, sku); err != nil {
		log.Printf("failed to get product info: %v", err)
		//TODO: вернуть ошибку надо тут
		//return err
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

	if err := s.cartAdder.Add(ctx, user, sku, count); err != nil {
		return fmt.Errorf("%s: %s", s.name, ErrAddItemToCart)
	}

	return nil
}
