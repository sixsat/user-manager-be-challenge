package grpchandler

import "github.com/sixsat/user-manager-be-challenge/domain"

type createUserReq struct {
	Name     string `validate:"required"`
	Email    string `validate:"email"`
	Password string `validate:"required"`
}

func (r *createUserReq) toDomain() *domain.CreateUserReq {
	return &domain.CreateUserReq{
		Name:         r.Name,
		Email:        r.Email,
		PasswordHash: r.Password,
	}
}

type getUserReq struct {
	ID string `validate:"mongodb"`
}
