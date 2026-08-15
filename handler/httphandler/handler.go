package httphandler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/sixsat/user-manager-be-challenge/port"
)

type handler struct {
	jwtSignKey string
	validate   *validator.Validate
	authSvc    port.AuthService
}

func New(jwtSignKey string, validate *validator.Validate, authSvc port.AuthService) *handler {
	return &handler{
		jwtSignKey: jwtSignKey,
		validate:   validate,
		authSvc:    authSvc,
	}
}

func (h *handler) RegisterRoutes(g *echo.Group) {
	auth := g.Group("/auth")
	{
		auth.POST("/register", h.registerUser)
		auth.POST("/login", h.loginUser)
	}

	users := g.Group("/users", echojwt.WithConfig(echojwt.Config{
		ErrorHandler: func(c *echo.Context, err error) error {
			return c.JSON(http.StatusUnauthorized, Res[any]{
				Code: CodeBadReq,
				Desc: DescBadReq,
			})
		},
		SigningKey:    h.jwtSignKey,
		SigningMethod: jwt.SigningMethodHS256.Name,
		TokenLookup:   "header:Authorization:Bearer ",
	}))
	{
		users.POST("/", h.createUser)
		users.GET("/", h.listUsers)
		users.GET("/:id", h.getUserByID)
		users.PATCH("/:id", h.updateUser)
		users.DELETE("/:id", h.deleteUser)
	}
}
