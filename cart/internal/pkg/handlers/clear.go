package handlers

import (
	"context"
	"encoding/json"
	"net/http"
)

type ClearRequest struct {
	User int64 `json:"user,omitempty"`
}

func (cl ClearRequest) Validate() error {
	if cl.User < 0 {
		return ErrIncorrectUser
	}
	return nil
}

type CartClear interface {
	Clear(ctx context.Context, user int64) error
}

type ClearHandler struct {
	name      string
	cartClear CartClear
}

func NewClearHandler(cartClear CartClear) *ClearHandler {
	return &ClearHandler{
		name:      "cart clear handler",
		cartClear: cartClear,
	}
}

func (h ClearHandler) Handle(w http.ResponseWriter, r *http.Request) {
	req := &ClearRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		GetErrorResponse(w, h.name, err, http.StatusInternalServerError)
		return
	}
	if err := req.Validate(); err != nil {
		GetErrorResponse(w, h.name, err, http.StatusBadRequest)
		return
	}

	err := h.cartClear.Clear(r.Context(), req.User)
	if err != nil {
		GetErrorResponse(w, h.name, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
