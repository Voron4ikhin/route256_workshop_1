package cart

import (
	"net/http"
	httpkit "route256/cart/internal/adapter/in/http/httpkit"
	portin "route256/cart/internal/port/in"
)

type ListHandler struct {
	name   string
	lister portin.CartLister
}

func NewListHandler(lister portin.CartLister) *ListHandler {
	return &ListHandler{
		name:   "cart.ListHandler",
		lister: lister,
	}
}

func (h *ListHandler) Handle(w http.ResponseWriter, r *http.Request) {
	req := &ListRequest{}
	if !httpkit.DecodeJSON(w, r, h.name, req) {
		return
	}
	if err := req.Validate(); err != nil {
		httpkit.WriteError(w, h.name, err, http.StatusBadRequest)
		return
	}

	items, total, err := h.lister.GetList(r.Context(), req.User)
	if err != nil {
		httpkit.WriteError(w, h.name, err, httpkit.ClassifyBusinessError(err))
		return
	}

	httpkit.WriteJSON(w, h.name, http.StatusOK, ListResponse{
		Item:       itemsFromDomain(items),
		TotalPrice: total,
	})
}
