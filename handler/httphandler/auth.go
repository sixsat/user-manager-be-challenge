package httphandler

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
)

func (h *handler) registerUser(c *echo.Context) error {
	var req RegisterUserReq
	err := c.Bind(&req)
	if err != nil {
		slog.Error("[handler] error binding request", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, Res[any]{
			Code: CodeBadReq,
			Desc: DescBadReq,
		})
	}

	err = h.validate.Struct(&req)
	if err != nil {
		slog.Error("[handler] error validating request", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, Res[any]{
			Code: CodeBadReq,
			Desc: DescBadReq,
		})
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
	return c.JSON(http.StatusOK, "TODO")
}
