package domain

import (
	"errors"
	"fmt"
)

type BizErr struct {
	Code string
	Desc string
}

func (e BizErr) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Desc)
}

var (
	ErrUserAlreadyExists  = BizErr{Code: "0002", Desc: "user already exists"}
	ErrInvalidCredentials = BizErr{Code: "0003", Desc: "invalid email or password"}

	ErrDuplicateUser = errors.New("duplicate user")
	ErrUserNotFound  = errors.New("user not found")
)
