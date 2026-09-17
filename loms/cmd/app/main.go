package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	grpcadapter "route256/loms/internal/adapter/in/grpc"
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

	createService := order.NewCreateService(orderRepo, stockRepo)
	payService := order.NewPayService(orderRepo, stockRepo)
	infoService := order.NewInfoService(orderRepo)
	queryService := stock.NewQueryService(stockRepo)

	httpDeps := httpadapter.Dependencies{
		OrderCreator:   createService,
		OrderPayer:     payService,
		OrderCanceler:  cancelService,
		OrderInformant: infoService,
		StockInformant: stock.NewQueryService(stockRepo),
	}

	grpcDeps := grpcadapter.Dependencies{
		OrderCreator:   createService,
		OrderPayer:     payService,
		OrderCanceler:  cancelService,
		OrderInformant: infoService,
		StockInformant: queryService,
	}

	httpServer := httpadapter.NewServer(cfg.Addr, httpDeps, logger)
	grpcServer, err := grpcadapter.NewGRPCServer(cfg.GRPCAddr, grpcDeps, logger)
	if err != nil {
		return err
	}

	autoCancelScheduler := scheduler.NewAutoCancelScheduler(autoCancelService, autoCancelInterval, logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go autoCancelScheduler.Run(ctx)

	errCh := make(chan error, 2)
	go func() { errCh <- httpServer.Run(ctx, shutdownTimeout) }()
	go func() { errCh <- grpcServer.Run(ctx, shutdownTimeout) }()

	if err := <-errCh; err != nil {
		return err
	}
	return <-errCh
}
