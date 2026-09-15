package item

import (
	"context"
	"fmt"
	"route256/cart/internal/domain"
	"route256/cart/internal/port/out"
)

type AddService struct {
	name           string
	productClient  out.ProductClient
	stockClient    out.StockClient
	cartRepository out.CartRepository
}

func NewAddService(productClient out.ProductClient, stockClient out.StockClient, cartRepository out.CartRepository) *AddService {
	return &AddService{
		name:           "item.AddService",
		productClient:  productClient,
		stockClient:    stockClient,
		cartRepository: cartRepository,
	}
}

func (s *AddService) Add(ctx context.Context, user int64, sku uint32, count uint16) error {
	if _, _, err := s.productClient.GetProductInfo(ctx, sku); err != nil {
		return fmt.Errorf("%s: get product info: %w", s.name, err)
	}

	stocksCount, err := s.stockClient.GetStocks(ctx, sku)
	if err != nil {
		return fmt.Errorf("%s: get stocks: %w", s.name, err)
	}
	if uint64(count) > stocksCount {
		return fmt.Errorf("%s: %w", s.name, domain.ErrInsufficientStocks)
	}

	if err := s.cartRepository.Add(ctx, user, sku, count); err != nil {
		return fmt.Errorf("%s: %w: %w", s.name, domain.ErrAddItemToCart, err)
	}

	return nil
}
