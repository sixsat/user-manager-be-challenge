package httphandler

import "github.com/labstack/echo/v5"

func registerUser(c *echo.Context) error {
	return c.JSON(201, "TODO")
}

func loginUser(c *echo.Context) error {
	return c.JSON(200, "TODO")
}
