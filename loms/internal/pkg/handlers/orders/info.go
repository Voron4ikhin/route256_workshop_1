package orders

import (
	"context"
	"encoding/json"
	"net/http"
	"route256/loms/internal/pkg/handlers"
)

type OrderInfoRequest struct {
	OrderID int64 `json:"order_id,omitempty"`
}

func (r OrderInfoRequest) Validate() error {
	if r.OrderID <= 0 {
		return handlers.ErrIncorrectOrderID
	}
	return nil
}

type OrderInfoResponse struct {
	Status string      `json:"status,omitempty"`
	User   int64       `json:"user,omitempty"`
	Items  []OrderItem `json:"items,omitempty"`
}

type OrderInfoHandler struct {
	name           string
	orderInformant OrderInformant
}

type OrderInformant interface {
	GetOrderInfo(ctx context.Context, orderID int64) (*OrderInfoResponse, error)
}

func NewOrderInfoHandler(orderInformant OrderInformant) *OrderInfoHandler {
	return &OrderInfoHandler{
		name:           "orderInfoHandler",
		orderInformant: orderInformant,
	}
}

func (h *OrderInfoHandler) Handle(w http.ResponseWriter, r *http.Request) {
	req := &OrderInfoRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handlers.GetErrorResponse(w, h.name, err, http.StatusInternalServerError)
		return
	}
	if err := req.Validate(); err != nil {
		handlers.GetErrorResponse(w, h.name, err, http.StatusBadRequest)
		return
	}

	order, err := h.orderInformant.GetOrderInfo(r.Context(), req.OrderID)
	if err != nil {
		handlers.GetErrorResponse(w, h.name, err, http.StatusPreconditionFailed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		handlers.GetErrorResponse(w, h.name, err, http.StatusInternalServerError)
		return
	}
}
