package domain

import (
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
	ErrUserAlreadyExists = BizErr{Code: "0002", Desc: "user already exists"}
)
