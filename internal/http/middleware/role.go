package middleware

import (
	"net/http"

	"room-booking/internal/http/response"
)

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			currentRole, ok := r.Context().Value(RoleKey).(string)
			if !ok || currentRole == "" {
				response.WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
				return
			}

			if currentRole != role {
				response.WriteError(w, http.StatusForbidden, "FORBIDDEN", "forbidden")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
