package handlers

import (
	"context"
	"net/http"

	"room-booking/internal/domain"
	"room-booking/internal/http/response"
)

type slotService interface {
	ListAvailable(ctx context.Context, roomID, date string) ([]domain.Slot, error)
}

type SlotHandler struct {
	service slotService
}

func NewSlotHandler(service slotService) *SlotHandler {
	return &SlotHandler{service: service}
}

func (h *SlotHandler) List(w http.ResponseWriter, r *http.Request) {
	roomID := r.PathValue("roomId")
	date := r.URL.Query().Get("date")

	slots, err := h.service.ListAvailable(r.Context(), roomID, date)
	if err != nil {
		response.WriteDomainError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]any{
		"slots": slots,
	})
}
