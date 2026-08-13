package port

import (
	"context"

	"github.com/sixsat/user-manager-be-challenge/domain"
)

type AuthService interface {
	Register(ctx context.Context, req *domain.RegisterUserReq) error
	Login(ctx context.Context, req *domain.LoginUserReq) (*domain.LoginUserRes, error)
}

type UserService interface {
	Create(ctx context.Context, req *domain.CreateUserReq) error
	GetByID(ctx context.Context, id string) (*domain.GetUserRes, error)
	GetByEmail(ctx context.Context, email string) (*domain.GetUserRes, error)
	List(ctx context.Context) ([]domain.GetUserRes, error)
	Update(ctx context.Context, req *domain.UpdateUserReq) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int, error)
}
