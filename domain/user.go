package domain

import (
	"time"
)

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

type CreateUserReq struct {
	Name     string
	Email    string
	Password string
}

type GetUserRes struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}

type UpdateUserReq struct {
	ID    string
	Name  *string
	Email *string
}
