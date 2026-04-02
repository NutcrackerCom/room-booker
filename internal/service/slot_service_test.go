package service

import (
	"context"
	"testing"
	"time"

	"room-booking/internal/domain"
)

type fakeSlotRoomRepo struct {
	existsFn func(ctx context.Context, roomID string) (bool, error)
}

func (f *fakeSlotRoomRepo) Exists(ctx context.Context, roomID string) (bool, error) {
	return f.existsFn(ctx, roomID)
}

type fakeSlotScheduleRepo struct {
	existsByRoomIDFn func(ctx context.Context, roomID string) (bool, error)
}

func (f *fakeSlotScheduleRepo) ExistsByRoomID(ctx context.Context, roomID string) (bool, error) {
	return f.existsByRoomIDFn(ctx, roomID)
}

type fakeSlotRepo struct {
	listAvailableFn func(ctx context.Context, roomID string, dayStart, dayEnd time.Time) ([]domain.Slot, error)
}

func (f *fakeSlotRepo) ListAvailable(ctx context.Context, roomID string, dayStart, dayEnd time.Time) ([]domain.Slot, error) {
	return f.listAvailableFn(ctx, roomID, dayStart, dayEnd)
}

func TestSlotServiceListAvailable_InvalidRequest(t *testing.T) {
	svc := NewSlotService(
		&fakeSlotRoomRepo{existsFn: func(ctx context.Context, roomID string) (bool, error) { return true, nil }},
		&fakeSlotScheduleRepo{existsByRoomIDFn: func(ctx context.Context, roomID string) (bool, error) { return true, nil }},
		&fakeSlotRepo{listAvailableFn: func(ctx context.Context, roomID string, dayStart, dayEnd time.Time) ([]domain.Slot, error) { return nil, nil }},
	)

	_, err := svc.ListAvailable(context.Background(), "", "")
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got %v", err)
	}
}

func TestSlotServiceListAvailable_RoomNotFound(t *testing.T) {
	svc := NewSlotService(
		&fakeSlotRoomRepo{existsFn: func(ctx context.Context, roomID string) (bool, error) { return false, nil }},
		&fakeSlotScheduleRepo{existsByRoomIDFn: func(ctx context.Context, roomID string) (bool, error) { return false, nil }},
		&fakeSlotRepo{listAvailableFn: func(ctx context.Context, roomID string, dayStart, dayEnd time.Time) ([]domain.Slot, error) { return nil, nil }},
	)

	_, err := svc.ListAvailable(context.Background(), "room-1", "2026-04-03")
	if err != domain.ErrRoomNotFound {
		t.Fatalf("expected ErrRoomNotFound, got %v", err)
	}
}

func TestSlotServiceListAvailable_NoSchedule(t *testing.T) {
	svc := NewSlotService(
		&fakeSlotRoomRepo{existsFn: func(ctx context.Context, roomID string) (bool, error) { return true, nil }},
		&fakeSlotScheduleRepo{existsByRoomIDFn: func(ctx context.Context, roomID string) (bool, error) { return false, nil }},
		&fakeSlotRepo{listAvailableFn: func(ctx context.Context, roomID string, dayStart, dayEnd time.Time) ([]domain.Slot, error) {
			t.Fatal("slot repo should not be called")
			return nil, nil
		}},
	)

	slots, err := svc.ListAvailable(context.Background(), "room-1", "2026-04-03")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(slots) != 0 {
		t.Fatalf("expected empty slots, got %d", len(slots))
	}
}

func TestSlotServiceListAvailable_Success(t *testing.T) {
	expected := []domain.Slot{{ID: "slot-1", RoomID: "room-1"}}

	svc := NewSlotService(
		&fakeSlotRoomRepo{existsFn: func(ctx context.Context, roomID string) (bool, error) { return true, nil }},
		&fakeSlotScheduleRepo{existsByRoomIDFn: func(ctx context.Context, roomID string) (bool, error) { return true, nil }},
		&fakeSlotRepo{listAvailableFn: func(ctx context.Context, roomID string, dayStart, dayEnd time.Time) ([]domain.Slot, error) {
			return expected, nil
		}},
	)

	slots, err := svc.ListAvailable(context.Background(), "room-1", "2026-04-03")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(slots) != 1 || slots[0].ID != "slot-1" {
		t.Fatalf("unexpected slots: %+v", slots)
	}
}
