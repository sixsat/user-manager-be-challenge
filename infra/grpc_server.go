package infra

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/sixsat/user-manager-be-challenge/config"
	"google.golang.org/grpc"
)

func NewGRPCServer(jwtSignKey string) *grpc.Server {
	return grpc.NewServer(grpc.ChainUnaryInterceptor(
		loggerInterceptor,
		jwtAuthInterceptor(jwtSignKey),
	))
}

func StartGRPCServer(ctx context.Context, cfg config.GRPCServer, server *grpc.Server) {
	listener, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		slog.Error("error listening for grpc server", slog.String("error", err.Error()))
		os.Exit(1)
	}

	go stopGRPCServer(ctx, server)

	if err := server.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		slog.Error("error starting grpc server", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("grpc server stopped")
}

func stopGRPCServer(ctx context.Context, server *grpc.Server) {
	<-ctx.Done()

	stopped := make(chan struct{})
	go func() {
		server.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
	case <-time.After(10 * time.Second):
		server.Stop()
	}
}
