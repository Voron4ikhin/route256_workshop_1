package http

import (
	"fmt"
	"log"
	"net/http"
	"route256/loms/internal/pkg/config"
	"route256/loms/internal/pkg/repository/inmemory"
)

type dependencies struct {
	ordersStorage *inmemory.OrdersStorage
	stocksStorage *inmemory.StocksStorage
}

func newDependencies(cfg config.Config) (dependencies, error) {
	ordersStorage, err := inmemory.NewOrdersStorage()
	if err != nil {
		return dependencies{}, fmt.Errorf("init orders storage: %w", err)
	}

	stocksStorage, err := inmemory.NewStocksStorage()
	if err != nil {
		return dependencies{}, fmt.Errorf("init stocks storage: %w", err)
	}

	return dependencies{
		ordersStorage: ordersStorage,
		stocksStorage: stocksStorage,
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
