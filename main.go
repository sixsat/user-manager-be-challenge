package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sixsat/user-manager-be-challenge/config"
	"github.com/sixsat/user-manager-be-challenge/infra"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("error loading config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	mongoClient, err := infra.ConnectMongoDB()
	if err != nil {
		slog.Error("error connecting mongodb", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := mongoClient.Disconnect(ctx)
		if err != nil {
			slog.Error("error disconnecting mongodb", slog.String("error", err.Error()))
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	infra.StartHTTPServer(ctx, cfg.HTTPServer)

	// TODO: start grpc server
}
