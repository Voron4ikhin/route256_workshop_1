package http

import (
	"net/http"
	"route256/loms/internal/adapter/in/http/httpkit"
	"route256/loms/internal/adapter/in/http/order"
	"route256/loms/internal/adapter/in/http/stock"
	portin "route256/loms/internal/port/in"
)

type Dependencies struct {
	OrderCreator   portin.OrderCreator
	OrderPayer     portin.OrderPayer
	OrderCanceler  portin.OrderCanceler
	OrderInformant portin.OrderInformant
	StockInformant portin.StockInformant
}

func post(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			httpkit.WriteError(w, "router", httpkit.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	}
}

func NewRouter(deps Dependencies) *http.ServeMux {
	orderCreateHandler := order.NewCreateHandler(deps.OrderCreator)
	orderInfoHandler := order.NewInfoHandler(deps.OrderInformant)
	orderPayHandler := order.NewPayHandler(deps.OrderPayer)
	orderCancelHandler := order.NewCancelHandler(deps.OrderCanceler)
	stockInfoHandler := stock.NewInfoHandler(deps.StockInformant)

	mux := http.NewServeMux()
	mux.HandleFunc("/order/create", post(orderCreateHandler.Handle))
	mux.HandleFunc("/order/info", post(orderInfoHandler.Handle))
	mux.HandleFunc("/order/pay", post(orderPayHandler.Handle))
	mux.HandleFunc("/order/cancel", post(orderCancelHandler.Handle))
	mux.HandleFunc("/stock/info", post(stockInfoHandler.Handle))

	return mux
}
