package config

import "testing"

func TestGetEnvReturnsFallbackWhenEmpty(t *testing.T) {
	t.Setenv("TEST_ENV_EMPTY", "")

	got := getEnv("TEST_ENV_EMPTY", "fallback")
	if got != "fallback" {
		t.Fatalf("expected fallback, got %q", got)
	}
}

func TestGetEnvReturnsValueWhenSet(t *testing.T) {
	t.Setenv("TEST_ENV_VALUE", "real-value")

	got := getEnv("TEST_ENV_VALUE", "fallback")
	if got != "real-value" {
		t.Fatalf("expected real-value, got %q", got)
	}
}

func TestLoadUsesEnvValues(t *testing.T) {
	t.Setenv("APP_PORT", "9090")
	t.Setenv("DATABASE_URL", "postgres://test")
	t.Setenv("JWT_SECRET", "secret")

	cfg := Load()

	if cfg.AppPort != "9090" {
		t.Fatalf("expected app port 9090, got %q", cfg.AppPort)
	}
	if cfg.DatabaseURL != "postgres://test" {
		t.Fatalf("expected database url postgres://test, got %q", cfg.DatabaseURL)
	}
	if cfg.JWTSecret != "secret" {
		t.Fatalf("expected jwt secret secret, got %q", cfg.JWTSecret)
	}
}
