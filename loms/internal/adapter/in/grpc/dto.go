package grpc

import (
	"route256/loms/internal/domain"
	lomsv1 "route256/loms/pkg/api/loms/v1"
)

func itemsToDomain(items []*lomsv1.Item) []domain.OrderItem {
	result := make([]domain.OrderItem, len(items))
	for i, item := range items {
		result[i] = domain.OrderItem{SKU: item.GetSku(), Count: uint16(item.GetCount())}
	}
	return result
}

func itemsFromDomain(items []domain.OrderItem) []*lomsv1.Item {
	result := make([]*lomsv1.Item, len(items))
	for i, item := range items {
		result[i] = &lomsv1.Item{Sku: item.SKU, Count: uint32(item.Count)}
	}
	return result
}
