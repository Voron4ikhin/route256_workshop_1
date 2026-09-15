package order

import (
	httpadapter "route256/loms/internal/adapter/in/http/httpkit"
	"route256/loms/internal/domain"
)

type ItemDTO struct {
	SKU   uint32 `json:"sku,omitempty"`
	Count uint16 `json:"count,omitempty"`
}

func itemsToDomain(items []ItemDTO) []domain.OrderItem {
	result := make([]domain.OrderItem, len(items))
	for i, item := range items {
		result[i] = domain.OrderItem{SKU: item.SKU, Count: item.Count}
	}
	return result
}

func itemsFromDomain(items []domain.OrderItem) []ItemDTO {
	result := make([]ItemDTO, len(items))
	for i, item := range items {
		result[i] = ItemDTO{SKU: item.SKU, Count: item.Count}
	}
	return result
}

type CreateRequest struct {
	User  int64     `json:"user,omitempty"`
	Items []ItemDTO `json:"items,omitempty"`
}

func (r CreateRequest) Validate() error {
	if r.User <= 0 {
		return httpadapter.ErrIncorrectUser
	}
	if len(r.Items) == 0 {
		return httpadapter.ErrIncorrectQuantity
	}
	for _, item := range r.Items {
		if item.SKU == 0 {
			return httpadapter.ErrIncorrectSKU
		}
		if item.Count <= 0 {
			return httpadapter.ErrIncorrectProductCount
		}
	}
	return nil
}

type CreateResponse struct {
	OrderID int64 `json:"orderID,omitempty"`
}

type InfoRequest struct {
	OrderID int64 `json:"orderID,omitempty"`
}

func (r InfoRequest) Validate() error {
	if r.OrderID <= 0 {
		return httpadapter.ErrIncorrectOrderID
	}
	return nil
}

type InfoResponse struct {
	Status string    `json:"status,omitempty"`
	User   int64     `json:"user,omitempty"`
	Items  []ItemDTO `json:"items,omitempty"`
}

func infoResponseFromDomain(o *domain.Order) InfoResponse {
	return InfoResponse{
		Status: string(o.Status),
		User:   o.UserID,
		Items:  itemsFromDomain(o.Items),
	}
}

type PayRequest struct {
	OrderID int64 `json:"orderID,omitempty"`
}

func (r PayRequest) Validate() error {
	if r.OrderID <= 0 {
		return httpadapter.ErrIncorrectOrderID
	}
	return nil
}

type CancelRequest struct {
	OrderID int64 `json:"orderID,omitempty"`
}

func (r CancelRequest) Validate() error {
	if r.OrderID <= 0 {
		return httpadapter.ErrIncorrectOrderID
	}
	return nil
}
