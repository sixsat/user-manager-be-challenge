package httphandler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/sixsat/user-manager-be-challenge/domain"
)

func (h *handler) createUser(c *echo.Context) error {
	var req CreateUserReq
	if err := c.Bind(&req); err != nil {
		slog.Error("[handler] error binding request", slog.String("error", err.Error()))
		return ResBadRequest(c)
	}

	if err := h.validate.Struct(&req); err != nil {
		slog.Error("[handler] error validating request", slog.String("error", err.Error()))
		return ResBadRequest(c)
	}

	if err := h.userSvc.Create(c.Request().Context(), req.toDomain()); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, Res[any]{
		Code: CodeOK,
		Desc: DescOK,
	})
}

func (h *handler) listUsers(c *echo.Context) error {
	users, err := h.userSvc.List(c.Request().Context())
	if err != nil {
		return err
	}

	res := make([]GetUserRes, 0, len(users))
	for _, u := range users {
		res = append(res, GetUserRes{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
		})
	}

	return c.JSON(http.StatusOK, Res[[]GetUserRes]{
		Code: CodeOK,
		Desc: DescOK,
		Data: &res,
	})
}

func (h *handler) getUserByID(c *echo.Context) error {
	var req GetUserReq
	if err := c.Bind(&req); err != nil {
		slog.Error("[handler] error binding request", slog.String("error", err.Error()))
		return ResBadRequest(c)
	}

	user, err := h.userSvc.GetByID(c.Request().Context(), req.ID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			ResNotFound(c)
		}
		return err
	}

	return c.JSON(http.StatusOK, Res[GetUserRes]{
		Code: CodeOK,
		Desc: DescOK,
		Data: &GetUserRes{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	})
}

func (h *handler) updateUser(c *echo.Context) error {
	var req UpdateUserReq
	if err := c.Bind(&req); err != nil {
		slog.Error("[handler] error binding request", slog.String("error", err.Error()))
		return ResBadRequest(c)
	}

	if req.Name == nil && req.Email == nil {
		slog.Error("[handler] invalid request", slog.String("error", "name or email is required"))
		return ResBadRequest(c)
	}

	if err := h.validate.Struct(&req); err != nil {
		slog.Error("[handler] error validating request", slog.String("error", err.Error()))
		return ResBadRequest(c)
	}

	if err := h.userSvc.Update(c.Request().Context(), req.toDomain()); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			ResNotFound(c)
		}
		return err
	}

	return c.JSON(http.StatusOK, Res[any]{
		Code: CodeOK,
		Desc: DescOK,
	})
}

func (h *handler) deleteUser(c *echo.Context) error {
	var req DeleteUserReq
	if err := c.Bind(&req); err != nil {
		slog.Error("[handler] error binding request", slog.String("error", err.Error()))
		return ResBadRequest(c)
	}

	if err := h.userSvc.Delete(c.Request().Context(), req.ID); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			ResNotFound(c)
		}
		return err
	}

	return c.JSON(http.StatusOK, Res[any]{
		Code: CodeOK,
		Desc: DescOK,
	})
}
