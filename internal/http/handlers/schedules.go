package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"room-booking/internal/domain"
	"room-booking/internal/http/response"
)

type scheduleService interface {
	Create(ctx context.Context, roomID string, daysOfWeek []int, startTime, endTime string) (*domain.Schedule, error)
}

type ScheduleHandler struct {
	service scheduleService
}

func NewScheduleHandler(service scheduleService) *ScheduleHandler {
	return &ScheduleHandler{service: service}
}

type createScheduleRequest struct {
	DaysOfWeek []int  `json:"daysOfWeek"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}

func (h *ScheduleHandler) Create(w http.ResponseWriter, r *http.Request) {
	roomID := r.PathValue("roomId")

	var req createScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request")
		return
	}

	schedule, err := h.service.Create(
		r.Context(),
		roomID,
		req.DaysOfWeek,
		req.StartTime,
		req.EndTime,
	)
	if err != nil {
		response.WriteDomainError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusCreated, map[string]any{
		"schedule": schedule,
	})
}
