package httphandler

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

func HandleRoutes(g *echo.Group, jwtSignKey string) {
	auth := g.Group("/auth")
	{
		auth.POST("/register", registerUser)
		auth.POST("/login", loginUser)
	}

	users := g.Group("/users", echojwt.WithConfig(echojwt.Config{
		ErrorHandler: func(c *echo.Context, err error) error {
			return c.JSON(401, "TODO")
		},
		SigningKey:    jwtSignKey,
		SigningMethod: jwt.SigningMethodHS256.Name,
		TokenLookup:   "header:Authorization:Bearer ",
	}))
	{
		users.POST("/", createUser)
		users.GET("/", listUsers)
		users.GET("/:id", getUserByID)
		users.PATCH("/:id", updateUser)
		users.DELETE("/:id", deleteUser)
	}
}
