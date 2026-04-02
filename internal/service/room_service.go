package service

import (
	"context"
	"strings"

	"room-booking/internal/domain"
	"room-booking/internal/repository"
)

type RoomService struct {
	repo *repository.RoomRepository
}

func NewRoomService(repo *repository.RoomRepository) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService) Create(ctx context.Context, name string, description *string, capacity *int) (*domain.Room, error) {
	if strings.TrimSpace(name) == "" {
		return nil, domain.ErrInvalidRequest
	}

	if capacity != nil && *capacity < 0 {
		return nil, domain.ErrInvalidRequest
	}

	return s.repo.Create(ctx, strings.TrimSpace(name), description, capacity)
}

func (s *RoomService) List(ctx context.Context) ([]domain.Room, error) {
	return s.repo.List(ctx)
}
