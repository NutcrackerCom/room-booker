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

	addr := ":" + cfg.AppPort
	fmt.Println("server started on", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}
}
