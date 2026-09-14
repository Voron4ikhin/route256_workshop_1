package http

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

	"route256/cart/internal/pkg/client/loms"
	"route256/cart/internal/pkg/client/product"
	productmock "route256/cart/internal/pkg/client/product/mock"
	"route256/cart/internal/pkg/handlers"
	hitem "route256/cart/internal/pkg/handlers/item"
	"route256/cart/internal/pkg/repository/inmemory"
	"route256/cart/internal/pkg/services"
	sitem "route256/cart/internal/pkg/services/item"
)

type config struct {
	addr        string
	lomsAddr    string
	productAddr string
	productMock bool
}

// envOrDefault returns the value of the env var if set, otherwise def.
// Flags below use this as their default, so the precedence is: flag > env > hardcoded default.
func envOrDefault(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func envOrDefaultBool(key string, def bool) bool {
	if v, ok := os.LookupEnv(key); ok {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}

func newConfigFromFlags() config {
	var (
		defaultAddr        = envOrDefault("ADDR", ":8080")
		defaultLomsAddr    = envOrDefault("LOMS_ADDR", "http://loms:8080")
		defaultProductAddr = envOrDefault("PRODUCT_ADDR", "http://route256.pavl.uk:8080") // это не работает
		defaultProductMock = envOrDefaultBool("PRODUCT_MOCK", true)                       // TODO: вернуть false, когда product-сервис снова станет доступен
	)

	result := config{}
	flag.StringVar(&result.addr, "addr", defaultAddr, "server address, default: "+defaultAddr)
	flag.StringVar(&result.lomsAddr, "loms_addr", defaultLomsAddr, "loms server address, default: "+defaultLomsAddr)
	flag.StringVar(&result.productAddr, "product_addr", defaultProductAddr, "product server address, default: "+defaultProductAddr)
	flag.BoolVar(&result.productMock, "product_mock", defaultProductMock, "use in-memory mock instead of the real product client")
	flag.Parse()
	return result
}

// productProvider is satisfied by both product.Client and its mock stand-in.
type productProvider interface {
	GetProductInfo(ctx context.Context, sku uint32) (string, uint32, error)
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
	var productClient productProvider
	if a.config.productMock {
		productClient = productmock.New("product client (mock)")
		log.Println("product client: using mock, real product service is not used")
	} else {
		realProductClient, err := product.New("product client", a.config.productAddr)
		if err != nil {
			log.Fatal(err)
		}
		productClient = realProductClient
	}
	cartRepository, err := inmemory.NewRepository()
	if err != nil {
		log.Fatal(err)
	}

	itemAddHandler := hitem.NewAddHandler(sitem.NewAddService(lomsClient, productClient, cartRepository))
	itemDeleteHandler := hitem.NewDeleteHandler(sitem.NewDeleteService(cartRepository))
	listHandler := handlers.NewListHandler(services.NewListService(productClient, cartRepository))
	clearHandler := handlers.NewClearHandler(services.NewClearService(cartRepository))
	checkoutHandler := handlers.NewCheckoutHandler(services.NewCheckoutService(cartRepository, lomsClient))

	http.HandleFunc("/item/add", itemAddHandler.Handle)
	http.HandleFunc("/item/delete", itemDeleteHandler.Handle)
	http.HandleFunc("/list", listHandler.Handle)
	http.HandleFunc("/clear", clearHandler.Handle)
	http.HandleFunc("/checkout", checkoutHandler.Handle)

	return http.ListenAndServe(a.config.addr, nil)
}
