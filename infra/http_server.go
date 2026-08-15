package infra

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/sixsat/user-manager-be-challenge/config"
	"github.com/sixsat/user-manager-be-challenge/domain"
	"github.com/sixsat/user-manager-be-challenge/handler/httphandler"
)

func NewHTTPServer() *echo.Echo {
	e := echo.New()
	e.Use(
		middleware.RequestLogger(),
		middleware.Recover(),
	)
	e.HTTPErrorHandler = httpErrorHandler
	return e
}

func StartHTTPServer(ctx context.Context, cfg config.HTTPServer, e *echo.Echo) {
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

func httpErrorHandler(c *echo.Context, err error) {
	if resp, err := echo.UnwrapResponse(c.Response()); err == nil {
		if resp.Committed {
			return
		}
	}

	if err, ok := errors.AsType[domain.BizErr](err); ok {
		_ = c.JSON(http.StatusConflict, httphandler.Res[any]{
			Code: err.Code,
			Desc: err.Desc,
		})
		return
	}

	_ = c.JSON(http.StatusInternalServerError, httphandler.Res[any]{
		Code: httphandler.CodeInternalErr,
		Desc: httphandler.DescInternalErr,
	})
}
