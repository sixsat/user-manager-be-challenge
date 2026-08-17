package httphandler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/sixsat/user-manager-be-challenge/port"
)

func TestRegisterRoutesProtectedEndpoints(t *testing.T) {
	e := echo.New()
	newTestHandler(nil, nil).RegisterRoutes(e.Group("/api"))
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)

	e.ServeHTTP(recorder, req)

	assertResponse(t, recorder, http.StatusUnauthorized, CodeBadReq, DescBadReq)
}

func newTestHandler(authSvc port.AuthService, userSvc port.UserService) *handler {
	return New(
		"test-signing-key",
		validator.New(validator.WithRequiredStructEnabled()),
		authSvc,
		userSvc,
	)
}

func assertResponse(t *testing.T, recorder *httptest.ResponseRecorder, wantStatus int, wantCode, wantDesc string) {
	t.Helper()
	if recorder.Code != wantStatus {
		t.Fatalf("status = %d, want %d", recorder.Code, wantStatus)
	}
	res := decodeResponse[any](t, recorder)
	if res.Code != wantCode || res.Desc != wantDesc {
		t.Fatalf("response = %#v, want code=%q desc=%q", res, wantCode, wantDesc)
	}
}

func decodeResponse[T any](t *testing.T, recorder *httptest.ResponseRecorder) Res[T] {
	t.Helper()
	var res Res[T]
	if err := json.NewDecoder(recorder.Body).Decode(&res); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return res
}

func newTestContext(method, target, body string) (*echo.Context, *httptest.ResponseRecorder) {
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, target, reader)
	if body != "" {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	recorder := httptest.NewRecorder()
	return echo.New().NewContext(req, recorder), recorder
}

type failingBinder struct {
	err error
}

func (b failingBinder) Bind(_ *echo.Context, _ any) error {
	return b.err
}
