package service

import (
	"context"
	"testing"
	"time"

	"room-booking/internal/domain"
)

type fakeScheduleRoomRepo struct {
	existsFn func(ctx context.Context, roomID string) (bool, error)
}

func (f *fakeScheduleRoomRepo) Exists(ctx context.Context, roomID string) (bool, error) {
	return f.existsFn(ctx, roomID)
}

type fakeScheduleRepo struct {
	existsByRoomIDFn func(ctx context.Context, roomID string) (bool, error)
	createFn         func(ctx context.Context, roomID string, daysOfWeek []int, startTime, endTime string) (*domain.Schedule, error)
}

func (f *fakeScheduleRepo) ExistsByRoomID(ctx context.Context, roomID string) (bool, error) {
	return f.existsByRoomIDFn(ctx, roomID)
}

func (f *fakeScheduleRepo) Create(ctx context.Context, roomID string, daysOfWeek []int, startTime, endTime string) (*domain.Schedule, error) {
	return f.createFn(ctx, roomID, daysOfWeek, startTime, endTime)
}

type fakeScheduleSlotRepo struct {
	createFn func(ctx context.Context, roomID string, startAt, endAt time.Time) error
}

func (f *fakeScheduleSlotRepo) Create(ctx context.Context, roomID string, startAt, endAt time.Time) error {
	return f.createFn(ctx, roomID, startAt, endAt)
}

func TestScheduleServiceCreate_InvalidRoomID(t *testing.T) {
	svc := NewScheduleService(
		&fakeScheduleRoomRepo{existsFn: func(ctx context.Context, roomID string) (bool, error) { return true, nil }},
		&fakeScheduleRepo{
			existsByRoomIDFn: func(ctx context.Context, roomID string) (bool, error) { return false, nil },
			createFn: func(ctx context.Context, roomID string, daysOfWeek []int, startTime, endTime string) (*domain.Schedule, error) {
				return nil, nil
			},
		},
		&fakeScheduleSlotRepo{createFn: func(ctx context.Context, roomID string, startAt, endAt time.Time) error { return nil }},
	)

	_, err := svc.Create(context.Background(), "", []int{1}, "09:00", "18:00")
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got %v", err)
	}
}

func TestScheduleServiceCreate_InvalidDays(t *testing.T) {
	svc := NewScheduleService(
		&fakeScheduleRoomRepo{existsFn: func(ctx context.Context, roomID string) (bool, error) { return true, nil }},
		&fakeScheduleRepo{
			existsByRoomIDFn: func(ctx context.Context, roomID string) (bool, error) { return false, nil },
			createFn: func(ctx context.Context, roomID string, daysOfWeek []int, startTime, endTime string) (*domain.Schedule, error) {
				return nil, nil
			},
		},
		&fakeScheduleSlotRepo{createFn: func(ctx context.Context, roomID string, startAt, endAt time.Time) error { return nil }},
	)

	_, err := svc.Create(context.Background(), "room-1", []int{0, 8}, "09:00", "18:00")
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got %v", err)
	}
}

func TestScheduleServiceCreate_InvalidTimeOrder(t *testing.T) {
	svc := NewScheduleService(
		&fakeScheduleRoomRepo{existsFn: func(ctx context.Context, roomID string) (bool, error) { return true, nil }},
		&fakeScheduleRepo{
			existsByRoomIDFn: func(ctx context.Context, roomID string) (bool, error) { return false, nil },
			createFn: func(ctx context.Context, roomID string, daysOfWeek []int, startTime, endTime string) (*domain.Schedule, error) {
				return nil, nil
			},
		},
		&fakeScheduleSlotRepo{createFn: func(ctx context.Context, roomID string, startAt, endAt time.Time) error { return nil }},
	)

	_, err := svc.Create(context.Background(), "room-1", []int{1, 2}, "18:00", "09:00")
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got %v", err)
	}
}

func TestScheduleServiceCreate_RoomNotFound(t *testing.T) {
	svc := NewScheduleService(
		&fakeScheduleRoomRepo{existsFn: func(ctx context.Context, roomID string) (bool, error) { return false, nil }},
		&fakeScheduleRepo{
			existsByRoomIDFn: func(ctx context.Context, roomID string) (bool, error) { return false, nil },
			createFn: func(ctx context.Context, roomID string, daysOfWeek []int, startTime, endTime string) (*domain.Schedule, error) {
				return nil, nil
			},
		},
		&fakeScheduleSlotRepo{createFn: func(ctx context.Context, roomID string, startAt, endAt time.Time) error { return nil }},
	)

	_, err := svc.Create(context.Background(), "room-1", []int{1, 2}, "09:00", "18:00")
	if err != domain.ErrRoomNotFound {
		t.Fatalf("expected ErrRoomNotFound, got %v", err)
	}
}

func TestScheduleServiceCreate_ScheduleExists(t *testing.T) {
	svc := NewScheduleService(
		&fakeScheduleRoomRepo{existsFn: func(ctx context.Context, roomID string) (bool, error) { return true, nil }},
		&fakeScheduleRepo{
			existsByRoomIDFn: func(ctx context.Context, roomID string) (bool, error) { return true, nil },
			createFn: func(ctx context.Context, roomID string, daysOfWeek []int, startTime, endTime string) (*domain.Schedule, error) {
				return nil, nil
			},
		},
		&fakeScheduleSlotRepo{createFn: func(ctx context.Context, roomID string, startAt, endAt time.Time) error { return nil }},
	)

	_, err := svc.Create(context.Background(), "room-1", []int{1, 2}, "09:00", "18:00")
	if err != domain.ErrScheduleExists {
		t.Fatalf("expected ErrScheduleExists, got %v", err)
	}
}

func TestScheduleServiceCreate_Success(t *testing.T) {
	slotCalls := 0

	svc := NewScheduleService(
		&fakeScheduleRoomRepo{existsFn: func(ctx context.Context, roomID string) (bool, error) { return true, nil }},
		&fakeScheduleRepo{
			existsByRoomIDFn: func(ctx context.Context, roomID string) (bool, error) { return false, nil },
			createFn: func(ctx context.Context, roomID string, daysOfWeek []int, startTime, endTime string) (*domain.Schedule, error) {
				return &domain.Schedule{
					ID:         "schedule-1",
					RoomID:     roomID,
					DaysOfWeek: daysOfWeek,
					StartTime:  startTime,
					EndTime:    endTime,
				}, nil
			},
		},
		&fakeScheduleSlotRepo{
			createFn: func(ctx context.Context, roomID string, startAt, endAt time.Time) error {
				slotCalls++
				return nil
			},
		},
	)

	schedule, err := svc.Create(context.Background(), "room-1", []int{1, 2, 3, 4, 5}, "09:00", "10:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if schedule.ID != "schedule-1" {
		t.Fatalf("unexpected schedule id: %s", schedule.ID)
	}
	if slotCalls == 0 {
		t.Fatal("expected generated slots to be stored")
	}
}
