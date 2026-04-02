package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"room-booking/internal/app"
	"room-booking/internal/auth"
	"room-booking/internal/config"
	"room-booking/internal/http/handlers"
	"room-booking/internal/http/middleware"
	"room-booking/internal/http/response"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := app.NewDB(ctx, cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		panic(err)
	}

	jwtManager := auth.NewJWTManager(cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(jwtManager)

	mux := http.NewServeMux()

	mux.HandleFunc("/_info", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/dummyLogin", authHandler.DummyLogin)

	protected := http.NewServeMux()
	protected.HandleFunc("/protected", func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(middleware.UserIDKey).(string)
		role, _ := r.Context().Value(middleware.RoleKey).(string)

		response.WriteJSON(w, http.StatusOK, map[string]string{
			"userId": userID,
			"role":   role,
		})
	})

	adminOnly := http.NewServeMux()
	adminOnly.HandleFunc("/admin-only", func(w http.ResponseWriter, r *http.Request) {
		response.WriteJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	mux.Handle("/protected", middleware.AuthRequired(jwtManager)(protected))
	mux.Handle("/admin-only", middleware.AuthRequired(jwtManager)(middleware.RequireRole("admin")(adminOnly)))

	addr := ":" + cfg.AppPort
	fmt.Println("server started on", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}
}
