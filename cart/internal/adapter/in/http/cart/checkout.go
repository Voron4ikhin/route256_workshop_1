package cart

import (
	"net/http"
	httpkit "route256/cart/internal/adapter/in/http/httpkit"
	portin "route256/cart/internal/port/in"
)

type CheckoutHandler struct {
	name       string
	checkouter portin.CartCheckouter
}

func NewCheckoutHandler(checkouter portin.CartCheckouter) *CheckoutHandler {
	return &CheckoutHandler{
		name:       "cart.CheckoutHandler",
		checkouter: checkouter,
	}
}

func (h *CheckoutHandler) Handle(w http.ResponseWriter, r *http.Request) {
	req := &CheckoutRequest{}
	if !httpkit.DecodeJSON(w, r, h.name, req) {
		return
	}
	if err := req.Validate(); err != nil {
		httpkit.WriteError(w, h.name, err, http.StatusBadRequest)
		return
	}

	orderID, err := h.checkouter.Checkout(r.Context(), req.User)
	if err != nil {
		httpkit.WriteError(w, h.name, err, httpkit.ClassifyBusinessError(err))
		return
	}

	httpkit.WriteJSON(w, h.name, http.StatusOK, orderID)
}
