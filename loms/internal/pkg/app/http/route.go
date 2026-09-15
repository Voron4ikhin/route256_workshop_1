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
	orderInfoHandler := horders.NewOrderInfoHandler(sorders.NewOrderInfoService(deps.ordersStorage))
	orderPayHandler := horders.NewOrderPayHandler(sorders.NewOrderPayService(deps.ordersStorage, deps.stocksStorage, deps.ordersStorage))
	orderCancelHandler := horders.NewOrderCancelHandler(sorders.NewOrderCancelService(deps.ordersStorage, deps.stocksStorage, deps.ordersStorage))
	stocksHandler := handlers.NewStocksHandler(services.NewStocksService(deps.stocksStorage))

	mux := http.NewServeMux()
	mux.HandleFunc("/order/create", orderCreateHandler.Handle)
	mux.HandleFunc("/order/info", orderInfoHandler.Handle)
	mux.HandleFunc("/order/pay", orderPayHandler.Handle)
	mux.HandleFunc("/order/cancel", orderCancelHandler.Handle)
	mux.HandleFunc("/stock/info", stocksHandler.Handle)

	return mux
}
