package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"room-booking/internal/auth"
)

func TestDummyLogin_Admin(t *testing.T) {
	handler := NewAuthHandler(auth.NewJWTManager("test-secret"))

	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBufferString(`{"role":"admin"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.DummyLogin(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", rr.Code, rr.Body.String())
	}

	body := rr.Body.String()
	if !strings.Contains(body, `"token":"`) {
		t.Fatalf("expected token in response, got %s", body)
	}
}

func TestDummyLogin_InvalidRole(t *testing.T) {
	handler := NewAuthHandler(auth.NewJWTManager("test-secret"))

	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBufferString(`{"role":"manager"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.DummyLogin(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", rr.Code, rr.Body.String())
	}

	body := rr.Body.String()
	if !strings.Contains(body, `"code":"INVALID_REQUEST"`) {
		t.Fatalf("expected INVALID_REQUEST, got %s", body)
	}
}

func TestDummyLogin_BadJSON(t *testing.T) {
	handler := NewAuthHandler(auth.NewJWTManager("test-secret"))

	req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewBufferString(`{bad json}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.DummyLogin(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body=%s", rr.Code, rr.Body.String())
	}
}
