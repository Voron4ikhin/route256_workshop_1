package item

import (
	"context"
	"encoding/json"
	"net/http"
	"route256/cart/internal/pkg/handlers"
)

func (r DeleteRequest) Validate() error {
	if r.User <= 0 {
		return handlers.ErrIncorrectUser
	}
	if r.SKU == 0 {
		return handlers.ErrIncorrectSKU
	}
	return nil
}

type DeleteRequest struct {
	User int64  `json:"user,omitempty"`
	SKU  uint32 `json:"sku,omitempty"`
}

type Deleter interface {
	Delete(ctx context.Context, user int64, sku uint32) error
}

type DeleteHandler struct {
	name    string
	deleter Deleter
}

func NewDeleteHandler(deleter Deleter) *DeleteHandler {
	return &DeleteHandler{
		name:    "item delete handler",
		deleter: deleter,
	}
}

func (h DeleteHandler) Handle(w http.ResponseWriter, r *http.Request) {
	req := &DeleteRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handlers.GetErrorResponse(w, h.name, err, http.StatusInternalServerError)
		return
	}
	if err := req.Validate(); err != nil {
		handlers.GetErrorResponse(w, h.name, err, http.StatusBadRequest)
		return
	}

	if err := h.deleter.Delete(r.Context(), req.User, req.SKU); err != nil {
		handlers.GetErrorResponse(w, h.name, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
