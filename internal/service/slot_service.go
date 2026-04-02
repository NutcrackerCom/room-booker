package service

import (
	"context"
	"time"

	"room-booking/internal/domain"
)

type slotRoomRepository interface {
	Exists(ctx context.Context, roomID string) (bool, error)
}

type slotScheduleRepository interface {
	ExistsByRoomID(ctx context.Context, roomID string) (bool, error)
}

type slotRepository interface {
	ListAvailable(ctx context.Context, roomID string, dayStart, dayEnd time.Time) ([]domain.Slot, error)
}

type SlotService struct {
	roomRepo     slotRoomRepository
	scheduleRepo slotScheduleRepository
	slotRepo     slotRepository
}

func NewSlotService(
	roomRepo slotRoomRepository,
	scheduleRepo slotScheduleRepository,
	slotRepo slotRepository,
) *SlotService {
	return &SlotService{
		roomRepo:     roomRepo,
		scheduleRepo: scheduleRepo,
		slotRepo:     slotRepo,
	}
}

func (s *SlotService) ListAvailable(ctx context.Context, roomID, date string) ([]domain.Slot, error) {
	if roomID == "" || date == "" {
		return nil, domain.ErrInvalidRequest
	}

	day, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, domain.ErrInvalidRequest
	}

	roomExists, err := s.roomRepo.Exists(ctx, roomID)
	if err != nil {
		return nil, err
	}
	if !roomExists {
		return nil, domain.ErrRoomNotFound
	}

	scheduleExists, err := s.scheduleRepo.ExistsByRoomID(ctx, roomID)
	if err != nil {
		return nil, err
	}
	if !scheduleExists {
		return []domain.Slot{}, nil
	}

	dayStart := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.Add(24 * time.Hour)

	return s.slotRepo.ListAvailable(ctx, roomID, dayStart, dayEnd)
}
