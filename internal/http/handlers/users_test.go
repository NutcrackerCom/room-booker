package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"room-booking/internal/auth"
	"room-booking/internal/domain"
	"room-booking/internal/service"

	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepoForHandler struct {
	createFn     func(ctx context.Context, email, passwordHash string) (*domain.User, error)
	getByEmailFn func(ctx context.Context, email string) (*domain.User, string, error)
}

func (f *fakeUserRepoForHandler) Create(ctx context.Context, email, passwordHash string) (*domain.User, error) {
	return f.createFn(ctx, email, passwordHash)
}

func (f *fakeUserRepoForHandler) GetByEmail(ctx context.Context, email string) (*domain.User, string, error) {
	return f.getByEmailFn(ctx, email)
}

func TestUserHandlerRegister_Success(t *testing.T) {
	repo := &fakeUserRepoForHandler{
		createFn: func(ctx context.Context, email, passwordHash string) (*domain.User, error) {
			return &domain.User{
				ID:    "user-1",
				Email: email,
				Role:  "user",
			}, nil
		},
		getByEmailFn: func(ctx context.Context, email string) (*domain.User, string, error) {
			return nil, "", nil
		},
	}

	handler := NewUserHandler(service.NewUserService(repo, auth.NewJWTManager("test-secret")))

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(`{"email":"newuser@example.com","password":"secret123"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), `"email":"newuser@example.com"`) {
		t.Fatalf("unexpected body: %s", rr.Body.String())
	}
}

func TestUserHandlerRegister_BadJSON(t *testing.T) {
	repo := &fakeUserRepoForHandler{
		createFn: func(ctx context.Context, email, passwordHash string) (*domain.User, error) {
			t.Fatal("repo.Create should not be called")
			return nil, nil
		},
		getByEmailFn: func(ctx context.Context, email string) (*domain.User, string, error) {
			return nil, "", nil
		},
	}

	handler := NewUserHandler(service.NewUserService(repo, auth.NewJWTManager("test-secret")))

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(`{bad json}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", rr.Code, rr.Body.String())
	}
}

func TestUserHandlerLogin_Success(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	repo := &fakeUserRepoForHandler{
		createFn: func(ctx context.Context, email, passwordHash string) (*domain.User, error) {
			return nil, nil
		},
		getByEmailFn: func(ctx context.Context, email string) (*domain.User, string, error) {
			return &domain.User{
				ID:    "992a856b-a86a-4046-bbb6-4ad3ef0efedf",
				Email: email,
				Role:  "user",
			}, string(hash), nil
		},
	}

	handler := NewUserHandler(service.NewUserService(repo, auth.NewJWTManager("test-secret")))

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(`{"email":"newuser@example.com","password":"secret123"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), `"token":"`) {
		t.Fatalf("unexpected body: %s", rr.Body.String())
	}
}

func TestUserHandlerLogin_Unauthorized(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	repo := &fakeUserRepoForHandler{
		createFn: func(ctx context.Context, email, passwordHash string) (*domain.User, error) {
			return nil, nil
		},
		getByEmailFn: func(ctx context.Context, email string) (*domain.User, string, error) {
			return &domain.User{
				ID:    "user-1",
				Email: email,
				Role:  "user",
			}, string(hash), nil
		},
	}

	handler := NewUserHandler(service.NewUserService(repo, auth.NewJWTManager("test-secret")))

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(`{"email":"newuser@example.com","password":"wrong"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d, body=%s", rr.Code, rr.Body.String())
	}
}
