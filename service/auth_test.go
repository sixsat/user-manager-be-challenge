package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sixsat/user-manager-be-challenge/domain"
	"github.com/sixsat/user-manager-be-challenge/mock"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthServiceRegister(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := mock.NewMockUserRepository(gomock.NewController(t))
		svc := NewAuthService(repo, "secret", time.Hour)
		repo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, req *domain.CreateUserReq) error {
				if req.Name != "Jane Doe" {
					t.Fatalf("name = %q, want %q", req.Name, "Jane Doe")
				}
				if req.Email != "jane@example.com" {
					t.Fatalf("email = %q, want %q", req.Email, "jane@example.com")
				}
				if req.PasswordHash == "password" {
					t.Fatal("password was not hashed")
				}
				if err := bcrypt.CompareHashAndPassword([]byte(req.PasswordHash), []byte("password")); err != nil {
					t.Fatalf("stored password hash does not match: %v", err)
				}
				return nil
			},
		)

		err := svc.Register(context.Background(), &domain.RegisterUserReq{
			Name:         "  Jane Doe  ",
			Email:        "  JANE@Example.COM ",
			PasswordHash: "password",
		})

		if err != nil {
			t.Fatalf("Register() error = %v", err)
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
			svc := NewAuthService(repo, "secret", time.Hour)
			repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(tt.repoErr)

			err := svc.Register(context.Background(), &domain.RegisterUserReq{PasswordHash: "password"})

			wantErr := tt.wantErr
			if wantErr == nil {
				wantErr = tt.repoErr
			}
			if !errors.Is(err, wantErr) {
				t.Fatalf("Register() error = %v, want %v", err, wantErr)
			}
		})
	}

	t.Run("fail password too long", func(t *testing.T) {
		repo := mock.NewMockUserRepository(gomock.NewController(t))
		svc := NewAuthService(repo, "secret", time.Hour)

		err := svc.Register(context.Background(), &domain.RegisterUserReq{PasswordHash: string(make([]byte, 73))})

		if !errors.Is(err, bcrypt.ErrPasswordTooLong) {
			t.Fatalf("Register() error = %v, want %v", err, bcrypt.ErrPasswordTooLong)
		}
	})
}

func TestAuthServiceLogin(t *testing.T) {
	const (
		userID    = "507f1f77bcf86cd799439011"
		signKey   = "test-signing-key"
		password  = "correct-password"
		jwtExpiry = 30 * time.Minute
	)

	t.Run("success", func(t *testing.T) {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			t.Fatalf("GenerateFromPassword() error = %v", err)
		}

		repo := mock.NewMockUserRepository(gomock.NewController(t))
		svc := NewAuthService(repo, signKey, jwtExpiry)
		repo.EXPECT().GetByEmail(gomock.Any(), "jane@example.com").Return(&domain.GetUserByEmailRes{
			ID:           userID,
			PasswordHash: string(hash),
		}, nil)
		before := time.Now()
		res, err := svc.Login(context.Background(), &domain.LoginUserReq{
			Email:    "  JANE@Example.COM ",
			Password: password,
		})
		if err != nil {
			t.Fatalf("Login() error = %v", err)
		}

		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(res.AccessToken, claims, func(token *jwt.Token) (any, error) {
			return []byte(signKey), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
		if err != nil {
			t.Fatalf("ParseWithClaims() error = %v", err)
		}
		if !token.Valid {
			t.Fatal("issued token is invalid")
		}
		if claims.Subject != userID {
			t.Fatalf("subject = %q, want %q", claims.Subject, userID)
		}
		if claims.IssuedAt == nil || claims.IssuedAt.Before(before.Add(-time.Second)) {
			t.Fatalf("issued_at = %v, want at or after %v", claims.IssuedAt, before)
		}
		if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Sub(claims.IssuedAt.Time) != jwtExpiry {
			t.Fatalf("token lifetime = %v, want %v", claims.ExpiresAt.Time.Sub(claims.IssuedAt.Time), jwtExpiry)
		}
	})

	tests := []struct {
		name       string
		password   string
		storedHash string
		repoErr    error
		wantErr    error
	}{
		{
			name:    "fail not found",
			repoErr: domain.ErrUserNotFound,
			wantErr: domain.BizErrInvalidCredentials,
		},
		{
			name:    "fail repository error",
			repoErr: errors.New("database unavailable"),
		},
		{
			name:       "fail invalid credentials",
			password:   "wrong-password",
			storedHash: hashPassword(t, password),
			wantErr:    domain.BizErrInvalidCredentials,
		},
		{
			name:       "fail password hash error",
			password:   password,
			storedHash: "not-a-bcrypt-hash",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mock.NewMockUserRepository(gomock.NewController(t))
			svc := NewAuthService(repo, signKey, jwtExpiry)
			var user *domain.GetUserByEmailRes
			if tt.repoErr == nil {
				user = &domain.GetUserByEmailRes{ID: userID, PasswordHash: tt.storedHash}
			}
			repo.EXPECT().GetByEmail(gomock.Any(), "jane@example.com").Return(user, tt.repoErr)

			_, err := svc.Login(context.Background(), &domain.LoginUserReq{
				Email:    "jane@example.com",
				Password: tt.password,
			})

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Login() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if tt.repoErr != nil && !errors.Is(err, tt.repoErr) {
				t.Fatalf("Login() error = %v, want %v", err, tt.repoErr)
			}
			if tt.repoErr == nil && err == nil {
				t.Fatal("Login() error = nil, want password hash error")
			}
		})
	}
}

func hashPassword(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}
	return string(hash)
}
