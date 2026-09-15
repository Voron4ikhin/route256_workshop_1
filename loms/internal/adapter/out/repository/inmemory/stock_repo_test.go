package inmemory_test

import (
	"context"
	"route256/loms/internal/adapter/out/repository/inmemory"
	"route256/loms/internal/domain"
	"testing"
)

func TestStockRepository_ReserveThenCancel_RestoresAvailability(t *testing.T) {
	repo := inmemory.NewStockRepository()
	const sku = 1000001

	before := repo.GetBySKU(context.Background(), sku)

	items := []domain.OrderItem{{SKU: sku, Count: 5}}
	if err := repo.Reserve(context.Background(), items); err != nil {
		t.Fatalf("reserve failed: %v", err)
	}
	if got := repo.GetBySKU(context.Background(), sku); got != before-5 {
		t.Fatalf("expected available %d after reserve, got %d", before-5, got)
	}

	if err := repo.ReserveCancel(context.Background(), items); err != nil {
		t.Fatalf("cancel failed: %v", err)
	}
	if got := repo.GetBySKU(context.Background(), sku); got != before {
		t.Fatalf("expected available %d after cancel, got %d", before, got)
	}
}

func TestStockRepository_ReserveCancel_RejectsOverCancel(t *testing.T) {
	repo := inmemory.NewStockRepository()
	const sku = 1000001

	items := []domain.OrderItem{{SKU: sku, Count: 5}}
	if err := repo.Reserve(context.Background(), items); err != nil {
		t.Fatalf("reserve failed: %v", err)
	}

	overCancel := []domain.OrderItem{{SKU: sku, Count: 10}}
	if err := repo.ReserveCancel(context.Background(), overCancel); err == nil {
		t.Fatal("expected error when cancelling more than reserved, got nil")
	}

	if got := repo.GetBySKU(context.Background(), sku); got == 0 {
		t.Fatalf("reserved counter must not have underflowed, available stock became %d", got)
	}
}
