package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/sixsat/user-manager-be-challenge/domain"
	"github.com/sixsat/user-manager-be-challenge/mock"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserServiceCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mock.NewMockUserRepository(gomock.NewController(t))
		svc := NewUserService(repo)
		repo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, req *domain.CreateUserReq) error {
				if req.Name != "John Doe" || req.Email != "john@example.com" {
					t.Fatalf("Create() request = %#v", req)
				}
				if err := bcrypt.CompareHashAndPassword([]byte(req.PasswordHash), []byte("password")); err != nil {
					t.Fatalf("stored password hash does not match: %v", err)
				}
				return nil
			},
		)

		err := svc.Create(context.Background(), &domain.CreateUserReq{
			Name:         "  John Doe ",
			Email:        " JOHN@Example.COM ",
			PasswordHash: "password",
		})

		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
	})

	tests := []struct {
		name    string
		repoErr error
		wantErr error
	}{
		{name: "fail duplicate user", repoErr: domain.ErrDuplicateUser, wantErr: domain.BizErrUserAlreadyExists},
		{name: "fail repository error", repoErr: errors.New("database unavailable")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mock.NewMockUserRepository(gomock.NewController(t))
			svc := NewUserService(repo)
			repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(tt.repoErr)

			err := svc.Create(context.Background(), &domain.CreateUserReq{PasswordHash: "password"})

			wantErr := tt.wantErr
			if wantErr == nil {
				wantErr = tt.repoErr
			}
			if !errors.Is(err, wantErr) {
				t.Fatalf("Create() error = %v, want %v", err, wantErr)
			}
		})
	}

	t.Run("fail password too long", func(t *testing.T) {
		repo := mock.NewMockUserRepository(gomock.NewController(t))
		svc := NewUserService(repo)

		err := svc.Create(context.Background(), &domain.CreateUserReq{PasswordHash: string(make([]byte, 73))})

		if !errors.Is(err, bcrypt.ErrPasswordTooLong) {
			t.Fatalf("Create() error = %v, want %v", err, bcrypt.ErrPasswordTooLong)
		}
	})
}

func TestUserServiceGetByID(t *testing.T) {
	want := &domain.GetUserRes{ID: "user-id", Name: "John"}
	repoErr := errors.New("database unavailable")
	tests := []struct {
		name     string
		repoUser *domain.GetUserRes
		repoErr  error
	}{
		{name: "success", repoUser: want},
		{name: "fail not found", repoErr: domain.ErrUserNotFound},
		{name: "fail repository error", repoErr: repoErr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mock.NewMockUserRepository(gomock.NewController(t))
			svc := NewUserService(repo)
			repo.EXPECT().GetByID(gomock.Any(), "user-id").Return(tt.repoUser, tt.repoErr)

			got, err := svc.GetByID(context.Background(), "user-id")

			if !errors.Is(err, tt.repoErr) {
				t.Fatalf("GetByID() error = %v, want %v", err, tt.repoErr)
			}
			if !reflect.DeepEqual(got, tt.repoUser) {
				t.Fatalf("GetByID() = %#v, want %#v", got, tt.repoUser)
			}
		})
	}
}

func TestUserServiceList(t *testing.T) {
	want := []domain.GetUserRes{{ID: "one"}, {ID: "two"}}
	repoErr := errors.New("database unavailable")
	tests := []struct {
		name      string
		repoUsers []domain.GetUserRes
		repoErr   error
	}{
		{name: "success", repoUsers: want},
		{name: "fail repository error", repoErr: repoErr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mock.NewMockUserRepository(gomock.NewController(t))
			svc := NewUserService(repo)
			repo.EXPECT().List(gomock.Any()).Return(tt.repoUsers, tt.repoErr)

			got, err := svc.List(context.Background())

			if !errors.Is(err, tt.repoErr) {
				t.Fatalf("List() error = %v, want %v", err, tt.repoErr)
			}
			if !reflect.DeepEqual(got, tt.repoUsers) {
				t.Fatalf("List() = %#v, want %#v", got, tt.repoUsers)
			}
		})
	}
}

func TestUserServiceUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mock.NewMockUserRepository(gomock.NewController(t))
		svc := NewUserService(repo)
		name, email := "  John Doe  ", "  JOHN@Example.COM "
		repo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, req *domain.UpdateUserReq) error {
				if req.ID != "user-id" || req.Name == nil || *req.Name != "John Doe" || req.Email == nil || *req.Email != "john@example.com" {
					t.Fatalf("Update() request = %#v", req)
				}
				return nil
			},
		)

		if err := svc.Update(context.Background(), &domain.UpdateUserReq{ID: "user-id", Name: &name, Email: &email}); err != nil {
			t.Fatalf("Update() error = %v", err)
		}
	})

	tests := []struct {
		name    string
		repoErr error
		wantErr error
	}{
		{name: "fail duplicate user", repoErr: domain.ErrDuplicateUser, wantErr: domain.BizErrUserAlreadyExists},
		{name: "fail not found", repoErr: domain.ErrUserNotFound},
		{name: "fail repository error", repoErr: errors.New("database unavailable")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mock.NewMockUserRepository(gomock.NewController(t))
			svc := NewUserService(repo)
			repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(tt.repoErr)

			err := svc.Update(context.Background(), &domain.UpdateUserReq{ID: "user-id"})
			wantErr := tt.wantErr

			if wantErr == nil {
				wantErr = tt.repoErr
			}
			if !errors.Is(err, wantErr) {
				t.Fatalf("Update() error = %v, want %v", err, wantErr)
			}
		})
	}
}

func TestUserServiceDelete(t *testing.T) {
	tests := []struct {
		name    string
		repoErr error
	}{
		{name: "success"},
		{name: "fail not found", repoErr: domain.ErrUserNotFound},
		{name: "fail repository error", repoErr: errors.New("database unavailable")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mock.NewMockUserRepository(gomock.NewController(t))
			svc := NewUserService(repo)
			repo.EXPECT().Delete(gomock.Any(), "user-id").Return(tt.repoErr)

			err := svc.Delete(context.Background(), "user-id")

			if !errors.Is(err, tt.repoErr) {
				t.Fatalf("Delete() error = %v, want %v", err, tt.repoErr)
			}
		})
	}
}

func TestUserServiceCount(t *testing.T) {
	tests := []struct {
		name      string
		repoCount int
		repoErr   error
		wantCount int
	}{
		{name: "success", repoCount: 42, wantCount: 42},
		{name: "fail repository error", repoErr: errors.New("database unavailable")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mock.NewMockUserRepository(gomock.NewController(t))
			svc := NewUserService(repo)
			repo.EXPECT().Count(gomock.Any()).Return(tt.repoCount, tt.repoErr)

			got, err := svc.Count(context.Background())

			if !errors.Is(err, tt.repoErr) {
				t.Fatalf("Count() error = %v, want %v", err, tt.repoErr)
			}
			if got != tt.wantCount {
				t.Fatalf("Count() = %d, want %d", got, tt.wantCount)
			}
		})
	}
}
