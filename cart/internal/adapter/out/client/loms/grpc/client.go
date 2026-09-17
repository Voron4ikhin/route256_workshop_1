package grpc

import (
	"context"
	"fmt"
	"route256/cart/internal/domain"
	lomsv1 "route256/loms/pkg/api/loms/v1"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const requestTimeout = 5 * time.Second

type Client struct {
	name string
	conn *grpc.ClientConn
	api  lomsv1.LOMSServiceClient
}

func New(name string, addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to dial %s: %w", name, addr, err)
	}

	return &Client{
		name: name,
		conn: conn,
		api:  lomsv1.NewLOMSServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetStocks(ctx context.Context, sku uint32) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	resp, err := c.api.GetStocks(ctx, &lomsv1.GetStocksRequest{Sku: sku})
	if err != nil {
		return 0, fmt.Errorf("%s: get stocks: %w", c.name, err)
	}

	return resp.GetCount(), nil
}

func (c *Client) Checkout(ctx context.Context, user int64, items []domain.CartItem) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	pbItems := make([]*lomsv1.Item, len(items))
	for i, item := range items {
		pbItems[i] = &lomsv1.Item{Sku: item.SKU, Count: uint32(item.Count)}
	}

	resp, err := c.api.CreateOrder(ctx, &lomsv1.CreateOrderRequest{User: user, Items: pbItems})
	if err != nil {
		return 0, fmt.Errorf("%s: create order: %w", c.name, err)
	}

	return resp.GetOrderId(), nil
}
