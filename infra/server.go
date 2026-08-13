package infra

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/sixsat/user-manager-be-challenge/config"
	"github.com/sixsat/user-manager-be-challenge/handler/httphandler"
)

func StartHTTPServer(ctx context.Context, cfg config.HTTPServer) {
	e := echo.New()
	e.Use(
		middleware.RequestLogger(),
		middleware.Recover(),
	)

	httphandler.HandleRoutes(e.Group("/api"), cfg.JWTSignKey)

	sc := echo.StartConfig{
		Address:         ":" + cfg.Port,
		HideBanner:      true,
		GracefulTimeout: 10 * time.Second,
		OnShutdownError: func(err error) {
			slog.Error("error shutting down http server", slog.String("error", err.Error()))
		},
	}
	if err := sc.Start(ctx, e); err != nil {
		slog.Error("error starting http server", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("http server stopped")
}

func StartGRPCServer(ctx context.Context, cfg config.GRPCServer) {
	// TODO: impl
}
