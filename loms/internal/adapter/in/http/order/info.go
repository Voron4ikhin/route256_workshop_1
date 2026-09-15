package order

import (
	"net/http"
	httpadapter "route256/loms/internal/adapter/in/http/httpkit"
	portin "route256/loms/internal/port/in"
)

type InfoHandler struct {
	name      string
	informant portin.OrderInformant
}

func NewInfoHandler(informant portin.OrderInformant) *InfoHandler {
	return &InfoHandler{
		name:      "order.InfoHandler",
		informant: informant,
	}
}

func (h *InfoHandler) Handle(w http.ResponseWriter, r *http.Request) {
	req := &InfoRequest{}
	if !httpadapter.DecodeJSON(w, r, h.name, req) {
		return
	}
	if err := req.Validate(); err != nil {
		httpadapter.WriteError(w, h.name, err, http.StatusBadRequest)
		return
	}

	order, err := h.informant.GetOrderInfo(r.Context(), req.OrderID)
	if err != nil {
		httpadapter.WriteError(w, h.name, err, httpadapter.ClassifyBusinessError(err))
		return
	}

	httpadapter.WriteJSON(w, h.name, http.StatusOK, infoResponseFromDomain(order))
}
