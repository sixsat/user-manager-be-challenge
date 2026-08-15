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
