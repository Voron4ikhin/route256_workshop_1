package mock

import (
	"context"
	"fmt"
)

type Client struct {
	name string
}

func New(name string) *Client {
	return &Client{name: name}
}

func (c *Client) GetProductInfo(_ context.Context, sku uint32) (string, uint32, error) {
	return fmt.Sprintf("mock product %d", sku), 100 + sku%900, nil
}
