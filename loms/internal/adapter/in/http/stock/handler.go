package stock

import (
	"net/http"
	httpadapter "route256/loms/internal/adapter/in/http/httpkit"
	portin "route256/loms/internal/port/in"
)

type InfoHandler struct {
	name      string
	informant portin.StockInformant
}

func NewInfoHandler(informant portin.StockInformant) *InfoHandler {
	return &InfoHandler{
		name:      "stock.InfoHandler",
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

	count := h.informant.GetStocks(r.Context(), req.SKU)

	httpadapter.WriteJSON(w, h.name, http.StatusOK, InfoResponse{Count: count})
}
