package order

import (
	"net/http"
	httpadapter "route256/loms/internal/adapter/in/http/httpkit"
	portin "route256/loms/internal/port/in"
)

type CreateHandler struct {
	name    string
	creator portin.OrderCreator
}

func NewCreateHandler(creator portin.OrderCreator) *CreateHandler {
	return &CreateHandler{
		name:    "order.CreateHandler",
		creator: creator,
	}
}

func (h *CreateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	req := &CreateRequest{}
	if !httpadapter.DecodeJSON(w, r, h.name, req) {
		return
	}
	if err := req.Validate(); err != nil {
		httpadapter.WriteError(w, h.name, err, http.StatusBadRequest)
		return
	}

	orderID, err := h.creator.CreateOrder(r.Context(), req.User, itemsToDomain(req.Items))
	if err != nil {
		httpadapter.WriteError(w, h.name, err, httpadapter.ClassifyBusinessError(err))
		return
	}

	httpadapter.WriteJSON(w, h.name, http.StatusOK, CreateResponse{OrderID: orderID})
}
