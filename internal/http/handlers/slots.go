package handlers

import (
	"net/http"

	"room-booking/internal/http/response"
	"room-booking/internal/service"
)

type SlotHandler struct {
	service *service.SlotService
}

func NewSlotHandler(service *service.SlotService) *SlotHandler {
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
