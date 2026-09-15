package cart

import (
	"net/http"
	httpkit "route256/cart/internal/adapter/in/http/httpkit"
	portin "route256/cart/internal/port/in"
)

type ClearHandler struct {
	name    string
	clearer portin.CartClearer
}

func NewClearHandler(clearer portin.CartClearer) *ClearHandler {
	return &ClearHandler{
		name:    "cart.ClearHandler",
		clearer: clearer,
	}
}

func (h *ClearHandler) Handle(w http.ResponseWriter, r *http.Request) {
	req := &ClearRequest{}
	if !httpkit.DecodeJSON(w, r, h.name, req) {
		return
	}
	if err := req.Validate(); err != nil {
		httpkit.WriteError(w, h.name, err, http.StatusBadRequest)
		return
	}

	if err := h.clearer.Clear(r.Context(), req.User); err != nil {
		httpkit.WriteError(w, h.name, err, httpkit.ClassifyBusinessError(err))
		return
	}

	httpkit.WriteJSON(w, h.name, http.StatusOK, nil)
}
