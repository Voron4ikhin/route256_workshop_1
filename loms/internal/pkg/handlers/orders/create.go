package orders

import (
	"context"
	"encoding/json"
	"net/http"
	"route256/loms/internal/pkg/handlers"
)

type OrderItem struct {
	SKU   uint32 `json:"sku,omitempty"`
	Count uint16 `json:"count,omitempty"`
}

type OrderCreateRequest struct {
	User  int64       `json:"user,omitempty"`
	Items []OrderItem `json:"items,omitempty"`
}

func (r OrderCreateRequest) Validate() error {
	if r.User <= 0 {
		return handlers.ErrIncorrectUser
	}
	if len(r.Items) == 0 {
		return handlers.ErrIncorrectQuantity
	}
	for _, item := range r.Items {
		if item.SKU == 0 {
			return handlers.ErrIncorrectSKU
		}
		if item.Count <= 0 {
			return handlers.ErrIncorrectProductCount
		}
	}
	return nil
}

type OrderCreateResponse struct {
	OrderId int64 `json:"order_id,omitempty"`
}

type OrderCreator interface {
	CreateOrder(ctx context.Context, user int64, items []OrderItem) (uint64, error)
}

type OrderCreateHandler struct {
	name         string
	orderCreator OrderCreator
}

func NewOrderCreateHandler(orderCreator OrderCreator) *OrderCreateHandler {
	return &OrderCreateHandler{
		name:         "orders create handler",
		orderCreator: orderCreator,
	}
}

func (h *OrderCreateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	req := &OrderCreateRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		handlers.GetErrorResponse(w, h.name, err, http.StatusInternalServerError)
		return
	}
	if err := req.Validate(); err != nil {
		handlers.GetErrorResponse(w, h.name, err, http.StatusBadRequest)
		return
	}

	order, err := h.orderCreator.CreateOrder(r.Context(), req.User, req.Items)
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
