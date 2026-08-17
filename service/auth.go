package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sixsat/user-manager-be-challenge/domain"
	"github.com/sixsat/user-manager-be-challenge/port"
	"golang.org/x/crypto/bcrypt"
)

type authSvc struct {
	userRepo   port.UserRepository
	jwtSignKey string
	jwtExpiry  time.Duration
}

func NewAuthService(userRepo port.UserRepository, jwtSignKey string, jwtExpiry time.Duration) port.AuthService {
	return &authSvc{
		userRepo:   userRepo,
		jwtSignKey: jwtSignKey,
		jwtExpiry:  jwtExpiry,
	}
}

func (s *authSvc) Register(ctx context.Context, req *domain.RegisterUserReq) error {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(req.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("[service] error bcrypt hashing password", slog.String("error", err.Error()))
		return err
	}

	err = s.userRepo.Create(ctx, &domain.CreateUserReq{
		Name:         strings.TrimSpace(req.Name),
		Email:        strings.ToLower(strings.TrimSpace(req.Email)),
		PasswordHash: string(passwordBytes),
	})
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateUser) {
			return domain.BizErrUserAlreadyExists
		}
		slog.Error("[service] error registering user", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (s *authSvc) Login(ctx context.Context, req *domain.LoginUserReq) (*domain.LoginUserRes, error) {
	user, err := s.userRepo.GetByEmail(ctx, strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.BizErrInvalidCredentials
		}
		slog.Error("[service] error getting user by email", slog.String("error", err.Error()))
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, domain.BizErrInvalidCredentials
		}
		slog.Error("[service] error comparing password hash", slog.String("error", err.Error()))
		return nil, err
	}

	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   user.ID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtExpiry)),
	})
	accessToken, err := token.SignedString([]byte(s.jwtSignKey))
	if err != nil {
		slog.Error("[service] error signing access token", slog.String("error", err.Error()))
		return nil, err
	}

	return &domain.LoginUserRes{AccessToken: accessToken}, nil
}
