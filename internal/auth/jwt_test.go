package auth

import "testing"

func TestCreateAndParseAdminToken(t *testing.T) {
	manager := NewJWTManager("test-secret")

	token, err := manager.CreateToken(AdminID, "admin")
	if err != nil {
		t.Fatalf("create token: %v", err)
	}

	claims, err := manager.ParseToken(token)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}

	if claims.Role != "admin" {
		t.Fatalf("expected role admin, got %q", claims.Role)
	}
	if claims.UserID != AdminID {
		t.Fatalf("expected admin id %q, got %q", AdminID, claims.UserID)
	}
}

func TestCreateAndParseUserToken(t *testing.T) {
	manager := NewJWTManager("test-secret")

	token, err := manager.CreateToken(UserID, "user")
	if err != nil {
		t.Fatalf("create token: %v", err)
	}

	claims, err := manager.ParseToken(token)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}

	if claims.Role != "user" {
		t.Fatalf("expected role user, got %q", claims.Role)
	}
	if claims.UserID != UserID {
		t.Fatalf("expected user id %q, got %q", UserID, claims.UserID)
	}
}

func TestParseInvalidToken(t *testing.T) {
	manager := NewJWTManager("test-secret")

	_, err := manager.ParseToken("not-a-token")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}
