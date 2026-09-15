package in

import "context"

type StockInformant interface {
	GetStocks(ctx context.Context, sku uint32) uint64
}
