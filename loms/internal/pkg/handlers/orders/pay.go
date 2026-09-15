package orders

import (
	"context"
	"encoding/json"
	"net/http"
	"route256/loms/internal/pkg/handlers"
)

type OrderPayRequest struct {
	OrderID int64 `json:"order_id,omitempty"`
}

func (r OrderPayRequest) Validate() error {
	if r.OrderID <= 0 {
		return handlers.ErrIncorrectOrderID
	}
	return nil
}

type OrderPayer interface {
	PayOrder(ctx context.Context, orderID int64) error
}

type OrderPayHandler struct {
	name       string
	orderPayer OrderPayer
}

func NewOrderPayHandler(orderPayer OrderPayer) *OrderPayHandler {
	return &OrderPayHandler{
		name:       "OrderPayHandler",
		orderPayer: orderPayer,
	}
}

func (h *OrderPayHandler) Handle(w http.ResponseWriter, r *http.Request) {
	req := &OrderPayRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handlers.GetErrorResponse(w, h.name, err, http.StatusInternalServerError)
		return
	}
	if err := req.Validate(); err != nil {
		handlers.GetErrorResponse(w, h.name, err, http.StatusBadRequest)
		return
	}

	err := h.orderPayer.PayOrder(r.Context(), req.OrderID)
	if err != nil {
		handlers.GetErrorResponse(w, h.name, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
