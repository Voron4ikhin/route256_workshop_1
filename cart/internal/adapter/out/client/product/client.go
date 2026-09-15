package product

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"route256/cart/internal/authctx"
	"time"
)

type Client struct {
	name       string
	path       string
	httpClient *http.Client
}

type getProductRequest struct {
	Token string `json:"token,omitempty"`
	SKU   uint32 `json:"sku,omitempty"`
}

type getProductResponse struct {
	Name  string `json:"name,omitempty"`
	Price uint32 `json:"price,omitempty"`
}

type getProductErrorResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func New(name string, basePath string) (*Client, error) {
	const handlerName = "get_product"
	path, err := url.JoinPath(basePath, handlerName)
	if err != nil {
		return nil, fmt.Errorf("%s: incorrect base path: %w", name, err)
	}

	return &Client{
		name:       name,
		path:       path,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}, nil
}

// GetProductInfo authenticates with the token carried on ctx, which the
// Auth middleware populates from the incoming request's Authorization
// header (see internal/authctx and internal/adapter/in/http/middleware).
func (c *Client) GetProductInfo(ctx context.Context, sku uint32) (string, uint32, error) {
	request := getProductRequest{
		Token: authctx.Token(ctx),
		SKU:   sku,
	}
	data, err := json.Marshal(request)
	if err != nil {
		return "", 0, fmt.Errorf("%s: failed to marshal request: %w", c.name, err)
	}
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.path, bytes.NewBuffer(data))
	if err != nil {
		return "", 0, fmt.Errorf("%s: failed to create http request: %w", c.name, err)
	}
	httpResponse, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return "", 0, fmt.Errorf("%s: failed to do http request: %w", c.name, err)
	}
	defer func() {
		_ = httpResponse.Body.Close()
	}()

	if httpResponse.StatusCode != http.StatusOK {
		errResp := &getProductErrorResponse{}
		if err := json.NewDecoder(httpResponse.Body).Decode(errResp); err != nil {
			return "", 0, fmt.Errorf("%s: failed to decode http response: %w", c.name, err)
		}
		return "", 0, fmt.Errorf("%s: http request responded with %d: %s", c.name, httpResponse.StatusCode, errResp.Message)
	}

	response := &getProductResponse{}
	if err := json.NewDecoder(httpResponse.Body).Decode(response); err != nil {
		return "", 0, fmt.Errorf("%s: failed to decode http response: %w", c.name, err)
	}
	return response.Name, response.Price, nil
}
