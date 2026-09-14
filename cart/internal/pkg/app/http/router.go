package http

import (
	"net/http"
	"route256/cart/internal/pkg/handlers"
	hitem "route256/cart/internal/pkg/handlers/item"
	"route256/cart/internal/pkg/services"
	sitem "route256/cart/internal/pkg/services/item"
)

func newRouter(deps dependencies) *http.ServeMux {
	itemAddHandler := hitem.NewAddHandler(sitem.NewAddService(deps.lomsClient, deps.productClient, deps.cartRepository))
	itemDeleteHandler := hitem.NewDeleteHandler(sitem.NewDeleteService(deps.cartRepository))
	listHandler := handlers.NewListHandler(services.NewListService(deps.productClient, deps.cartRepository))
	clearHandler := handlers.NewClearHandler(services.NewClearService(deps.cartRepository))
	checkoutHandler := handlers.NewCheckoutHandler(services.NewCheckoutService(deps.cartRepository, deps.lomsClient))

	mux := http.NewServeMux()
	mux.HandleFunc("/item/add", itemAddHandler.Handle)
	mux.HandleFunc("/item/delete", itemDeleteHandler.Handle)
	mux.HandleFunc("/list", listHandler.Handle)
	mux.HandleFunc("/clear", clearHandler.Handle)
	mux.HandleFunc("/checkout", checkoutHandler.Handle)

	return mux
}
