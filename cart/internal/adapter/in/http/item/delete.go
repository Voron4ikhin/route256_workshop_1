package item

import (
	"net/http"
	httpkit "route256/cart/internal/adapter/in/http/httpkit"
	portin "route256/cart/internal/port/in"
)

type DeleteHandler struct {
	name    string
	deleter portin.ItemDeleter
}

func NewDeleteHandler(deleter portin.ItemDeleter) *DeleteHandler {
	return &DeleteHandler{
		name:    "item.DeleteHandler",
		deleter: deleter,
	}
}

func (h *DeleteHandler) Handle(w http.ResponseWriter, r *http.Request) {
	req := &DeleteRequest{}
	if !httpkit.DecodeJSON(w, r, h.name, req) {
		return
	}
	if err := req.Validate(); err != nil {
		httpkit.WriteError(w, h.name, err, http.StatusBadRequest)
		return
	}

	if err := h.deleter.Delete(r.Context(), req.User, req.SKU); err != nil {
		httpkit.WriteError(w, h.name, err, httpkit.ClassifyBusinessError(err))
		return
	}

	httpkit.WriteJSON(w, h.name, http.StatusOK, nil)
}
