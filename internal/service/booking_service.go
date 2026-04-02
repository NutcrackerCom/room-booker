package service

import (
	"context"
	"errors"
	"time"

	"room-booking/internal/domain"
	"room-booking/internal/repository"
)

type BookingService struct {
	slotRepo    *repository.SlotRepository
	bookingRepo *repository.BookingRepository
}

func NewBookingService(slotRepo *repository.SlotRepository, bookingRepo *repository.BookingRepository) *BookingService {
	return &BookingService{
		slotRepo:    slotRepo,
		bookingRepo: bookingRepo,
	}
}

func (s *BookingService) Create(ctx context.Context, slotID, userID string, createConferenceLink bool) (*domain.Booking, error) {
	if slotID == "" || userID == "" {
		return nil, domain.ErrInvalidRequest
	}

	slot, err := s.slotRepo.GetByID(ctx, slotID)
	if err != nil {
		return nil, domain.ErrSlotNotFound
	}

	if slot.Start.Before(time.Now().UTC()) {
		return nil, domain.ErrInvalidRequest
	}

	var conferenceLink *string
	if createConferenceLink {
		link := "https://conference.local/mock-link"
		conferenceLink = &link
	}

	booking, err := s.bookingRepo.Create(ctx, slotID, userID, conferenceLink)
	if err != nil {
		if errors.Is(err, domain.ErrSlotBooked) {
			return nil, domain.ErrSlotBooked
		}
		return nil, err
	}

	return booking, nil
}

func (s *BookingService) ListMyUpcoming(ctx context.Context, userID string) ([]domain.Booking, error) {
	if userID == "" {
		return nil, domain.ErrUnauthorized
	}

	return s.bookingRepo.ListMyUpcoming(ctx, userID)
}

func (s *BookingService) Cancel(ctx context.Context, bookingID, userID string) (*domain.Booking, error) {
	if bookingID == "" || userID == "" {
		return nil, domain.ErrInvalidRequest
	}

	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return nil, domain.ErrBookingNotFound
	}

	if booking.UserID != userID {
		return nil, domain.ErrForbidden
	}

	if booking.Status == "cancelled" {
		return booking, nil
	}

	return s.bookingRepo.Cancel(ctx, bookingID)
}

func (s *BookingService) ListAll(ctx context.Context, page, pageSize int) ([]domain.Booking, int, error) {
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}

	if page < 1 || pageSize < 1 || pageSize > 100 {
		return nil, 0, domain.ErrInvalidRequest
	}

	return s.bookingRepo.ListAll(ctx, page, pageSize)
}
