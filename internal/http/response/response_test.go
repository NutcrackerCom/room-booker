package response

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"room-booking/internal/domain"
)

func TestWriteJSON(t *testing.T) {
	rr := httptest.NewRecorder()

	WriteJSON(rr, http.StatusCreated, map[string]string{
		"status": "ok",
	})

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json, got %q", ct)
	}

	body := rr.Body.String()
	if !strings.Contains(body, `"status":"ok"`) {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestWriteError(t *testing.T) {
	rr := httptest.NewRecorder()

	WriteError(rr, http.StatusBadRequest, "INVALID_REQUEST", "invalid request")

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, `"code":"INVALID_REQUEST"`) {
		t.Fatalf("unexpected body: %s", body)
	}
	if !strings.Contains(body, `"message":"invalid request"`) {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestWriteDomainError_InvalidRequest(t *testing.T) {
	rr := httptest.NewRecorder()

	WriteDomainError(rr, domain.ErrInvalidRequest)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, `"code":"INVALID_REQUEST"`) {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestWriteDomainError_Forbidden(t *testing.T) {
	rr := httptest.NewRecorder()

	WriteDomainError(rr, domain.ErrForbidden)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, `"code":"FORBIDDEN"`) {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestWriteDomainError_Internal(t *testing.T) {
	rr := httptest.NewRecorder()

	WriteDomainError(rr, errors.New("boom"))

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, `"code":"INTERNAL_ERROR"`) {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestWriteDomainError_Unauthorized(t *testing.T) {
	rr := httptest.NewRecorder()

	WriteDomainError(rr, domain.ErrUnauthorized)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), `"code":"UNAUTHORIZED"`) {
		t.Fatalf("unexpected body: %s", rr.Body.String())
	}
}

func TestWriteDomainError_RoomNotFound(t *testing.T) {
	rr := httptest.NewRecorder()

	WriteDomainError(rr, domain.ErrRoomNotFound)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), `"code":"ROOM_NOT_FOUND"`) {
		t.Fatalf("unexpected body: %s", rr.Body.String())
	}
}

func TestWriteDomainError_SlotNotFound(t *testing.T) {
	rr := httptest.NewRecorder()

	WriteDomainError(rr, domain.ErrSlotNotFound)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), `"code":"SLOT_NOT_FOUND"`) {
		t.Fatalf("unexpected body: %s", rr.Body.String())
	}
}

func TestWriteDomainError_SlotBooked(t *testing.T) {
	rr := httptest.NewRecorder()

	WriteDomainError(rr, domain.ErrSlotBooked)

	if rr.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), `"code":"SLOT_ALREADY_BOOKED"`) {
		t.Fatalf("unexpected body: %s", rr.Body.String())
	}
}

func TestWriteDomainError_BookingNotFound(t *testing.T) {
	rr := httptest.NewRecorder()

	WriteDomainError(rr, domain.ErrBookingNotFound)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), `"code":"BOOKING_NOT_FOUND"`) {
		t.Fatalf("unexpected body: %s", rr.Body.String())
	}
}

func TestWriteDomainError_ScheduleExists(t *testing.T) {
	rr := httptest.NewRecorder()

	WriteDomainError(rr, domain.ErrScheduleExists)

	if rr.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, rr.Code)
	}
	if !strings.Contains(rr.Body.String(), `"code":"SCHEDULE_EXISTS"`) {
		t.Fatalf("unexpected body: %s", rr.Body.String())
	}
}
