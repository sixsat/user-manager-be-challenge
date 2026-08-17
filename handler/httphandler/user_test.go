package httphandler

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/sixsat/user-manager-be-challenge/domain"
	"github.com/sixsat/user-manager-be-challenge/mock"
	"go.uber.org/mock/gomock"
)

const userID = "6a82eb18aec4953bd25e7556"

func TestCreateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userSvc := mock.NewMockUserService(gomock.NewController(t))
		h := newTestHandler(nil, userSvc)
		userSvc.EXPECT().Create(gomock.Any(), &domain.CreateUserReq{
			Name:         "John Doe",
			Email:        "john@example.com",
			PasswordHash: "password",
		}).Return(nil)
		c, recorder := newTestContext(http.MethodPost, "/api/users", `{"name":"John Doe","email":"john@example.com","password":"password"}`)

		if err := h.createUser(c); err != nil {
			t.Fatalf("createUser() error = %v", err)
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
			userSvc := mock.NewMockUserService(gomock.NewController(t))
			h := newTestHandler(nil, userSvc)
			c, recorder := newTestContext(http.MethodPost, "/api/users", tt.body)

			if err := h.createUser(c); err != nil {
				t.Fatalf("createUser() error = %v", err)
			}

			assertResponse(t, recorder, http.StatusBadRequest, CodeBadReq, DescBadReq)
		})
	}

	t.Run("fail service error", func(t *testing.T) {
		wantErr := errors.New("service unavailable")
		userSvc := mock.NewMockUserService(gomock.NewController(t))
		h := newTestHandler(nil, userSvc)
		userSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(wantErr)
		c, _ := newTestContext(http.MethodPost, "/api/users", `{"name":"John Doe","email":"john@example.com","password":"password"}`)

		if err := h.createUser(c); !errors.Is(err, wantErr) {
			t.Fatalf("createUser() error = %v, want %v", err, wantErr)
		}
	})
}

func TestListUsers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		createdAt := time.Date(2026, time.August, 17, 12, 0, 0, 0, time.UTC)
		userSvc := mock.NewMockUserService(gomock.NewController(t))
		h := newTestHandler(nil, userSvc)
		userSvc.EXPECT().List(gomock.Any()).Return([]domain.GetUserRes{{
			ID: userID, Name: "John Doe", Email: "john@example.com", CreatedAt: createdAt,
		}}, nil)
		c, recorder := newTestContext(http.MethodGet, "/api/users", "")

		if err := h.listUsers(c); err != nil {
			t.Fatalf("listUsers() error = %v", err)
		}

		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
		}
		res := decodeResponse[[]GetUserRes](t, recorder)
		if res.Data == nil || len(*res.Data) != 1 || (*res.Data)[0].ID != userID || !(*res.Data)[0].CreatedAt.Equal(createdAt) {
			t.Fatalf("response = %#v", res)
		}
	})

	t.Run("fail service error", func(t *testing.T) {
		wantErr := errors.New("service unavailable")
		userSvc := mock.NewMockUserService(gomock.NewController(t))
		h := newTestHandler(nil, userSvc)
		userSvc.EXPECT().List(gomock.Any()).Return(nil, wantErr)
		c, _ := newTestContext(http.MethodGet, "/api/users", "")

		if err := h.listUsers(c); !errors.Is(err, wantErr) {
			t.Fatalf("listUsers() error = %v, want %v", err, wantErr)
		}
	})
}

func TestGetUserByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		createdAt := time.Date(2026, time.August, 17, 12, 0, 0, 0, time.UTC)
		userSvc := mock.NewMockUserService(gomock.NewController(t))
		h := newTestHandler(nil, userSvc)
		userSvc.EXPECT().GetByID(gomock.Any(), userID).Return(&domain.GetUserRes{
			ID: userID, Name: "John Doe", Email: "john@example.com", CreatedAt: createdAt,
		}, nil)
		c, recorder := newTestContext(http.MethodGet, "/api/users/"+userID, "")
		c.SetPathValues(echo.PathValues{{Name: "id", Value: userID}})

		if err := h.getUserByID(c); err != nil {
			t.Fatalf("getUserByID() error = %v", err)
		}

		res := decodeResponse[GetUserRes](t, recorder)
		if recorder.Code != http.StatusOK || res.Data == nil || res.Data.ID != userID || res.Data.Email != "john@example.com" || !res.Data.CreatedAt.Equal(createdAt) {
			t.Fatalf("status = %d, response = %#v", recorder.Code, res)
		}
	})

	t.Run("fail binding request error", func(t *testing.T) {
		userSvc := mock.NewMockUserService(gomock.NewController(t))
		h := newTestHandler(nil, userSvc)
		c, recorder := newTestContext(http.MethodGet, "/api/users/"+userID, "")
		c.Echo().Binder = failingBinder{err: errors.New("binding failed")}

		if err := h.getUserByID(c); err != nil {
			t.Fatalf("getUserByID() error = %v", err)
		}

		assertResponse(t, recorder, http.StatusBadRequest, CodeBadReq, DescBadReq)
	})

	t.Run("fail invalid request", func(t *testing.T) {
		userSvc := mock.NewMockUserService(gomock.NewController(t))
		h := newTestHandler(nil, userSvc)
		c, recorder := newTestContext(http.MethodGet, "/api/users/invalid", "")
		c.SetPathValues(echo.PathValues{{Name: "id", Value: "invalid"}})

		if err := h.getUserByID(c); err != nil {
			t.Fatalf("getUserByID() error = %v", err)
		}

		assertResponse(t, recorder, http.StatusBadRequest, CodeBadReq, DescBadReq)
	})

	t.Run("fail not found", func(t *testing.T) {
		userSvc := mock.NewMockUserService(gomock.NewController(t))
		h := newTestHandler(nil, userSvc)
		userSvc.EXPECT().GetByID(gomock.Any(), userID).Return(nil, domain.ErrUserNotFound)
		c, recorder := newTestContext(http.MethodGet, "/api/users/"+userID, "")
		c.SetPathValues(echo.PathValues{{Name: "id", Value: userID}})

		if err := h.getUserByID(c); err != nil {
			t.Fatalf("getUserByID() error = %v", err)
		}

		assertResponse(t, recorder, http.StatusNotFound, CodeBadReq, DescBadReq)
	})

	t.Run("fail service error", func(t *testing.T) {
		wantErr := errors.New("service unavailable")
		userSvc := mock.NewMockUserService(gomock.NewController(t))
		h := newTestHandler(nil, userSvc)
		userSvc.EXPECT().GetByID(gomock.Any(), userID).Return(nil, wantErr)
		c, _ := newTestContext(http.MethodGet, "/api/users/"+userID, "")
		c.SetPathValues(echo.PathValues{{Name: "id", Value: userID}})

		if err := h.getUserByID(c); !errors.Is(err, wantErr) {
			t.Fatalf("getUserByID() error = %v, want %v", err, wantErr)
		}
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userSvc := mock.NewMockUserService(gomock.NewController(t))
		h := newTestHandler(nil, userSvc)
		userSvc.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, req *domain.UpdateUserReq) error {
				if req.ID != userID || req.Name == nil || *req.Name != "John Doe" || req.Email == nil || *req.Email != "john@example.com" {
					t.Fatalf("Update() request = %#v", req)
				}
				return nil
			},
		)
		c, recorder := newTestContext(http.MethodPatch, "/api/users/"+userID, `{"name":"John Doe","email":"john@example.com"}`)
		c.SetPathValues(echo.PathValues{{Name: "id", Value: userID}})

		if err := h.updateUser(c); err != nil {
			t.Fatalf("updateUser() error = %v", err)
		}

		assertResponse(t, recorder, http.StatusOK, CodeOK, DescOK)
	})

	tests := []struct {
		name string
		body string
	}{
		{name: "fail malformed JSON", body: `{"name":`},
		{name: "fail empty payload", body: `{}`},
		{name: "fail invalid request", body: `{"email":"invalid"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userSvc := mock.NewMockUserService(gomock.NewController(t))
			h := newTestHandler(nil, userSvc)
			c, recorder := newTestContext(http.MethodPatch, "/api/users/"+userID, tt.body)
			c.SetPathValues(echo.PathValues{{Name: "id", Value: userID}})

			if err := h.updateUser(c); err != nil {
				t.Fatalf("updateUser() error = %v", err)
			}

			assertResponse(t, recorder, http.StatusBadRequest, CodeBadReq, DescBadReq)
		})
	}

	t.Run("fail not found", func(t *testing.T) {
		userSvc := mock.NewMockUserService(gomock.NewController(t))
		h := newTestHandler(nil, userSvc)
		userSvc.EXPECT().Update(gomock.Any(), gomock.Any()).Return(domain.ErrUserNotFound)
		c, recorder := newTestContext(http.MethodPatch, "/api/users/"+userID, `{"name":"John"}`)
		c.SetPathValues(echo.PathValues{{Name: "id", Value: userID}})

		if err := h.updateUser(c); err != nil {
			t.Fatalf("updateUser() error = %v", err)
		}

		assertResponse(t, recorder, http.StatusNotFound, CodeBadReq, DescBadReq)
	})

	t.Run("fail service error", func(t *testing.T) {
		wantErr := errors.New("service unavailable")
		userSvc := mock.NewMockUserService(gomock.NewController(t))
		h := newTestHandler(nil, userSvc)
		userSvc.EXPECT().Update(gomock.Any(), gomock.Any()).Return(wantErr)
		c, _ := newTestContext(http.MethodPatch, "/api/users/"+userID, `{"name":"John"}`)
		c.SetPathValues(echo.PathValues{{Name: "id", Value: userID}})

		if err := h.updateUser(c); !errors.Is(err, wantErr) {
			t.Fatalf("updateUser() error = %v, want %v", err, wantErr)
		}
	})
}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name       string
		serviceErr error
		wantStatus int
		wantCode   string
		wantDesc   string
	}{
		{name: "success", wantStatus: http.StatusOK, wantCode: CodeOK, wantDesc: DescOK},
		{name: "fail not found", serviceErr: domain.ErrUserNotFound, wantStatus: http.StatusNotFound, wantCode: CodeBadReq, wantDesc: DescBadReq},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userSvc := mock.NewMockUserService(gomock.NewController(t))
			h := newTestHandler(nil, userSvc)
			userSvc.EXPECT().Delete(gomock.Any(), userID).Return(tt.serviceErr)
			c, recorder := newTestContext(http.MethodDelete, "/api/users/"+userID, "")
			c.SetPathValues(echo.PathValues{{Name: "id", Value: userID}})

			if err := h.deleteUser(c); err != nil {
				t.Fatalf("deleteUser() error = %v", err)
			}

			assertResponse(t, recorder, tt.wantStatus, tt.wantCode, tt.wantDesc)
		})
	}

	t.Run("fail binding request error", func(t *testing.T) {
		userSvc := mock.NewMockUserService(gomock.NewController(t))
		h := newTestHandler(nil, userSvc)
		c, recorder := newTestContext(http.MethodDelete, "/api/users/"+userID, "")
		c.Echo().Binder = failingBinder{err: errors.New("binding failed")}

		if err := h.deleteUser(c); err != nil {
			t.Fatalf("deleteUser() error = %v", err)
		}

		assertResponse(t, recorder, http.StatusBadRequest, CodeBadReq, DescBadReq)
	})

	t.Run("fail invalid request", func(t *testing.T) {
		userSvc := mock.NewMockUserService(gomock.NewController(t))
		h := newTestHandler(nil, userSvc)
		c, recorder := newTestContext(http.MethodDelete, "/api/users/invalid", "")
		c.SetPathValues(echo.PathValues{{Name: "id", Value: "invalid"}})

		if err := h.deleteUser(c); err != nil {
			t.Fatalf("deleteUser() error = %v", err)
		}

		assertResponse(t, recorder, http.StatusBadRequest, CodeBadReq, DescBadReq)
	})

	t.Run("fail service error", func(t *testing.T) {
		wantErr := errors.New("service unavailable")
		userSvc := mock.NewMockUserService(gomock.NewController(t))
		h := newTestHandler(nil, userSvc)
		userSvc.EXPECT().Delete(gomock.Any(), userID).Return(wantErr)
		c, _ := newTestContext(http.MethodDelete, "/api/users/"+userID, "")
		c.SetPathValues(echo.PathValues{{Name: "id", Value: userID}})

		if err := h.deleteUser(c); !errors.Is(err, wantErr) {
			t.Fatalf("deleteUser() error = %v, want %v", err, wantErr)
		}
	})
}
