package httphandler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/sixsat/user-manager-be-challenge/domain"
)

func (h *handler) registerUser(c *echo.Context) error {
	var req RegisterUserReq
	err := c.Bind(&req)
	if err != nil {
		slog.Error("[handler] error binding request", slog.String("error", err.Error()))
		return ResBadRequest(c)
	}

	err = h.validate.Struct(&req)
	if err != nil {
		slog.Error("[handler] error validating request", slog.String("error", err.Error()))
		return ResBadRequest(c)
	}

	err = h.authSvc.Register(c.Request().Context(), req.toDomain())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, Res[any]{
		Code: CodeOK,
		Desc: DescOK,
	})
}

func (h *handler) loginUser(c *echo.Context) error {
	var req LoginUserReq
	if err := c.Bind(&req); err != nil {
		slog.Error("[handler] error binding request", slog.String("error", err.Error()))
		return ResBadRequest(c)
	}

	if err := h.validate.Struct(&req); err != nil {
		slog.Error("[handler] error validating request", slog.String("error", err.Error()))
		return ResBadRequest(c)
	}

	res, err := h.authSvc.Login(c.Request().Context(), req.toDomain())
	if err != nil {
		if errors.Is(err, domain.BizErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, Res[any]{
				Code: domain.BizErrInvalidCredentials.Code,
				Desc: domain.BizErrInvalidCredentials.Desc,
			})
		}
		return err
	}

	return c.JSON(http.StatusOK, Res[LoginUserRes]{
		Code: CodeOK,
		Desc: DescOK,
		Data: &LoginUserRes{AccessToken: res.AccessToken},
	})
}
