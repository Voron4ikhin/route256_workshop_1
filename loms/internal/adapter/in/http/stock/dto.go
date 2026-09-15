package stock

import httpadapter "route256/loms/internal/adapter/in/http/httpkit"

type InfoRequest struct {
	SKU uint32 `json:"sku"`
}

func (r InfoRequest) Validate() error {
	if r.SKU == 0 {
		return httpadapter.ErrIncorrectSKU
	}
	return nil
}

type InfoResponse struct {
	Count uint64 `json:"count"`
}
