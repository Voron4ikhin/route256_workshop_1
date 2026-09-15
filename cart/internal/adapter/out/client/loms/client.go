package loms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"route256/cart/internal/domain"
	"time"
)

type Client struct {
	name            string
	stockInfoPath   string
	orderCreatePath string
	httpClient      *http.Client
}

type stocksRequest struct {
	SKU uint32 `json:"sku,omitempty"`
}

type stocksResponse struct {
	Count uint64 `json:"count,omitempty"`
}

type itemDTO struct {
	SKU   uint32 `json:"sku,omitempty"`
	Count uint16 `json:"count,omitempty"`
}

type orderCreateRequest struct {
	User  int64     `json:"user,omitempty"`
	Items []itemDTO `json:"items,omitempty"`
}

type orderCreateResponse struct {
	OrderID int64 `json:"orderID,omitempty"`
}

func New(name string, basePath string) (*Client, error) {
	stockInfoPath, err := url.JoinPath(basePath, "stock/info")
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
		httpClient:      &http.Client{Timeout: 5 * time.Second},
	}, nil
}

func (c *Client) doJSONPost(ctx context.Context, path string, request, response any) error {
	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("%s: failed to marshal request: %w", c.name, err)
	}
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("%s: failed to create http request: %w", c.name, err)
	}
	httpResponse, err := c.httpClient.Do(httpRequest)
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

func (c *Client) GetStocks(ctx context.Context, sku uint32) (uint64, error) {
	request := stocksRequest{SKU: sku}
	response := &stocksResponse{}
	if err := c.doJSONPost(ctx, c.stockInfoPath, request, response); err != nil {
		return 0, err
	}
	return response.Count, nil
}

func (c *Client) Checkout(ctx context.Context, user int64, items []domain.CartItem) (int64, error) {
	dtoItems := make([]itemDTO, len(items))
	for i, item := range items {
		dtoItems[i] = itemDTO{SKU: item.SKU, Count: item.Count}
	}

	request := orderCreateRequest{User: user, Items: dtoItems}
	response := &orderCreateResponse{}
	if err := c.doJSONPost(ctx, c.orderCreatePath, request, response); err != nil {
		return 0, err
	}
	return response.OrderID, nil
}
