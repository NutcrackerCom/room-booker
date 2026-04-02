package service

import (
	"context"
	"testing"
	"time"

	"room-booking/internal/domain"
)

type fakeBookingSlotRepo struct {
	getByIDFn func(ctx context.Context, slotID string) (*domain.Slot, error)
}

func (f *fakeBookingSlotRepo) GetByID(ctx context.Context, slotID string) (*domain.Slot, error) {
	return f.getByIDFn(ctx, slotID)
}

type fakeBookingRepo struct {
	createFn         func(ctx context.Context, slotID, userID string, conferenceLink *string) (*domain.Booking, error)
	listMyUpcomingFn func(ctx context.Context, userID string) ([]domain.Booking, error)
	getByIDFn        func(ctx context.Context, bookingID string) (*domain.Booking, error)
	cancelFn         func(ctx context.Context, bookingID string) (*domain.Booking, error)
	listAllFn        func(ctx context.Context, page, pageSize int) ([]domain.Booking, int, error)
}

func (f *fakeBookingRepo) Create(ctx context.Context, slotID, userID string, conferenceLink *string) (*domain.Booking, error) {
	return f.createFn(ctx, slotID, userID, conferenceLink)
}
func (f *fakeBookingRepo) ListMyUpcoming(ctx context.Context, userID string) ([]domain.Booking, error) {
	return f.listMyUpcomingFn(ctx, userID)
}
func (f *fakeBookingRepo) GetByID(ctx context.Context, bookingID string) (*domain.Booking, error) {
	return f.getByIDFn(ctx, bookingID)
}
func (f *fakeBookingRepo) Cancel(ctx context.Context, bookingID string) (*domain.Booking, error) {
	return f.cancelFn(ctx, bookingID)
}
func (f *fakeBookingRepo) ListAll(ctx context.Context, page, pageSize int) ([]domain.Booking, int, error) {
	return f.listAllFn(ctx, page, pageSize)
}

func TestBookingServiceCreate_Success(t *testing.T) {
	slotStart := time.Now().UTC().Add(24 * time.Hour)

	svc := NewBookingService(
		&fakeBookingSlotRepo{
			getByIDFn: func(ctx context.Context, slotID string) (*domain.Slot, error) {
				return &domain.Slot{ID: slotID, Start: slotStart}, nil
			},
		},
		&fakeBookingRepo{
			createFn: func(ctx context.Context, slotID, userID string, conferenceLink *string) (*domain.Booking, error) {
				return &domain.Booking{
					ID:     "booking-1",
					SlotID: slotID,
					UserID: userID,
					Status: "active",
				}, nil
			},
			listMyUpcomingFn: func(ctx context.Context, userID string) ([]domain.Booking, error) { return nil, nil },
			getByIDFn:        func(ctx context.Context, bookingID string) (*domain.Booking, error) { return nil, nil },
			cancelFn:         func(ctx context.Context, bookingID string) (*domain.Booking, error) { return nil, nil },
			listAllFn:        func(ctx context.Context, page, pageSize int) ([]domain.Booking, int, error) { return nil, 0, nil },
		},
	)

	booking, err := svc.Create(context.Background(), "slot-1", "user-1", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if booking.Status != "active" {
		t.Fatalf("expected active booking, got %+v", booking)
	}
}

func TestBookingServiceCreate_PastSlot(t *testing.T) {
	slotStart := time.Now().UTC().Add(-1 * time.Hour)

	svc := NewBookingService(
		&fakeBookingSlotRepo{
			getByIDFn: func(ctx context.Context, slotID string) (*domain.Slot, error) {
				return &domain.Slot{ID: slotID, Start: slotStart}, nil
			},
		},
		&fakeBookingRepo{
			createFn:         func(ctx context.Context, slotID, userID string, conferenceLink *string) (*domain.Booking, error) { return nil, nil },
			listMyUpcomingFn: func(ctx context.Context, userID string) ([]domain.Booking, error) { return nil, nil },
			getByIDFn:        func(ctx context.Context, bookingID string) (*domain.Booking, error) { return nil, nil },
			cancelFn:         func(ctx context.Context, bookingID string) (*domain.Booking, error) { return nil, nil },
			listAllFn:        func(ctx context.Context, page, pageSize int) ([]domain.Booking, int, error) { return nil, 0, nil },
		},
	)

	_, err := svc.Create(context.Background(), "slot-1", "user-1", false)
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got %v", err)
	}
}

