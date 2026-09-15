package stock_test

import (
	"context"
	"route256/loms/internal/mocks"
	"route256/loms/internal/usecase/stock"
	"testing"
)

func TestQueryService_GetStocks(t *testing.T) {
	repo := &mocks.StockRepository{
		GetBySKUFunc: func(ctx context.Context, sku uint32) uint64 {
			if sku != 1000001 {
				t.Fatalf("unexpected sku: %d", sku)
			}
			return 15
		},
	}

	svc := stock.NewQueryService(repo)
	got := svc.GetStocks(context.Background(), 1000001)
	if got != 15 {
		t.Fatalf("expected 15, got %d", got)
	}
}
