package services

import (
	"context"
	"route256/cart/internal/pkg/handlers"
)

type ProductProvider interface {
	GetProductInfo(ctx context.Context, sku uint32) (string, uint32, error)
}

type CartLister interface {
	GetList(ctx context.Context, user int64) ([]handlers.FullCartItem, error)
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
	//TODO: здесь запрос на получение цен будет (когда пофиксят проблему с product_service)
	resp := handlers.ListResponse{
		Item:       make([]handlers.FullCartItem, 0),
		TotalPrice: 123, // будет 0 по умолчанию
	}
	list, err := s.cartLister.GetList(ctx, user)
	if err != nil {
		return resp, err
	}
	resp.Item = append(resp.Item, list...)
	return resp, nil
}
