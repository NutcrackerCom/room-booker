package domain

import "errors"

var (
	ErrInvalidRequest  = errors.New("invalid request")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrForbidden       = errors.New("forbidden")
	ErrRoomNotFound    = errors.New("room not found")
	ErrSlotNotFound    = errors.New("slot not found")
	ErrSlotBooked      = errors.New("slot already booked")
	ErrBookingNotFound = errors.New("booking not found")
	ErrScheduleExists  = errors.New("schedule exists")
)
