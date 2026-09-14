package handlers

import (
	"context"
	"encoding/json"
	"net/http"
)

type CheckoutRequest struct {
	User int64
}

func (cr CheckoutRequest) Validate() error {
	if cr.User < 0 {
		return ErrIncorrectUser
	}
	return nil
}

type Checkouter interface {
	Checkout(ctx context.Context, user int64) (int64, error)
}

type CheckoutHandler struct {
	name       string
	checkouter Checkouter
}

func NewCheckoutHandler(checkouter Checkouter) *CheckoutHandler {
	return &CheckoutHandler{
		name:       "checkout handler",
		checkouter: checkouter,
	}
}

func (h *CheckoutHandler) Handle(w http.ResponseWriter, r *http.Request) {
	req := &CheckoutRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		GetErrorResponse(w, h.name, err, http.StatusInternalServerError)
		return
	}
	if err := req.Validate(); err != nil {
		GetErrorResponse(w, h.name, err, http.StatusBadRequest)
		return
	}

	order, err := h.checkouter.Checkout(r.Context(), req.User)
	if err != nil {
		GetErrorResponse(w, h.name, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(order); err != nil {
		GetErrorResponse(w, h.name, err, http.StatusInternalServerError)
		return
	}
}
