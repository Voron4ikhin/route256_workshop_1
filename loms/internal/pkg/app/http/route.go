package http

import (
	"net/http"
	"route256/loms/internal/pkg/handlers"
	horders "route256/loms/internal/pkg/handlers/orders"
	sorders "route256/loms/internal/pkg/services/orders"

	"route256/loms/internal/pkg/services"
)

func newRouter(deps dependencies) *http.ServeMux {
	orderCreateHandler := horders.NewOrderCreateHandler(sorders.NewOrderCreateService(deps.ordersStorage, deps.ordersStorage, deps.stocksStorage))
	stocksHandler := handlers.NewStocksHandler(services.NewStocksService(deps.stocksStorage))

	mux := http.NewServeMux()
	mux.HandleFunc("/order/create", orderCreateHandler.Handle)
	mux.HandleFunc("/stock/info", stocksHandler.Handle)

	return mux
}
