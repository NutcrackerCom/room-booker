package handlers

import (
	"encoding/json"
	"net/http"

	"room-booking/internal/auth"
	"room-booking/internal/http/response"
)

type AuthHandler struct {
	jwt *auth.JWTManager
}

func NewAuthHandler(jwt *auth.JWTManager) *AuthHandler {
	return &AuthHandler{jwt: jwt}
}

type dummyLoginRequest struct {
	Role string `json:"role"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	var req dummyLoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request")
		return
	}

	if req.Role != "admin" && req.Role != "user" {
		response.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request")
		return
	}

	token, err := h.jwt.CreateToken(req.Role)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	response.WriteJSON(w, http.StatusOK, tokenResponse{Token: token})
}
