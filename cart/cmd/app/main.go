package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	httpadapter "route256/cart/internal/adapter/in/http"
	"route256/cart/internal/adapter/out/client/loms"
	"route256/cart/internal/adapter/out/client/product"
	productmock "route256/cart/internal/adapter/out/client/product/mock"
	"route256/cart/internal/adapter/out/repository/inmemory"
	"route256/cart/internal/config"
	"route256/cart/internal/port/out"
	usecasecart "route256/cart/internal/usecase/cart"
	usecaseitem "route256/cart/internal/usecase/item"
	"syscall"
	"time"
)

const shutdownTimeout = 10 * time.Second

func main() {
	if err := run(); err != nil {
		slog.Error("cart exited with error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := config.NewFromFlags()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	lomsClient, err := loms.New("loms client", cfg.LomsAddr)
	if err != nil {
		return fmt.Errorf("init loms client: %w", err)
	}

	var productClient out.ProductClient
	if cfg.ProductMock {
		productClient = productmock.New("product client (mock)")
		logger.Info("product client: using mock, real product service is not used")
	} else {
		realProductClient, err := product.New("product client", cfg.ProductAddr)
		if err != nil {
			return fmt.Errorf("init product client: %w", err)
		}
		productClient = realProductClient
	}

	cartRepository := inmemory.NewCartRepository()

	deps := httpadapter.Dependencies{
		ItemAdder:      usecaseitem.NewAddService(productClient, lomsClient, cartRepository),
		ItemDeleter:    usecaseitem.NewDeleteService(cartRepository),
		CartLister:     usecasecart.NewListService(productClient, cartRepository),
		CartClearer:    usecasecart.NewClearService(cartRepository),
		CartCheckouter: usecasecart.NewCheckoutService(cartRepository, lomsClient),
	}

	server := httpadapter.NewServer(cfg.Addr, deps, logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	return server.Run(ctx, shutdownTimeout)
}
