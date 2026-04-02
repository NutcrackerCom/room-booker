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
	"room-booking/internal/http/middleware"
)

type fakeBookingServiceForHandler struct {
	createFn         func(ctx context.Context, slotID, userID string, createConferenceLink bool) (*domain.Booking, error)
	listMyUpcomingFn func(ctx context.Context, userID string) ([]domain.Booking, error)
	cancelFn         func(ctx context.Context, bookingID, userID string) (*domain.Booking, error)
	listAllFn        func(ctx context.Context, page, pageSize int) ([]domain.Booking, int, error)
}

func (f *fakeBookingServiceForHandler) Create(ctx context.Context, slotID, userID string, createConferenceLink bool) (*domain.Booking, error) {
	return f.createFn(ctx, slotID, userID, createConferenceLink)
}
func (f *fakeBookingServiceForHandler) ListMyUpcoming(ctx context.Context, userID string) ([]domain.Booking, error) {
	return f.listMyUpcomingFn(ctx, userID)
}
func (f *fakeBookingServiceForHandler) Cancel(ctx context.Context, bookingID, userID string) (*domain.Booking, error) {
	return f.cancelFn(ctx, bookingID, userID)
}
func (f *fakeBookingServiceForHandler) ListAll(ctx context.Context, page, pageSize int) ([]domain.Booking, int, error) {
	return f.listAllFn(ctx, page, pageSize)
}

func TestBookingHandler_Create_Success(t *testing.T) {
	handler := NewBookingHandler(&fakeBookingServiceForHandler{
		createFn: func(ctx context.Context, slotID, userID string, createConferenceLink bool) (*domain.Booking, error) {
			return &domain.Booking{
				ID:     "booking-1",
				SlotID: slotID,
				UserID: userID,
				Status: "active",
			}, nil
		},
		listMyUpcomingFn: func(ctx context.Context, userID string) ([]domain.Booking, error) { return nil, nil },
		cancelFn:         func(ctx context.Context, bookingID, userID string) (*domain.Booking, error) { return nil, nil },
		listAllFn:        func(ctx context.Context, page, pageSize int) ([]domain.Booking, int, error) { return nil, 0, nil },
	})

	req := httptest.NewRequest(http.MethodPost, "/bookings/create", bytes.NewBufferString(`{"slotId":"slot-1"}`))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.Create(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), `"status":"active"`) {
		t.Fatalf("unexpected body: %s", rr.Body.String())
	}
}

func TestBookingHandler_My_Success(t *testing.T) {
	now := time.Now().UTC()

	handler := NewBookingHandler(&fakeBookingServiceForHandler{
		createFn: func(ctx context.Context, slotID, userID string, createConferenceLink bool) (*domain.Booking, error) {
			return nil, nil
		},
		listMyUpcomingFn: func(ctx context.Context, userID string) ([]domain.Booking, error) {
			return []domain.Booking{
				{
					ID:        "booking-1",
					SlotID:    "slot-1",
					UserID:    userID,
					Status:    "active",
					CreatedAt: &now,
				},
			}, nil
		},
		cancelFn:  func(ctx context.Context, bookingID, userID string) (*domain.Booking, error) { return nil, nil },
		listAllFn: func(ctx context.Context, page, pageSize int) ([]domain.Booking, int, error) { return nil, 0, nil },
	})

	req := httptest.NewRequest(http.MethodGet, "/bookings/my", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.My(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", rr.Code, rr.Body.String())
	}
}

func TestBookingHandler_Cancel_Success(t *testing.T) {
	handler := NewBookingHandler(&fakeBookingServiceForHandler{
		createFn:         func(ctx context.Context, slotID, userID string, createConferenceLink bool) (*domain.Booking, error) { return nil, nil },
		listMyUpcomingFn: func(ctx context.Context, userID string) ([]domain.Booking, error) { return nil, nil },
		cancelFn: func(ctx context.Context, bookingID, userID string) (*domain.Booking, error) {
			return &domain.Booking{
				ID:     bookingID,
				UserID: userID,
				Status: "cancelled",
			}, nil
		},
		listAllFn: func(ctx context.Context, page, pageSize int) ([]domain.Booking, int, error) { return nil, 0, nil },
	})

	req := httptest.NewRequest(http.MethodPost, "/bookings/booking-1/cancel", nil)
	req.SetPathValue("bookingId", "booking-1")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, "user-1")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.Cancel(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), `"status":"cancelled"`) {
		t.Fatalf("unexpected body: %s", rr.Body.String())
	}
}

func TestBookingHandler_List_Success(t *testing.T) {
	now := time.Now().UTC()

	handler := NewBookingHandler(&fakeBookingServiceForHandler{
		createFn:         func(ctx context.Context, slotID, userID string, createConferenceLink bool) (*domain.Booking, error) { return nil, nil },
		listMyUpcomingFn: func(ctx context.Context, userID string) ([]domain.Booking, error) { return nil, nil },
		cancelFn:         func(ctx context.Context, bookingID, userID string) (*domain.Booking, error) { return nil, nil },
		listAllFn: func(ctx context.Context, page, pageSize int) ([]domain.Booking, int, error) {
			return []domain.Booking{
				{
					ID:        "booking-1",
					SlotID:    "slot-1",
					UserID:    "user-1",
					Status:    "active",
					CreatedAt: &now,
				},
			}, 1, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/bookings/list?page=1&pageSize=1", nil)
	rr := httptest.NewRecorder()

	handler.List(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), `"pagination"`) {
		t.Fatalf("unexpected body: %s", rr.Body.String())
	}
}

func TestBookingHandler_List_InvalidPage(t *testing.T) {
	handler := NewBookingHandler(&fakeBookingServiceForHandler{
		createFn:         func(ctx context.Context, slotID, userID string, createConferenceLink bool) (*domain.Booking, error) { return nil, nil },
		listMyUpcomingFn: func(ctx context.Context, userID string) ([]domain.Booking, error) { return nil, nil },
		cancelFn:         func(ctx context.Context, bookingID, userID string) (*domain.Booking, error) { return nil, nil },
		listAllFn:        func(ctx context.Context, page, pageSize int) ([]domain.Booking, int, error) { return nil, 0, nil },
	})

	req := httptest.NewRequest(http.MethodGet, "/bookings/list?page=bad", nil)
	rr := httptest.NewRecorder()

	handler.List(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", rr.Code, rr.Body.String())
	}
}
