package http

import (
	"net/http"
	"route256/cart/internal/adapter/in/http/cart"
	"route256/cart/internal/adapter/in/http/httpkit"
	"route256/cart/internal/adapter/in/http/item"
	portin "route256/cart/internal/port/in"
)

type Dependencies struct {
	ItemAdder      portin.ItemAdder
	ItemDeleter    portin.ItemDeleter
	CartLister     portin.CartLister
	CartClearer    portin.CartClearer
	CartCheckouter portin.CartCheckouter
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
	itemAddHandler := item.NewAddHandler(deps.ItemAdder)
	itemDeleteHandler := item.NewDeleteHandler(deps.ItemDeleter)
	listHandler := cart.NewListHandler(deps.CartLister)
	clearHandler := cart.NewClearHandler(deps.CartClearer)
	checkoutHandler := cart.NewCheckoutHandler(deps.CartCheckouter)

	mux := http.NewServeMux()
	mux.HandleFunc("/cart/item/add", post(itemAddHandler.Handle))
	mux.HandleFunc("/cart/item/delete", post(itemDeleteHandler.Handle))
	mux.HandleFunc("/cart/list", post(listHandler.Handle))
	mux.HandleFunc("/cart/clear", post(clearHandler.Handle))
	mux.HandleFunc("/cart/checkout", post(checkoutHandler.Handle))

	return mux
}
