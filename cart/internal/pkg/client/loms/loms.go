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
	name string
	path string
}

type StocksRequest struct {
	SKU uint32 `json:"sku,omitempty"`
}

type StocksResponse struct {
	Count uint64 `json:"count,omitempty"`
}

func New(name string, basePath string) (*Client, error) {
	const stocksPath = "stocks"
	path, err := url.JoinPath(basePath, stocksPath)
	if err != nil {
		return nil, fmt.Errorf("%s: incorrect base path: %w", name, err)
	}

	return &Client{
		name: name,
		path: path,
	}, nil
}

func (c Client) GetStocks(ctx context.Context, sku uint32) (uint64, error) {
	request := StocksRequest{
		SKU: sku,
	}
	data, err := json.Marshal(request)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to marshal request: %w", c.name, err)
	}
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.path, bytes.NewBuffer(data))
	if err != nil {
		return 0, fmt.Errorf("%s: failed to create http request: %w", c.name, err)
	}
	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to do http request: %w", c.name, err)
	}
	defer func() {
		_ = httpResponse.Body.Close()
	}()

	if httpResponse.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("%s: http request failed with status code: %d", c.name, httpResponse.StatusCode)
	}

	response := &StocksResponse{}
	err = json.NewDecoder(httpResponse.Body).Decode(response)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to decode http response: %w", c.name, err)
	}

	return response.Count, nil
}

func (c Client) Checkout(ctx context.Context, user int64, items []handlers.FullCartItem) (int64, error) {
	return 12, nil
}
