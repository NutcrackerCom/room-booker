package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"room-booking/internal/domain"
)

type fakeRoomService struct {
	createFn func(ctx context.Context, name string, description *string, capacity *int) (*domain.Room, error)
	listFn   func(ctx context.Context) ([]domain.Room, error)
}

func (f *fakeRoomService) Create(ctx context.Context, name string, description *string, capacity *int) (*domain.Room, error) {
	return f.createFn(ctx, name, description, capacity)
}

func (f *fakeRoomService) List(ctx context.Context) ([]domain.Room, error) {
	return f.listFn(ctx)
}

func TestRoomHandler_Create_Success(t *testing.T) {
	now := time.Now().UTC()

	handler := NewRoomHandler(&fakeRoomService{
		createFn: func(ctx context.Context, name string, description *string, capacity *int) (*domain.Room, error) {
			return &domain.Room{
				ID:          "room-1",
				Name:        name,
				Description: description,
				Capacity:    capacity,
				CreatedAt:   &now,
			}, nil
		},
		listFn: func(ctx context.Context) ([]domain.Room, error) { return nil, nil },
	})

	req := httptest.NewRequest(http.MethodPost, "/rooms/create", bytes.NewBufferString(`{"name":"Blue Room","description":"Small room","capacity":4}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Create(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body=%s", rr.Code, rr.Body.String())
	}

	body := rr.Body.String()
	if !strings.Contains(body, `"name":"Blue Room"`) {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestRoomHandler_Create_InvalidJSON(t *testing.T) {
	handler := NewRoomHandler(&fakeRoomService{
		createFn: func(ctx context.Context, name string, description *string, capacity *int) (*domain.Room, error) {
			return nil, nil
		},
		listFn: func(ctx context.Context) ([]domain.Room, error) { return nil, nil },
	})

	req := httptest.NewRequest(http.MethodPost, "/rooms/create", bytes.NewBufferString(`{bad json}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Create(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", rr.Code, rr.Body.String())
	}
}

func TestRoomHandler_Create_DomainError(t *testing.T) {
	handler := NewRoomHandler(&fakeRoomService{
		createFn: func(ctx context.Context, name string, description *string, capacity *int) (*domain.Room, error) {
			return nil, domain.ErrInvalidRequest
		},
		listFn: func(ctx context.Context) ([]domain.Room, error) { return nil, nil },
	})

	req := httptest.NewRequest(http.MethodPost, "/rooms/create", bytes.NewBufferString(`{"name":""}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Create(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", rr.Code, rr.Body.String())
	}
}

func TestRoomHandler_List_Success(t *testing.T) {
	now := time.Now().UTC()

	handler := NewRoomHandler(&fakeRoomService{
		createFn: func(ctx context.Context, name string, description *string, capacity *int) (*domain.Room, error) {
			return nil, nil
		},
		listFn: func(ctx context.Context) ([]domain.Room, error) {
			return []domain.Room{
				{
					ID:        "room-1",
					Name:      "Blue Room",
					CreatedAt: &now,
				},
			}, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/rooms/list", nil)
	rr := httptest.NewRecorder()

	handler.List(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", rr.Code, rr.Body.String())
	}

	body := rr.Body.String()
	if !strings.Contains(body, `"rooms":[`) {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestRoomHandler_List_InternalError(t *testing.T) {
	handler := NewRoomHandler(&fakeRoomService{
		createFn: func(ctx context.Context, name string, description *string, capacity *int) (*domain.Room, error) {
			return nil, nil
		},
		listFn: func(ctx context.Context) ([]domain.Room, error) {
			return nil, context.DeadlineExceeded
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/rooms/list", nil)
	rr := httptest.NewRecorder()

	handler.List(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d, body=%s", rr.Code, rr.Body.String())
	}
}
