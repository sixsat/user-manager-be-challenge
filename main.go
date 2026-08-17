package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sixsat/user-manager-be-challenge/adapter/mongodb"
	"github.com/sixsat/user-manager-be-challenge/config"
	"github.com/sixsat/user-manager-be-challenge/handler/httphandler"
	"github.com/sixsat/user-manager-be-challenge/infra"
	"github.com/sixsat/user-manager-be-challenge/service"
)

func main() {
	infra.InitLogger()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("error loading config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	mongoClient, err := infra.ConnectMongoDB(cfg.Mongo.URI)
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

	userRepo := mongodb.NewUserRepository(mongoClient)

	authSvc := service.NewAuthService(userRepo, cfg.JWT.SignKey, cfg.JWT.Expiry)
	userSvc := service.NewUserService(userRepo)

	e := infra.NewHTTPServer()
	httphandler.
		New(
			cfg.JWT.SignKey,
			validator.New(validator.WithRequiredStructEnabled()),
			authSvc,
			userSvc,
		).
		RegisterRoutes(e.Group("/api"))

	go infra.StartHTTPServer(ctx, cfg.HTTPServer, e)

	infra.StartBackgroundJob(ctx, userSvc)
	// TODO: start grpc server
}
