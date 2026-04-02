package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"room-booking/internal/domain"
)

type fakeScheduleService struct {
	createFn func(ctx context.Context, roomID string, daysOfWeek []int, startTime, endTime string) (*domain.Schedule, error)
}

func (f *fakeScheduleService) Create(ctx context.Context, roomID string, daysOfWeek []int, startTime, endTime string) (*domain.Schedule, error) {
	return f.createFn(ctx, roomID, daysOfWeek, startTime, endTime)
}

func TestScheduleHandler_Create_Success(t *testing.T) {
	handler := NewScheduleHandler(&fakeScheduleService{
		createFn: func(ctx context.Context, roomID string, daysOfWeek []int, startTime, endTime string) (*domain.Schedule, error) {
			return &domain.Schedule{
				ID:         "schedule-1",
				RoomID:     roomID,
				DaysOfWeek: daysOfWeek,
				StartTime:  startTime,
				EndTime:    endTime,
			}, nil
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/rooms/room-1/schedule/create", bytes.NewBufferString(`{"daysOfWeek":[1,2,3],"startTime":"09:00","endTime":"18:00"}`))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("roomId", "room-1")
	rr := httptest.NewRecorder()

	handler.Create(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), `"roomId":"room-1"`) {
		t.Fatalf("unexpected body: %s", rr.Body.String())
	}
}

func TestScheduleHandler_Create_InvalidJSON(t *testing.T) {
	handler := NewScheduleHandler(&fakeScheduleService{
		createFn: func(ctx context.Context, roomID string, daysOfWeek []int, startTime, endTime string) (*domain.Schedule, error) {
			return nil, nil
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/rooms/room-1/schedule/create", bytes.NewBufferString(`{bad json}`))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("roomId", "room-1")
	rr := httptest.NewRecorder()

	handler.Create(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", rr.Code, rr.Body.String())
	}
}

func TestScheduleHandler_Create_DomainError(t *testing.T) {
	handler := NewScheduleHandler(&fakeScheduleService{
		createFn: func(ctx context.Context, roomID string, daysOfWeek []int, startTime, endTime string) (*domain.Schedule, error) {
			return nil, domain.ErrScheduleExists
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/rooms/room-1/schedule/create", bytes.NewBufferString(`{"daysOfWeek":[1],"startTime":"09:00","endTime":"18:00"}`))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("roomId", "room-1")
	rr := httptest.NewRecorder()

	handler.Create(rr, req)

	if rr.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d, body=%s", rr.Code, rr.Body.String())
	}
}
