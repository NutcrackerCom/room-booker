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
	"room-booking/internal/repository"
	"room-booking/internal/service"
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

	roomRepo := repository.NewRoomRepository(db)
	roomService := service.NewRoomService(roomRepo)
	roomHandler := handlers.NewRoomHandler(roomService)

	scheduleRepo := repository.NewScheduleRepository(db)

	slotRepo := repository.NewSlotRepository(db)
	slotService := service.NewSlotService(roomRepo, scheduleRepo, slotRepo)
	slotHandler := handlers.NewSlotHandler(slotService)

	scheduleService := service.NewScheduleService(roomRepo, scheduleRepo, slotRepo)
	scheduleHandler := handlers.NewScheduleHandler(scheduleService)

	bookingRepo := repository.NewBookingRepository(db)
	bookingService := service.NewBookingService(slotRepo, bookingRepo)
	bookingHandler := handlers.NewBookingHandler(bookingService)

	mux := http.NewServeMux()

	mux.HandleFunc("/_info", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/dummyLogin", authHandler.DummyLogin)

	protected := middleware.AuthRequired(jwtManager)

	mux.Handle("/rooms/list", protected(http.HandlerFunc(roomHandler.List)))
	mux.Handle("/rooms/create", protected(middleware.RequireRole("admin")(http.HandlerFunc(roomHandler.Create))))
	mux.Handle("/bookings/my", protected(middleware.RequireRole("user")(http.HandlerFunc(bookingHandler.My))))
	mux.Handle("/bookings/{bookingId}/cancel", protected(middleware.RequireRole("user")(http.HandlerFunc(bookingHandler.Cancel))))
	mux.Handle("/rooms/{roomId}/schedule/create", protected(middleware.RequireRole("admin")(http.HandlerFunc(scheduleHandler.Create))))
	mux.Handle("/rooms/{roomId}/slots/list", protected(http.HandlerFunc(slotHandler.List)))

	mux.Handle("/bookings/create", protected(middleware.RequireRole("user")(http.HandlerFunc(bookingHandler.Create))))

	mux.Handle("/bookings/list", protected(middleware.RequireRole("admin")(http.HandlerFunc(bookingHandler.List))))

	addr := ":" + cfg.AppPort
	fmt.Println("server started on", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}
}
