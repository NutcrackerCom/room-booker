package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"room-booking/internal/domain"
	"room-booking/internal/http/response"
)

type roomService interface {
	Create(ctx context.Context, name string, description *string, capacity *int) (*domain.Room, error)
	List(ctx context.Context) ([]domain.Room, error)
}

type RoomHandler struct {
	service roomService
}

func NewRoomHandler(service roomService) *RoomHandler {
	return &RoomHandler{service: service}
}

type createRoomRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Capacity    *int    `json:"capacity"`
}

func (h *RoomHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request")
		return
	}

	room, err := h.service.Create(r.Context(), req.Name, req.Description, req.Capacity)
	if err != nil {
		response.WriteDomainError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusCreated, map[string]any{
		"room": room,
	})
}

func (h *RoomHandler) List(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.service.List(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]any{
		"rooms": rooms,
	})
}
