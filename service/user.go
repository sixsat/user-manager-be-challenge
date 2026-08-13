package service

import (
	"context"

	"github.com/sixsat/user-manager-be-challenge/domain"
	"github.com/sixsat/user-manager-be-challenge/port"
)

type userSvc struct {
	userRepo port.UserRepository
}

func NewUserService() port.UserService {
	return &userSvc{}
}

func (s *userSvc) Create(ctx context.Context, req *domain.CreateUserReq) error {
	return nil
}

func (s *userSvc) GetByID(ctx context.Context, id string) (*domain.GetUserRes, error) {
	return nil, nil
}

func (s *userSvc) GetByEmail(ctx context.Context, email string) (*domain.GetUserRes, error) {
	return nil, nil
}

func (s *userSvc) List(ctx context.Context) ([]domain.GetUserRes, error) {
	return nil, nil
}

func (s *userSvc) Update(ctx context.Context, req *domain.UpdateUserReq) error {
	return nil
}

func (s *userSvc) Delete(ctx context.Context, id string) error {
	return nil
}

func (s *userSvc) Count(ctx context.Context) (int, error) {
	return 0, nil
}
