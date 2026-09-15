package item

import (
	"net/http"
	httpkit "route256/cart/internal/adapter/in/http/httpkit"
	portin "route256/cart/internal/port/in"
)

type AddHandler struct {
	name  string
	adder portin.ItemAdder
}

func NewAddHandler(adder portin.ItemAdder) *AddHandler {
	return &AddHandler{
		name:  "item.AddHandler",
		adder: adder,
	}
}

func (h *AddHandler) Handle(w http.ResponseWriter, r *http.Request) {
	req := &AddRequest{}
	if !httpkit.DecodeJSON(w, r, h.name, req) {
		return
	}
	if err := req.Validate(); err != nil {
		httpkit.WriteError(w, h.name, err, http.StatusBadRequest)
		return
	}

	if err := h.adder.Add(r.Context(), req.User, req.SKU, req.Count); err != nil {
		httpkit.WriteError(w, h.name, err, httpkit.ClassifyBusinessError(err))
		return
	}

	httpkit.WriteJSON(w, h.name, http.StatusOK, nil)
}
