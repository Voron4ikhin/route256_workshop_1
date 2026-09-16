package grpc

import (
	"context"
	"log"
	"log/slog"
	"net"
	"time"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	grpcServer *grpc.Server
	listener   net.Listener
	logger     *slog.Logger
}

func NewGRPCServer(addr string, deps Dependencies, logger *slog.Logger) (*GRPCServer, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	Register(srv, deps)

	return &GRPCServer{grpcServer: srv, listener: lis, logger: logger}, nil
}

func (s *GRPCServer) Run(ctx context.Context, shutdownTimeout time.Duration) error {
	errCh := make(chan error, 1)
	go func() {
		s.logger.Info("grpc server listening", "addr", s.listener.Addr().String())
		errCh <- s.grpcServer.Serve(s.listener)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	s.logger.Info("grpc server shutting down")
	stopped := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
	case <-time.After(shutdownTimeout):
		s.grpcServer.Stop()
	}

	return <-errCh
}
