package cart

import (
	httpkit "route256/cart/internal/adapter/in/http/httpkit"
	"route256/cart/internal/domain"
)

type ItemDTO struct {
	SKU   uint32
	Count uint16
	Name  string
	Price uint32
}

func itemsFromDomain(items []domain.FullCartItem) []ItemDTO {
	result := make([]ItemDTO, len(items))
	for i, item := range items {
		result[i] = ItemDTO{SKU: item.SKU, Count: item.Count, Name: item.Name, Price: item.Price}
	}
	return result
}

type ListRequest struct {
	User int64
}

func (r ListRequest) Validate() error {
	if r.User < 0 {
		return httpkit.ErrIncorrectUser
	}
	return nil
}

type ListResponse struct {
	Item       []ItemDTO
	TotalPrice uint32
}

type ClearRequest struct {
	User int64 `json:"user,omitempty"`
}

func (r ClearRequest) Validate() error {
	if r.User < 0 {
		return httpkit.ErrIncorrectUser
	}
	return nil
}

type CheckoutRequest struct {
	User int64 `json:"user,omitempty"`
}

func (r CheckoutRequest) Validate() error {
	if r.User < 0 {
		return httpkit.ErrIncorrectUser
	}
	return nil
}
