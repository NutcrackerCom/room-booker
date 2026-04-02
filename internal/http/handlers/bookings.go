package handlers

import (
	"encoding/json"
	"net/http"

	"room-booking/internal/http/middleware"
	"room-booking/internal/http/response"
	"room-booking/internal/service"
)

type BookingHandler struct {
	service *service.BookingService
}

func NewBookingHandler(service *service.BookingService) *BookingHandler {
	return &BookingHandler{service: service}
}

type createBookingRequest struct {
	SlotID               string `json:"slotId"`
	CreateConferenceLink bool   `json:"createConferenceLink"`
}

func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request")
		return
	}

	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	booking, err := h.service.Create(r.Context(), req.SlotID, userID, req.CreateConferenceLink)
	if err != nil {
		response.WriteDomainError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusCreated, map[string]any{
		"booking": booking,
	})
}
