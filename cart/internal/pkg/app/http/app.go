package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"route256/cart/internal/pkg/client/loms"
	"route256/cart/internal/pkg/client/product"
	productmock "route256/cart/internal/pkg/client/product/mock"
	"route256/cart/internal/pkg/config"
	"route256/cart/internal/pkg/repository/inmemory"
)

type productProvider interface {
	GetProductInfo(ctx context.Context, sku uint32) (string, uint32, error)
}

type dependencies struct {
	lomsClient     *loms.Client
	productClient  productProvider
	cartRepository *inmemory.InMemoryRepo
}

func newDependencies(cfg config.Config) (dependencies, error) {
	lomsClient, err := loms.New("loms client", cfg.LomsAddr)
	if err != nil {
		return dependencies{}, fmt.Errorf("init loms client: %w", err)
	}

	var productClient productProvider
	if cfg.ProductMock {
		productClient = productmock.New("product client (mock)")
		log.Println("product client: using mock, real product service is not used")
	} else {
		realProductClient, err := product.New("product client", cfg.ProductAddr)
		if err != nil {
			return dependencies{}, fmt.Errorf("init product client: %w", err)
		}
		productClient = realProductClient
	}

	cartRepository, err := inmemory.NewRepository()
	if err != nil {
		return dependencies{}, fmt.Errorf("init cart repository: %w", err)
	}

	return dependencies{
		lomsClient:     lomsClient,
		productClient:  productClient,
		cartRepository: cartRepository,
	}, nil
}

type App struct {
	config config.Config
}

func NewApp(cfg config.Config) *App {
	return &App{
		config: cfg,
	}
}

func (a App) Run() error {
	deps, err := newDependencies(a.config)
	if err != nil {
		log.Fatal(err)
	}

	mux := newRouter(deps)

	return http.ListenAndServe(a.config.Addr, mux)
}
