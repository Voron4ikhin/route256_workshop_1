package order

import (
	"net/http"
	httpadapter "route256/loms/internal/adapter/in/http/httpkit"
	portin "route256/loms/internal/port/in"
)

type PayHandler struct {
	name  string
	payer portin.OrderPayer
}

func NewPayHandler(payer portin.OrderPayer) *PayHandler {
	return &PayHandler{
		name:  "order.PayHandler",
		payer: payer,
	}
}

func (h *PayHandler) Handle(w http.ResponseWriter, r *http.Request) {
	req := &PayRequest{}
	if !httpadapter.DecodeJSON(w, r, h.name, req) {
		return
	}
	if err := req.Validate(); err != nil {
		httpadapter.WriteError(w, h.name, err, http.StatusBadRequest)
		return
	}

	if err := h.payer.PayOrder(r.Context(), req.OrderID); err != nil {
		httpadapter.WriteError(w, h.name, err, httpadapter.ClassifyBusinessError(err))
		return
	}

	httpadapter.WriteJSON(w, h.name, http.StatusOK, nil)
}
