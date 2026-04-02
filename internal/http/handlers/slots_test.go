package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"room-booking/internal/domain"
)

type fakeSlotServiceForHandler struct {
	listAvailableFn func(ctx context.Context, roomID, date string) ([]domain.Slot, error)
}

func (f *fakeSlotServiceForHandler) ListAvailable(ctx context.Context, roomID, date string) ([]domain.Slot, error) {
	return f.listAvailableFn(ctx, roomID, date)
}

func TestSlotHandler_List_Success(t *testing.T) {
	handler := NewSlotHandler(&fakeSlotServiceForHandler{
		listAvailableFn: func(ctx context.Context, roomID, date string) ([]domain.Slot, error) {
			return []domain.Slot{
				{
					ID:     "slot-1",
					RoomID: roomID,
					Start:  time.Date(2026, 4, 3, 9, 0, 0, 0, time.UTC),
					End:    time.Date(2026, 4, 3, 9, 30, 0, 0, time.UTC),
				},
			}, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/rooms/room-1/slots/list?date=2026-04-03", nil)
	req.SetPathValue("roomId", "room-1")
	rr := httptest.NewRecorder()

	handler.List(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), `"slots":[`) {
		t.Fatalf("unexpected body: %s", rr.Body.String())
	}
}

func TestSlotHandler_List_DomainError(t *testing.T) {
	handler := NewSlotHandler(&fakeSlotServiceForHandler{
		listAvailableFn: func(ctx context.Context, roomID, date string) ([]domain.Slot, error) {
			return nil, domain.ErrInvalidRequest
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/rooms/room-1/slots/list", nil)
	req.SetPathValue("roomId", "room-1")
	rr := httptest.NewRecorder()

	handler.List(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", rr.Code, rr.Body.String())
	}
}
