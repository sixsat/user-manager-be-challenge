package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"github.com/sixsat/user-manager-be-challenge/domain"
	"github.com/sixsat/user-manager-be-challenge/port"
	"golang.org/x/crypto/bcrypt"
)

type userSvc struct {
	userRepo port.UserRepository
}

func NewUserService(userRepo port.UserRepository) port.UserService {
	return &userSvc{
		userRepo: userRepo,
	}
}

func (s *userSvc) Create(ctx context.Context, req *domain.CreateUserReq) error {
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
		slog.Error("[service] error creating user", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (s *userSvc) GetByID(ctx context.Context, id string) (*domain.GetUserRes, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			slog.Error("[service] error getting user by id", slog.String("error", err.Error()))
		}
		return nil, err
	}
	return user, nil
}

func (s *userSvc) List(ctx context.Context) ([]domain.GetUserRes, error) {
	users, err := s.userRepo.List(ctx)
	if err != nil {
		slog.Error("[service] error listing users", slog.String("error", err.Error()))
		return nil, err
	}
	return users, nil
}

func (s *userSvc) Update(ctx context.Context, req *domain.UpdateUserReq) error {
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		req.Name = &name
	}
	if req.Email != nil {
		email := strings.ToLower(strings.TrimSpace(*req.Email))
		req.Email = &email
	}

	err := s.userRepo.Update(ctx, req)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateUser) {
			return domain.BizErrUserAlreadyExists
		}
		if !errors.Is(err, domain.ErrUserNotFound) {
			slog.Error("[service] error updating user", slog.String("error", err.Error()))
		}
		return err
	}
	return nil
}

func (s *userSvc) Delete(ctx context.Context, id string) error {
	err := s.userRepo.Delete(ctx, id)
	if err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			slog.Error("[service] error deleting user", slog.String("error", err.Error()))
		}
		return err
	}
	return nil
}

func (s *userSvc) Count(ctx context.Context) (int, error) {
	count, err := s.userRepo.Count(ctx)
	if err != nil {
		slog.Error("[service] error counting users", slog.String("error", err.Error()))
		return 0, err
	}
	return count, nil
}
