package httphandler

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func (h *handler) createUser(c *echo.Context) error {
	return c.JSON(http.StatusCreated, "TODO")
}

func (h *handler) listUsers(c *echo.Context) error {
	return c.JSON(http.StatusOK, "TODO")
}

func (h *handler) getUserByID(c *echo.Context) error {
	return c.JSON(http.StatusOK, "TODO")
}

func (h *handler) updateUser(c *echo.Context) error {
	return c.JSON(http.StatusOK, "TODO")
}

func (h *handler) deleteUser(c *echo.Context) error {
	return c.JSON(http.StatusNoContent, "TODO")
}
