package service

import (
	"context"
	"strings"

	"room-booking/internal/auth"
	"room-booking/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type userRepository interface {
	Create(ctx context.Context, email, passwordHash string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, string, error)
}

type UserService struct {
	repo       userRepository
	jwtManager *auth.JWTManager
}

func NewUserService(repo userRepository, jwtManager *auth.JWTManager) *UserService {
	return &UserService{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

func (s *UserService) Register(ctx context.Context, email, password string) (*domain.User, error) {
	email = strings.TrimSpace(email)
	if email == "" || password == "" {
		return nil, domain.ErrInvalidRequest
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, email, string(hash))
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	email = strings.TrimSpace(email)
	if email == "" || password == "" {
		return "", domain.ErrInvalidRequest
	}

	user, passwordHash, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", domain.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return "", domain.ErrUnauthorized
	}

	token, err := s.jwtManager.CreateToken(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}
