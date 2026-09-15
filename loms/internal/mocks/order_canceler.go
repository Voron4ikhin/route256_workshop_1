package mocks

import "context"

type OrderCanceler struct {
	CancelOrderFunc func(ctx context.Context, orderID int64) error
}

func (m *OrderCanceler) CancelOrder(ctx context.Context, orderID int64) error {
	return m.CancelOrderFunc(ctx, orderID)
}
