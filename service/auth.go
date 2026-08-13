package service

import (
	"context"

	"github.com/sixsat/user-manager-be-challenge/domain"
	"github.com/sixsat/user-manager-be-challenge/port"
)

type authSvc struct {
}

func NewAuthService() port.AuthService {
	return &authSvc{}
}

func (s *authSvc) Register(ctx context.Context, req *domain.RegisterUserReq) error {
	return nil
}

func (s *authSvc) Login(ctx context.Context, req *domain.LoginUserReq) (*domain.LoginUserRes, error) {
	return nil, nil
}
