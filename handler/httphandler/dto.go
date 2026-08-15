package httphandler

import "github.com/sixsat/user-manager-be-challenge/domain"

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
