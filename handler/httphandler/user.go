package httphandler

import "github.com/labstack/echo/v5"

func createUser(c *echo.Context) error {
	return c.JSON(201, "TODO")
}

func listUsers(c *echo.Context) error {
	return c.JSON(200, "TODO")
}

func getUserByID(c *echo.Context) error {
	return c.JSON(200, "TODO")
}

func updateUser(c *echo.Context) error {
	return c.JSON(200, "TODO")
}

func deleteUser(c *echo.Context) error {
	return c.JSON(204, "TODO")
}
