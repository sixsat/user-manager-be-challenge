package httphandler

import (
	"time"

	"github.com/sixsat/user-manager-be-challenge/domain"
)

type RegisterUserReq struct {
	Name     string `json:"name"     validate:"required"`
	Email    string `json:"email"    validate:"email"`
	Password string `json:"password" validate:"required"`
}

func (r *RegisterUserReq) toDomain() *domain.RegisterUserReq {
	return &domain.RegisterUserReq{
		Name:         r.Name,
		Email:        r.Email,
		PasswordHash: r.Password,
	}
}

type LoginUserReq struct {
	Email    string `json:"email"    validate:"email"`
	Password string `json:"password" validate:"required"`
}

func (r *LoginUserReq) toDomain() *domain.LoginUserReq {
	return &domain.LoginUserReq{
		Email:    r.Email,
		Password: r.Password,
	}
}

type LoginUserRes struct {
	AccessToken string `json:"access_token"`
}

type CreateUserReq struct {
	Name     string `json:"name"     validate:"required"`
	Email    string `json:"email"    validate:"email"`
	Password string `json:"password" validate:"required"`
}

func (r *CreateUserReq) toDomain() *domain.CreateUserReq {
	return &domain.CreateUserReq{
		Name:         r.Name,
		Email:        r.Email,
		PasswordHash: r.Password,
	}
}

type UpdateUserReq struct {
	ID    string  `json:"-"     param:"id" validate:"mongodb"`
	Name  *string `json:"name"  validate:"omitempty,min=1"`
	Email *string `json:"email" validate:"omitempty,email"`
}

func (r *UpdateUserReq) toDomain() *domain.UpdateUserReq {
	return &domain.UpdateUserReq{
		ID:    r.ID,
		Name:  r.Name,
		Email: r.Email,
	}
}

type GetUserReq struct {
	ID string `json:"-" param:"id" validate:"mongodb"`
}

type GetUserRes struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type DeleteUserReq struct {
	ID string `json:"-" param:"id" validate:"mongodb"`
}
