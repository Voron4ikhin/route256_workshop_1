package item

import httpkit "route256/cart/internal/adapter/in/http/httpkit"

type AddRequest struct {
	User  int64  `json:"user,omitempty"`
	SKU   uint32 `json:"sku,omitempty"`
	Count uint16 `json:"count,omitempty"`
}

func (r AddRequest) Validate() error {
	if r.User <= 0 {
		return httpkit.ErrIncorrectUser
	}
	if r.SKU == 0 {
		return httpkit.ErrIncorrectSKU
	}
	if r.Count == 0 {
		return httpkit.ErrIncorrectQuantity
	}
	return nil
}

type DeleteRequest struct {
	User int64  `json:"user,omitempty"`
	SKU  uint32 `json:"sku,omitempty"`
}

func (r DeleteRequest) Validate() error {
	if r.User <= 0 {
		return httpkit.ErrIncorrectUser
	}
	if r.SKU == 0 {
		return httpkit.ErrIncorrectSKU
	}
	return nil
}
