package scheduler

import (
	"context"
	"log/slog"
	portin "route256/loms/internal/port/in"
	"time"
)

type AutoCancelScheduler struct {
	autoCanceler portin.OrderAutoCanceler
	interval     time.Duration
	logger       *slog.Logger
}

func NewAutoCancelScheduler(autoCanceler portin.OrderAutoCanceler, interval time.Duration, logger *slog.Logger) *AutoCancelScheduler {
	return &AutoCancelScheduler{
		autoCanceler: autoCanceler,
		interval:     interval,
		logger:       logger,
	}
}

func (s *AutoCancelScheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			cancelled, err := s.autoCanceler.SweepUnpaidOrders(ctx)
			if err != nil {
				s.logger.Error("auto-cancel sweep finished with errors", "error", err, "cancelled", cancelled)
				continue
			}
			if cancelled > 0 {
				s.logger.Info("auto-cancel sweep cancelled unpaid orders", "count", cancelled)
			}
		}
	}
}
