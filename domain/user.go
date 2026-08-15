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
	Name         string
	Email        string
	PasswordHash string
}

type GetUserRes struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}

type GetUserByEmailRes struct {
	ID           string
	PasswordHash string
}

type UpdateUserReq struct {
	ID    string
	Name  *string
	Email *string
}