func TestBookingServiceCreate_SlotBooked(t *testing.T) {
	slotStart := time.Now().UTC().Add(24 * time.Hour)

	svc := NewBookingService(
		&fakeBookingSlotRepo{
			getByIDFn: func(ctx context.Context, slotID string) (*domain.Slot, error) {
				return &domain.Slot{ID: slotID, Start: slotStart}, nil
			},
		},
		&fakeBookingRepo{
			createFn: func(ctx context.Context, slotID, userID string, conferenceLink *string) (*domain.Booking, error) {
				return nil, domain.ErrSlotBooked
			},
			listMyUpcomingFn: func(ctx context.Context, userID string) ([]domain.Booking, error) { return nil, nil },
			getByIDFn:        func(ctx context.Context, bookingID string) (*domain.Booking, error) { return nil, nil },
			cancelFn:         func(ctx context.Context, bookingID string) (*domain.Booking, error) { return nil, nil },
			listAllFn:        func(ctx context.Context, page, pageSize int) ([]domain.Booking, int, error) { return nil, 0, nil },
		},
	)

	_, err := svc.Create(context.Background(), "slot-1", "user-1", false)
	if err != domain.ErrSlotBooked {
		t.Fatalf("expected ErrSlotBooked, got %v", err)
	}
}

func TestBookingServiceListMyUpcoming_Unauthorized(t *testing.T) {
	svc := NewBookingService(
		&fakeBookingSlotRepo{getByIDFn: func(ctx context.Context, slotID string) (*domain.Slot, error) { return nil, nil }},
		&fakeBookingRepo{
			createFn:         func(ctx context.Context, slotID, userID string, conferenceLink *string) (*domain.Booking, error) { return nil, nil },
			listMyUpcomingFn: func(ctx context.Context, userID string) ([]domain.Booking, error) { return nil, nil },
			getByIDFn:        func(ctx context.Context, bookingID string) (*domain.Booking, error) { return nil, nil },
			cancelFn:         func(ctx context.Context, bookingID string) (*domain.Booking, error) { return nil, nil },
			listAllFn:        func(ctx context.Context, page, pageSize int) ([]domain.Booking, int, error) { return nil, 0, nil },
		},
	)

	_, err := svc.ListMyUpcoming(context.Background(), "")
	if err != domain.ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestBookingServiceCancel_AlreadyCancelled(t *testing.T) {
	svc := NewBookingService(
		&fakeBookingSlotRepo{getByIDFn: func(ctx context.Context, slotID string) (*domain.Slot, error) { return nil, nil }},
		&fakeBookingRepo{
			createFn:         func(ctx context.Context, slotID, userID string, conferenceLink *string) (*domain.Booking, error) { return nil, nil },
			listMyUpcomingFn: func(ctx context.Context, userID string) ([]domain.Booking, error) { return nil, nil },
			getByIDFn: func(ctx context.Context, bookingID string) (*domain.Booking, error) {
				return &domain.Booking{ID: bookingID, UserID: "user-1", Status: "cancelled"}, nil
			},
			cancelFn:  func(ctx context.Context, bookingID string) (*domain.Booking, error) { return nil, nil },
			listAllFn: func(ctx context.Context, page, pageSize int) ([]domain.Booking, int, error) { return nil, 0, nil },
		},
	)

	booking, err := svc.Cancel(context.Background(), "booking-1", "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if booking.Status != "cancelled" {
		t.Fatalf("expected cancelled, got %+v", booking)
	}
}

func TestBookingServiceCancel_Forbidden(t *testing.T) {
	svc := NewBookingService(
		&fakeBookingSlotRepo{getByIDFn: func(ctx context.Context, slotID string) (*domain.Slot, error) { return nil, nil }},
		&fakeBookingRepo{
			createFn:         func(ctx context.Context, slotID, userID string, conferenceLink *string) (*domain.Booking, error) { return nil, nil },
			listMyUpcomingFn: func(ctx context.Context, userID string) ([]domain.Booking, error) { return nil, nil },
			getByIDFn: func(ctx context.Context, bookingID string) (*domain.Booking, error) {
				return &domain.Booking{ID: bookingID, UserID: "another-user", Status: "active"}, nil
			},
			cancelFn:  func(ctx context.Context, bookingID string) (*domain.Booking, error) { return nil, nil },
			listAllFn: func(ctx context.Context, page, pageSize int) ([]domain.Booking, int, error) { return nil, 0, nil },
		},
	)

	_, err := svc.Cancel(context.Background(), "booking-1", "user-1")
	if err != domain.ErrForbidden {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestBookingServiceListAll_InvalidPageSize(t *testing.T) {
	svc := NewBookingService(
		&fakeBookingSlotRepo{getByIDFn: func(ctx context.Context, slotID string) (*domain.Slot, error) { return nil, nil }},
		&fakeBookingRepo{
			createFn:         func(ctx context.Context, slotID, userID string, conferenceLink *string) (*domain.Booking, error) { return nil, nil },
			listMyUpcomingFn: func(ctx context.Context, userID string) ([]domain.Booking, error) { return nil, nil },
			getByIDFn:        func(ctx context.Context, bookingID string) (*domain.Booking, error) { return nil, nil },
			cancelFn:         func(ctx context.Context, bookingID string) (*domain.Booking, error) { return nil, nil },
			listAllFn:        func(ctx context.Context, page, pageSize int) ([]domain.Booking, int, error) { return nil, 0, nil },
		},
	)

	_, _, err := svc.ListAll(context.Background(), 1, 101)
	if err != domain.ErrInvalidRequest {
		t.Fatalf("expected ErrInvalidRequest, got %v", err)
	}
}
