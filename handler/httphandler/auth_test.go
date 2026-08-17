package httphandler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/sixsat/user-manager-be-challenge/domain"
	"github.com/sixsat/user-manager-be-challenge/mock"
	"go.uber.org/mock/gomock"
)

func TestRegisterUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		authSvc := mock.NewMockAuthService(gomock.NewController(t))
		h := newTestHandler(authSvc, nil)
		authSvc.EXPECT().Register(gomock.Any(), &domain.RegisterUserReq{
			Name:         "Jane Doe",
			Email:        "jane@example.com",
			PasswordHash: "password",
		}).Return(nil)
		c, recorder := newTestContext(http.MethodPost, "/api/auth/register", `{"name":"Jane Doe","email":"jane@example.com","password":"password"}`)

		if err := h.registerUser(c); err != nil {
			t.Fatalf("registerUser() error = %v", err)
		}

		assertResponse(t, recorder, http.StatusCreated, CodeOK, DescOK)
	})

	tests := []struct {
		name string
		body string
	}{
		{name: "fail malformed JSON", body: `{"name":`},
		{name: "fail invalid request", body: `{"name":"","email":"invalid","password":""}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authSvc := mock.NewMockAuthService(gomock.NewController(t))
			h := newTestHandler(authSvc, nil)
			c, recorder := newTestContext(http.MethodPost, "/api/auth/register", tt.body)

			if err := h.registerUser(c); err != nil {
				t.Fatalf("registerUser() error = %v", err)
			}

			assertResponse(t, recorder, http.StatusBadRequest, CodeBadReq, DescBadReq)
		})
	}

	t.Run("fail service error", func(t *testing.T) {
		wantErr := errors.New("service unavailable")
		authSvc := mock.NewMockAuthService(gomock.NewController(t))
		h := newTestHandler(authSvc, nil)
		authSvc.EXPECT().Register(gomock.Any(), gomock.Any()).Return(wantErr)
		c, _ := newTestContext(http.MethodPost, "/api/auth/register", `{"name":"Jane Doe","email":"jane@example.com","password":"password"}`)

		if err := h.registerUser(c); !errors.Is(err, wantErr) {
			t.Fatalf("registerUser() error = %v, want %v", err, wantErr)
		}
	})
}

func TestLoginUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		authSvc := mock.NewMockAuthService(gomock.NewController(t))
		h := newTestHandler(authSvc, nil)
		authSvc.EXPECT().Login(gomock.Any(), &domain.LoginUserReq{
			Email:    "jane@example.com",
			Password: "password",
		}).Return(&domain.LoginUserRes{AccessToken: "signed-token"}, nil)
		c, recorder := newTestContext(http.MethodPost, "/api/auth/login", `{"email":"jane@example.com","password":"password"}`)

		if err := h.loginUser(c); err != nil {
			t.Fatalf("loginUser() error = %v", err)
		}

		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
		}
		res := decodeResponse[LoginUserRes](t, recorder)
		if res.Code != CodeOK || res.Desc != DescOK || res.Data == nil || res.Data.AccessToken != "signed-token" {
			t.Fatalf("response = %#v", res)
		}
	})

	tests := []struct {
		name string
		body string
	}{
		{name: "fail malformed JSON", body: `{"email":`},
		{name: "fail invalid request", body: `{"email":"invalid","password":""}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authSvc := mock.NewMockAuthService(gomock.NewController(t))
			h := newTestHandler(authSvc, nil)
			c, recorder := newTestContext(http.MethodPost, "/api/auth/login", tt.body)

			if err := h.loginUser(c); err != nil {
				t.Fatalf("loginUser() error = %v", err)
			}

			assertResponse(t, recorder, http.StatusBadRequest, CodeBadReq, DescBadReq)
		})
	}

	t.Run("fail invalid credentials", func(t *testing.T) {
		authSvc := mock.NewMockAuthService(gomock.NewController(t))
		h := newTestHandler(authSvc, nil)
		authSvc.EXPECT().Login(gomock.Any(), gomock.Any()).Return(nil, domain.BizErrInvalidCredentials)
		c, recorder := newTestContext(http.MethodPost, "/api/auth/login", `{"email":"jane@example.com","password":"wrong"}`)

		if err := h.loginUser(c); err != nil {
			t.Fatalf("loginUser() error = %v", err)
		}

		assertResponse(t, recorder, http.StatusUnauthorized, domain.BizErrInvalidCredentials.Code, domain.BizErrInvalidCredentials.Desc)
	})

	t.Run("fail service error", func(t *testing.T) {
		wantErr := errors.New("service unavailable")
		authSvc := mock.NewMockAuthService(gomock.NewController(t))
		h := newTestHandler(authSvc, nil)
		authSvc.EXPECT().Login(gomock.Any(), gomock.Any()).Return(nil, wantErr)
		c, _ := newTestContext(http.MethodPost, "/api/auth/login", `{"email":"jane@example.com","password":"password"}`)

		if err := h.loginUser(c); !errors.Is(err, wantErr) {
			t.Fatalf("loginUser() error = %v, want %v", err, wantErr)
		}
	})
}
