package product

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	name string
	path string
}

type GetProductRequest struct {
	Token string `json:"token,omitempty"`
	SKU   uint32 `json:"sku,omitempty"`
}

type GetProductResponse struct {
	Name  string `json:"name,omitempty"`
	Price uint32 `json:"price,omitempty"`
}

type GetProductErrorResponse struct {
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
		name: name,
		path: path,
	}, nil
}

func (c Client) GetProductInfo(sku uint32) (string, uint32, error) {
	request := GetProductRequest{
		Token: "testtoken", // TODO: get from request and pass into ctx
		SKU:   sku,
	}
	data, err := json.Marshal(request)
	if err != nil {
		return "", 0, fmt.Errorf("%s failed to marshal request: %w", c.name, err)
	}
	httpRequest, err := http.NewRequest(http.MethodPost, c.path, bytes.NewBuffer(data))
	if err != nil {
		return "", 0, fmt.Errorf("%s failed to create http request: %w", c.name, err)
	}
	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return "", 0, fmt.Errorf("%s failed to do http request: %w", c.name, err)
	}

	defer func() {
		_ = httpResponse.Body.Close()
	}()

	if httpResponse.StatusCode != http.StatusOK {
		response := &GetProductErrorResponse{}
		err = json.NewDecoder(httpResponse.Body).Decode(response)
		if err != nil {
			return "", 0, fmt.Errorf("%s failed to decode http response: %w", c.name, err)
		}
		return "", 0, fmt.Errorf("%s HTTP request responded with: %d, message: %s", c.name, httpResponse.StatusCode, response.Message)
	}
	response := &GetProductResponse{}
	err = json.NewDecoder(httpResponse.Body).Decode(response)
	if err != nil {
		return "", 0, fmt.Errorf("%s failed to decode http response: %w", c.name, err)
	}
	return response.Name, response.Price, nil
}
