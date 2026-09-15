package order

import (
	"net/http"
	httpadapter "route256/loms/internal/adapter/in/http/httpkit"
	portin "route256/loms/internal/port/in"
)

type CancelHandler struct {
	name     string
	canceler portin.OrderCanceler
}

func NewCancelHandler(canceler portin.OrderCanceler) *CancelHandler {
	return &CancelHandler{
		name:     "order.CancelHandler",
		canceler: canceler,
	}
}

func (h *CancelHandler) Handle(w http.ResponseWriter, r *http.Request) {
	req := &CancelRequest{}
	if !httpadapter.DecodeJSON(w, r, h.name, req) {
		return
	}
	if err := req.Validate(); err != nil {
		httpadapter.WriteError(w, h.name, err, http.StatusBadRequest)
		return
	}

	if err := h.canceler.CancelOrder(r.Context(), req.OrderID); err != nil {
		httpadapter.WriteError(w, h.name, err, httpadapter.ClassifyBusinessError(err))
		return
	}

	httpadapter.WriteJSON(w, h.name, http.StatusOK, nil)
}
