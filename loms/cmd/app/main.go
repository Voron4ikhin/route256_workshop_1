package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	httpadapter "route256/loms/internal/adapter/in/http"
	"route256/loms/internal/adapter/in/scheduler"
	"route256/loms/internal/adapter/out/repository/inmemory"
	"route256/loms/internal/config"
	"route256/loms/internal/usecase/order"
	"route256/loms/internal/usecase/stock"
	"syscall"
	"time"
)

const (
	shutdownTimeout    = 10 * time.Second
	unpaidOrderTTL     = 10 * time.Minute
	autoCancelInterval = 30 * time.Second
)

func main() {
	if err := run(); err != nil {
		slog.Error("loms exited with error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := config.NewFromFlags()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	orderRepo := inmemory.NewOrderRepository()
	stockRepo := inmemory.NewStockRepository()

	cancelService := order.NewCancelService(orderRepo, stockRepo)
	autoCancelService := order.NewAutoCancelService(orderRepo, cancelService, unpaidOrderTTL)

	deps := httpadapter.Dependencies{
		OrderCreator:   order.NewCreateService(orderRepo, stockRepo),
		OrderPayer:     order.NewPayService(orderRepo, stockRepo),
		OrderCanceler:  cancelService,
		OrderInformant: order.NewInfoService(orderRepo),
		StockInformant: stock.NewQueryService(stockRepo),
	}

	server := httpadapter.NewServer(cfg.Addr, deps, logger)
	autoCancelScheduler := scheduler.NewAutoCancelScheduler(autoCancelService, autoCancelInterval, logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go autoCancelScheduler.Run(ctx)

	return server.Run(ctx, shutdownTimeout)
}
