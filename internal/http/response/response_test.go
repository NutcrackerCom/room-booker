package response

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"room-booking/internal/domain"
)

func TestWriteJSON(t *testing.T) {
	rr := httptest.NewRecorder()

	WriteJSON(rr, http.StatusCreated, map[string]string{
		"status": "ok",
	})

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json, got %q", ct)
	}

	body := rr.Body.String()
	if !strings.Contains(body, `"status":"ok"`) {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestWriteError(t *testing.T) {
	rr := httptest.NewRecorder()

	WriteError(rr, http.StatusBadRequest, "INVALID_REQUEST", "invalid request")

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, `"code":"INVALID_REQUEST"`) {
		t.Fatalf("unexpected body: %s", body)
	}
	if !strings.Contains(body, `"message":"invalid request"`) {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestWriteDomainError_InvalidRequest(t *testing.T) {
	rr := httptest.NewRecorder()

	WriteDomainError(rr, domain.ErrInvalidRequest)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, `"code":"INVALID_REQUEST"`) {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestWriteDomainError_Forbidden(t *testing.T) {
	rr := httptest.NewRecorder()

	WriteDomainError(rr, domain.ErrForbidden)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, `"code":"FORBIDDEN"`) {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestWriteDomainError_Internal(t *testing.T) {
	rr := httptest.NewRecorder()

	WriteDomainError(rr, errors.New("boom"))

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, `"code":"INTERNAL_ERROR"`) {
		t.Fatalf("unexpected body: %s", body)
	}
}
