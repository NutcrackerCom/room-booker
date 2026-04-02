package config

import "os"

type Config struct {
	AppPort     string
	DatabaseURL string
	JWTSecret   string
}

func Load() Config {
	return Config{
		AppPort:     getEnv("APP_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@db:5432/room_booking?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "super-secret-key"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
