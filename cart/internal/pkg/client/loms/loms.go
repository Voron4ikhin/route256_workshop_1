package loms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"route256/cart/internal/pkg/handlers"
)

type Client struct {
	name            string
	stockInfoPath   string
	orderCreatePath string
}

type StocksRequest struct {
	SKU uint32 `json:"sku,omitempty"`
}

type StocksResponse struct {
	Count uint64 `json:"count,omitempty"`
}

type OrderCreateRequest struct {
	User  int64               `json:"user,omitempty"`
	Items []handlers.CartItem `json:"items,omitempty"`
}

type OrderCreateResponse struct {
	OrderId int64 `json:"order_id,omitempty"`
}

func New(name string, basePath string) (*Client, error) {
	stockInfoPath, err := url.JoinPath(basePath, "stocks/info")
	if err != nil {
		return nil, fmt.Errorf("%s: incorrect base path: %w", name, err)
	}

	orderCreatePath, err := url.JoinPath(basePath, "order/create")
	if err != nil {
		return nil, fmt.Errorf("%s: incorrect base path: %w", name, err)
	}

	return &Client{
		name:            name,
		stockInfoPath:   stockInfoPath,
		orderCreatePath: orderCreatePath,
	}, nil
}

func (c Client) doJSONPost(ctx context.Context, path string, request, response any) error {
	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("%s: failed to marshal request: %w", c.name, err)
	}
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("%s: failed to create http request: %w", c.name, err)
	}
	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return fmt.Errorf("%s: failed to do http request: %w", c.name, err)
	}
	defer func() {
		_ = httpResponse.Body.Close()
	}()

	if httpResponse.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: http request failed with status code: %d", c.name, httpResponse.StatusCode)
	}

	if err := json.NewDecoder(httpResponse.Body).Decode(response); err != nil {
		return fmt.Errorf("%s: failed to decode http response: %w", c.name, err)
	}

	return nil
}

func (c Client) GetStocks(ctx context.Context, sku uint32) (uint64, error) {
	request := StocksRequest{
		SKU: sku,
	}
	response := &StocksResponse{}
	if err := c.doJSONPost(ctx, c.stockInfoPath, request, response); err != nil {
		return 0, err
	}

	return response.Count, nil
}

func (c Client) Checkout(ctx context.Context, user int64, items []handlers.CartItem) (int64, error) {
	request := OrderCreateRequest{
		User:  user,
		Items: items,
	}
	response := &OrderCreateResponse{}
	if err := c.doJSONPost(ctx, c.orderCreatePath, request, response); err != nil {
		return 0, err
	}

	return response.OrderId, nil
}
