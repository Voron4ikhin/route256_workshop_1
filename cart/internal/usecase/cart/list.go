package cart

import (
	"context"
	"fmt"
	"route256/cart/internal/domain"
	"route256/cart/internal/port/out"
)

type ListService struct {
	name           string
	productClient  out.ProductClient
	cartRepository out.CartRepository
}

func NewListService(productClient out.ProductClient, cartRepository out.CartRepository) *ListService {
	return &ListService{
		name:           "cart.ListService",
		productClient:  productClient,
		cartRepository: cartRepository,
	}
}

func (s *ListService) GetList(ctx context.Context, user int64) ([]domain.FullCartItem, uint32, error) {
	items, err := s.cartRepository.GetList(ctx, user)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", s.name, err)
	}

	result := make([]domain.FullCartItem, 0, len(items))
	var totalPrice uint32
	for _, item := range items {
		name, price, err := s.productClient.GetProductInfo(ctx, item.SKU)
		if err != nil {
			return nil, 0, fmt.Errorf("%s: get product info for sku %d: %w", s.name, item.SKU, err)
		}
		result = append(result, domain.FullCartItem{
			SKU:   item.SKU,
			Count: item.Count,
			Name:  name,
			Price: price,
		})
		totalPrice += price * uint32(item.Count)
	}

	return result, totalPrice, nil
}
