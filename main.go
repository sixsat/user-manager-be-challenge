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
	"github.com/sixsat/user-manager-be-challenge/handler/grpchandler"
	"github.com/sixsat/user-manager-be-challenge/handler/httphandler"
	"github.com/sixsat/user-manager-be-challenge/infra"
	userproto "github.com/sixsat/user-manager-be-challenge/proto"
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
	validate := validator.New(validator.WithRequiredStructEnabled())

	e := infra.NewHTTPServer()
	httphandler.
		New(
			cfg.JWT.SignKey,
			validate,
			authSvc,
			userSvc,
		).
		RegisterRoutes(e.Group("/api"))
	go infra.StartHTTPServer(ctx, cfg.HTTPServer, e)

	grpcServer := infra.NewGRPCServer(cfg.JWT.SignKey)
	userproto.RegisterUserServiceServer(grpcServer, grpchandler.New(validate, userSvc))
	go infra.StartGRPCServer(ctx, cfg.GRPCServer, grpcServer)

	infra.StartBackgroundJob(ctx, userSvc)
}
