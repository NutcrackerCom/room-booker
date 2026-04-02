package response

import (
	"errors"
	"net/http"

	"room-booking/internal/domain"
)

type errorBody struct {
	Error errorItem `json:"error"`
}

type errorItem struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WriteError(w http.ResponseWriter, status int, code, message string) {
	WriteJSON(w, status, errorBody{
		Error: errorItem{
			Code:    code,
			Message: message,
		},
	})
}

func WriteDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidRequest):
		WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request")
	case errors.Is(err, domain.ErrUnauthorized):
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
	case errors.Is(err, domain.ErrForbidden):
		WriteError(w, http.StatusForbidden, "FORBIDDEN", "forbidden")
	case errors.Is(err, domain.ErrRoomNotFound):
		WriteError(w, http.StatusNotFound, "ROOM_NOT_FOUND", "room not found")
	case errors.Is(err, domain.ErrSlotNotFound):
		WriteError(w, http.StatusNotFound, "SLOT_NOT_FOUND", "slot not found")
	case errors.Is(err, domain.ErrSlotBooked):
		WriteError(w, http.StatusConflict, "SLOT_ALREADY_BOOKED", "slot is already booked")
	case errors.Is(err, domain.ErrBookingNotFound):
		WriteError(w, http.StatusNotFound, "BOOKING_NOT_FOUND", "booking not found")
	case errors.Is(err, domain.ErrScheduleExists):
		WriteError(w, http.StatusConflict, "SCHEDULE_EXISTS", "schedule for this room already exists and cannot be changed")
	default:
		WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
	}
}
