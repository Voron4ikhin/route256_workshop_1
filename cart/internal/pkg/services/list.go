package services

import (
	"context"
	"fmt"
	"route256/cart/internal/pkg/handlers"
)

type ProductProvider interface {
	GetProductInfo(ctx context.Context, sku uint32) (string, uint32, error)
}

type CartLister interface {
	GetList(ctx context.Context, user int64) ([]handlers.CartItem, error)
}

type ListService struct {
	name            string
	productProvider ProductProvider
	cartLister      CartLister
}

func NewListService(productProvider ProductProvider, cartLister CartLister) *ListService {
	return &ListService{
		name:            "cart list service",
		productProvider: productProvider,
		cartLister:      cartLister,
	}
}

func (s ListService) GetList(ctx context.Context, user int64) (handlers.ListResponse, error) {
	resp := handlers.ListResponse{
		Item: make([]handlers.FullCartItem, 0),
	}

	items, err := s.cartLister.GetList(ctx, user)
	if err != nil {
		return resp, err
	}

	var totalPrice uint32
	for _, item := range items {
		name, price, err := s.productProvider.GetProductInfo(ctx, item.SKU)
		if err != nil {
			return handlers.ListResponse{}, fmt.Errorf("%s: get product info for sku %d: %w", s.name, item.SKU, err)
		}
		resp.Item = append(resp.Item, handlers.FullCartItem{
			SKU:   item.SKU,
			Count: item.Count,
			Name:  name,
			Price: price,
		})
		totalPrice += price * uint32(item.Count)
	}
	resp.TotalPrice = totalPrice

	return resp, nil
}
