package http

import (
	"flag"
	"log"
	"net/http"
	"route256/cart/internal/pkg/client/loms"
	"route256/cart/internal/pkg/client/product"
	hitem "route256/cart/internal/pkg/handlers/item"
	sitem "route256/cart/internal/pkg/services/item"
)

func newConfigFromFlags() config {
	const (
		defaultAddr        = ":8080"
		defaultLomsAddr    = "http://loms:8080"
		defaultProductAddr = "http://route256.pavl.uk:8080"
	)

	result := config{}
	flag.StringVar(&result.addr, "addr", defaultAddr, "server address, default: "+defaultAddr)
	flag.StringVar(&result.lomsAddr, "loms_addr", defaultLomsAddr, "loms server address, default: "+defaultLomsAddr)
	flag.StringVar(&result.productAddr, "product_addr", defaultProductAddr, "product server address, default: "+defaultProductAddr)
	flag.Parse()
	return result
}

type config struct {
	addr        string
	lomsAddr    string
	productAddr string
}

type App struct {
	config config
}

func NewApp() *App {
	return &App{
		config: newConfigFromFlags(),
	}
}

func (a App) Run() error {
	lomsClient, err := loms.New("loms client", a.config.lomsAddr)
	if err != nil {
		log.Fatal(err)
	}
	productClient, err := product.New("product client", a.config.productAddr)
	if err != nil {
		log.Fatal(err)
	}

	itemAddHandler := hitem.NewAddHandler(sitem.NewAddService(lomsClient, productClient))
	http.HandleFunc("/item/add", itemAddHandler.Handle)
	log.Fatal(http.ListenAndServe(a.config.addr, nil))
	return nil
}
