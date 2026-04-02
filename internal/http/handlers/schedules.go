package handlers

import (
	"encoding/json"
	"net/http"

	"room-booking/internal/http/response"
	"room-booking/internal/service"
)

type ScheduleHandler struct {
	service *service.ScheduleService
}

func NewScheduleHandler(service *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{service: service}
}

type createScheduleRequest struct {
	RoomID     string `json:"roomId"`
	DaysOfWeek []int  `json:"daysOfWeek"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
}

func (h *ScheduleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request")
		return
	}

	schedule, err := h.service.Create(
		r.Context(),
		req.RoomID,
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
