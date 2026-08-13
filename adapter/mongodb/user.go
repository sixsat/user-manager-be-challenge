package mongodb

import (
	"context"

	"github.com/sixsat/user-manager-be-challenge/domain"
	"github.com/sixsat/user-manager-be-challenge/port"
)

type userRepo struct {
}

func NewUserRepository() port.UserRepository {
	return &userRepo{}
}

func (s *userRepo) Create(ctx context.Context, req *domain.CreateUserReq) error {
	return nil
}

func (s *userRepo) GetByID(ctx context.Context, id string) (*domain.GetUserRes, error) {
	return nil, nil
}

func (s *userRepo) GetByEmail(ctx context.Context, email string) (*domain.GetUserRes, error) {
	return nil, nil
}

func (s *userRepo) List(ctx context.Context) ([]domain.GetUserRes, error) {
	return nil, nil
}

func (s *userRepo) Update(ctx context.Context, req *domain.UpdateUserReq) error {
	return nil
}

func (s *userRepo) Delete(ctx context.Context, id string) error {
	return nil
}

func (s *userRepo) Count(ctx context.Context) (int, error) {
	return 0, nil
}
