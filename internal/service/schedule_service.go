package service

import (
	"context"
	"sort"
	"time"

	"room-booking/internal/domain"
	"room-booking/internal/repository"
	"room-booking/internal/slots"
)

type ScheduleService struct {
	roomRepo     *repository.RoomRepository
	scheduleRepo *repository.ScheduleRepository
	slotRepo     *repository.SlotRepository
}

func NewScheduleService(
	roomRepo *repository.RoomRepository,
	scheduleRepo *repository.ScheduleRepository,
	slotRepo *repository.SlotRepository,
) *ScheduleService {
	return &ScheduleService{
		roomRepo:     roomRepo,
		scheduleRepo: scheduleRepo,
		slotRepo:     slotRepo,
	}
}

func (s *ScheduleService) Create(ctx context.Context, roomID string, daysOfWeek []int, startTime, endTime string) (*domain.Schedule, error) {
	if roomID == "" {
		return nil, domain.ErrInvalidRequest
	}

	if len(daysOfWeek) == 0 {
		return nil, domain.ErrInvalidRequest
	}

	normalizedDays := make([]int, 0, len(daysOfWeek))
	seen := make(map[int]bool)

	for _, day := range daysOfWeek {
		if day < 1 || day > 7 {
			return nil, domain.ErrInvalidRequest
		}
		if !seen[day] {
			seen[day] = true
			normalizedDays = append(normalizedDays, day)
		}
	}

	sort.Ints(normalizedDays)

	startParsed, err := time.Parse("15:04", startTime)
	if err != nil {
		return nil, domain.ErrInvalidRequest
	}

	endParsed, err := time.Parse("15:04", endTime)
	if err != nil {
		return nil, domain.ErrInvalidRequest
	}

	if !startParsed.Before(endParsed) {
		return nil, domain.ErrInvalidRequest
	}

	exists, err := s.roomRepo.Exists(ctx, roomID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrRoomNotFound
	}

	scheduleExists, err := s.scheduleRepo.ExistsByRoomID(ctx, roomID)
	if err != nil {
		return nil, err
	}
	if scheduleExists {
		return nil, domain.ErrScheduleExists
	}

	schedule, err := s.scheduleRepo.Create(ctx, roomID, normalizedDays, startTime, endTime)
	if err != nil {
		return nil, err
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	for i := 0; i < 30; i++ {
		date := today.AddDate(0, 0, i)
		generated, err := slots.BuildSlotsForDate(date, normalizedDays, startTime, endTime)
		if err != nil {
			return nil, err
		}

		for _, slot := range generated {
			if err := s.slotRepo.Create(ctx, roomID, slot.Start, slot.End); err != nil {
				return nil, err
			}
		}
	}

	return schedule, nil
}
