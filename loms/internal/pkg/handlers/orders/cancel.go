package orders

import (
	"context"
	"encoding/json"
	"net/http"
	"route256/loms/internal/pkg/handlers"
)

type OrderCancelRequest struct {
	OrderID int64 `json:"order_id,omitempty"`
}

func (r OrderCancelRequest) Validate() error {
	if r.OrderID <= 0 {
		return handlers.ErrIncorrectOrderID
	}
	return nil
}

type OrderCanceler interface {
	CancelOrder(ctx context.Context, orderID int64) error
}

type OrderCancelHandler struct {
	name          string
	orderCanceler OrderCanceler
}

func NewOrderCancelHandler(orderCanceler OrderCanceler) *OrderCancelHandler {
	return &OrderCancelHandler{
		name:          "OrderCancelHandler",
		orderCanceler: orderCanceler,
	}
}

func (h *OrderCancelHandler) Handle(w http.ResponseWriter, r *http.Request) {
	req := &OrderCancelRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handlers.GetErrorResponse(w, h.name, err, http.StatusInternalServerError)
		return
	}
	if err := req.Validate(); err != nil {
		handlers.GetErrorResponse(w, h.name, err, http.StatusBadRequest)
		return
	}

	err := h.orderCanceler.CancelOrder(r.Context(), req.OrderID)
	if err != nil {
		handlers.GetErrorResponse(w, h.name, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
