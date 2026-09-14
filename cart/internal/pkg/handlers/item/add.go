package item

import (
	"context"
	"encoding/json"
	"net/http"
	"route256/cart/internal/pkg/handlers"
)

func (r AddRequest) Validate() error {
	if r.User <= 0 {
		return handlers.ErrIncorrectUser
	}
	if r.SKU == 0 {
		return handlers.ErrIncorrectSKU
	}
	if r.Count == 0 {
		return handlers.ErrIncorrectQuantity
	}
	return nil
}

type AddRequest struct {
	User  int64  `json:"user,omitempty"`
	SKU   uint32 `json:"sku,omitempty"`
	Count uint16 `json:"count,omitempty"`
}

type Adder interface {
	Add(ctx context.Context, user int64, sku uint32, count uint16) error
}

type AddHandler struct {
	name  string
	adder Adder
}

func NewAddHandler(adder Adder) *AddHandler {
	return &AddHandler{
		name:  "item add handler",
		adder: adder,
	}
}

func (h AddHandler) Handle(w http.ResponseWriter, r *http.Request) {
	req := &AddRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handlers.GetErrorResponse(w, h.name, err, http.StatusInternalServerError)
		return
	}
	if err := req.Validate(); err != nil {
		handlers.GetErrorResponse(w, h.name, err, http.StatusBadRequest)
		return
	}

	if err := h.adder.Add(r.Context(), req.User, req.SKU, req.Count); err != nil {
		handlers.GetErrorResponse(w, h.name, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
