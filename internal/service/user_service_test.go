package service

import (
	"context"
	"testing"

	"room-booking/internal/auth"
	"room-booking/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepo struct {
	createFn     func(ctx context.Context, email, passwordHash string) (*domain.User, error)
	getByEmailFn func(ctx context.Context, email string) (*domain.User, string, error)
}

func (f *fakeUserRepo) Create(ctx context.Context, email, passwordHash string) (*domain.User, error) {
	return f.createFn(ctx, email, passwordHash)
}

func (f *fakeUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, string, error) {
	return f.getByEmailFn(ctx, email)
}

func TestUserServiceRegister_Success(t *testing.T) {
	repo := &fakeUserRepo{
		createFn: func(ctx context.Context, email, passwordHash string) (*domain.User, error) {
			if email != "newuser@example.com" {
				t.Fatalf("unexpected email: %s", email)
			}
			if passwordHash == "" {
				t.Fatal("expected password hash")
			}
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

	svc := NewUserService(repo, auth.NewJWTManager("test-secret"))

	user, err := svc.Register(context.Background(), "newuser@example.com", "secret123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Email != "newuser@example.com" {
		t.Fatalf("unexpected user email: %s", user.Email)
	}
}

func TestUserServiceRegister_InvalidRequest(t *testing.T) {
	repo := &fakeUserRepo{
		createFn: func(ctx context.Context, email, passwordHash string) (*domain.User, error) {
			t.Fatal("repo.Create should not be called")
			return nil, nil
		},
		getByEmailFn: func(ctx context.Context, email string) (*domain.User, string, error) {
			return nil, "", nil
		},
	}

	svc := NewUserService(repo, auth.NewJWTManager("test-secret"))

	_, err := svc.Register(context.Background(), "", "")
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got %v", err)
	}
}

func TestUserServiceLogin_Success(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	repo := &fakeUserRepo{
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

	jwtManager := auth.NewJWTManager("test-secret")
	svc := NewUserService(repo, jwtManager)

	token, err := svc.Login(context.Background(), "newuser@example.com", "secret123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	claims, err := jwtManager.ParseToken(token)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}

	if claims.UserID != "992a856b-a86a-4046-bbb6-4ad3ef0efedf" {
		t.Fatalf("unexpected user id in token: %s", claims.UserID)
	}
	if claims.Role != "user" {
		t.Fatalf("unexpected role in token: %s", claims.Role)
	}
}

func TestUserServiceLogin_Unauthorized(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	repo := &fakeUserRepo{
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

	svc := NewUserService(repo, auth.NewJWTManager("test-secret"))

	_, err = svc.Login(context.Background(), "newuser@example.com", "wrong")
	if err != domain.ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}
