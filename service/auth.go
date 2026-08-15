package service

import (
	"context"
	"log/slog"
	"strings"

	"github.com/sixsat/user-manager-be-challenge/domain"
	"github.com/sixsat/user-manager-be-challenge/port"
	"golang.org/x/crypto/bcrypt"
)

type authSvc struct {
	userRepo port.UserRepository
}

func NewAuthService(userRepo port.UserRepository) port.AuthService {
	return &authSvc{
		userRepo: userRepo,
	}
}

func (s *authSvc) Register(ctx context.Context, req *domain.RegisterUserReq) error {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(req.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("[service] error bcrypt hashing password", slog.String("error", err.Error()))
		return err
	}

	err = s.userRepo.Create(ctx, &domain.CreateUserReq{
		Name:         req.Name,
		Email:        strings.ToLower(req.Email),
		PasswordHash: string(passwordBytes),
	})
	if err != nil {
		slog.Error("[service] error creating user", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (s *authSvc) Login(ctx context.Context, req *domain.LoginUserReq) (*domain.LoginUserRes, error) {
	return nil, nil
}
