package service

import (
	"context"
	"testing"

	"room-booking/internal/domain"
)

type fakeRoomRepo struct {
	createFn func(ctx context.Context, name string, description *string, capacity *int) (*domain.Room, error)
	listFn   func(ctx context.Context) ([]domain.Room, error)
}

func (f *fakeRoomRepo) Create(ctx context.Context, name string, description *string, capacity *int) (*domain.Room, error) {
	return f.createFn(ctx, name, description, capacity)
}

func (f *fakeRoomRepo) List(ctx context.Context) ([]domain.Room, error) {
	return f.listFn(ctx)
}

func TestRoomServiceCreate_Success(t *testing.T) {
	repo := &fakeRoomRepo{
		createFn: func(ctx context.Context, name string, description *string, capacity *int) (*domain.Room, error) {
			return &domain.Room{Name: name, Description: description, Capacity: capacity}, nil
		},
		listFn: func(ctx context.Context) ([]domain.Room, error) { return nil, nil },
	}

	svc := NewRoomService(repo)

	desc := "Small room"
	capacity := 4
	room, err := svc.Create(context.Background(), "  Blue Room  ", &desc, &capacity)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if room.Name != "Blue Room" {
		t.Fatalf("expected trimmed name, got %q", room.Name)
	}
}

func TestRoomServiceCreate_EmptyName(t *testing.T) {
	repo := &fakeRoomRepo{
		createFn: func(ctx context.Context, name string, description *string, capacity *int) (*domain.Room, error) {
			t.Fatal("repo.Create should not be called")
			return nil, nil
		},
		listFn: func(ctx context.Context) ([]domain.Room, error) { return nil, nil },
	}

	svc := NewRoomService(repo)

	_, err := svc.Create(context.Background(), "   ", nil, nil)
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got %v", err)
	}
}

func TestRoomServiceCreate_NegativeCapacity(t *testing.T) {
	repo := &fakeRoomRepo{
		createFn: func(ctx context.Context, name string, description *string, capacity *int) (*domain.Room, error) {
			t.Fatal("repo.Create should not be called")
			return nil, nil
		},
		listFn: func(ctx context.Context) ([]domain.Room, error) { return nil, nil },
	}

	svc := NewRoomService(repo)

	capacity := -1
	_, err := svc.Create(context.Background(), "Blue Room", nil, &capacity)
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got %v", err)
	}
}

func TestRoomServiceList(t *testing.T) {
	repo := &fakeRoomRepo{
		createFn: func(ctx context.Context, name string, description *string, capacity *int) (*domain.Room, error) {
			return nil, nil
		},
		listFn: func(ctx context.Context) ([]domain.Room, error) {
			return []domain.Room{{ID: "room-1", Name: "Blue Room"}}, nil
		},
	}

	svc := NewRoomService(repo)

	rooms, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(rooms) != 1 {
		t.Fatalf("expected 1 room, got %d", len(rooms))
	}
}
