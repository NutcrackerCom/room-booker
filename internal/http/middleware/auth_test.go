package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"room-booking/internal/auth"
)

func TestAuthRequired_WithoutHeader(t *testing.T) {
	manager := auth.NewJWTManager("test-secret")

	h := AuthRequired(manager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}

func TestAuthRequired_WithInvalidHeader(t *testing.T) {
	manager := auth.NewJWTManager("test-secret")

	h := AuthRequired(manager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bad token")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}

func TestAuthRequired_WithValidToken(t *testing.T) {
	manager := auth.NewJWTManager("test-secret")

	token, err := manager.CreateToken("user")
	if err != nil {
		t.Fatalf("create token: %v", err)
	}

	h := AuthRequired(manager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(UserIDKey).(string)
		role, _ := r.Context().Value(RoleKey).(string)

		if userID == "" {
			t.Fatal("expected userID in context")
		}
		if role != "user" {
			t.Fatalf("expected role user, got %q", role)
		}

		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}
