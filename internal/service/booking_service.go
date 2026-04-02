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
