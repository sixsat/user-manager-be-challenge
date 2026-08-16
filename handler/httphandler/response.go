package httphandler

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

const (
	CodeOK          = "0000"
	DescOK          = "success"
	CodeBadReq      = "0001"
	DescBadReq      = "bad request"
	CodeInternalErr = "9999"
	DescInternalErr = "internal server error"
)

type Res[T any] struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
	Data *T     `json:"data,omitempty"`
}

func ResBadRequest(c *echo.Context) error {
	return c.JSON(http.StatusBadRequest, Res[any]{
		Code: CodeBadReq,
		Desc: DescBadReq,
	})
}

func ResNotFound(c *echo.Context) error {
	return c.JSON(http.StatusNotFound, Res[any]{
		Code: CodeBadReq,
		Desc: DescBadReq,
	})
}
