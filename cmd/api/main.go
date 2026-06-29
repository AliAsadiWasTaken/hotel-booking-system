package main

import (
	"log/slog"
	"os"

	"github.com/aliasadiwastaken/hotel-booking-system/internal/booking"
	"github.com/aliasadiwastaken/hotel-booking-system/internal/config"
	"github.com/aliasadiwastaken/hotel-booking-system/internal/database"
	"github.com/aliasadiwastaken/hotel-booking-system/internal/hotel"
	"github.com/aliasadiwastaken/hotel-booking-system/internal/logger"
	"github.com/aliasadiwastaken/hotel-booking-system/internal/room"
	"github.com/aliasadiwastaken/hotel-booking-system/internal/router"
	"github.com/aliasadiwastaken/hotel-booking-system/internal/server"
	"github.com/aliasadiwastaken/hotel-booking-system/internal/user"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	appLogger, err := logger.New(cfg.Logger)
	if err != nil {
		slog.Error("failed to initialize logger", "error", err)
		os.Exit(1)
	}

	db, err := database.Connect(cfg.Database)
	if err != nil {
		appLogger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	appLogger.Info("database connected")

	// Repositories
	hotelRepo := hotel.NewRepository(db)
	roomRepo := room.NewRepository(db)
	userRepo := user.NewRepository(db)
	bookingRepo := booking.NewRepository(db)

	// Services
	hotelService := hotel.NewService(hotelRepo)
	roomService := room.NewService(roomRepo)
	userService := user.NewService(userRepo)
	bookingService := booking.NewService(db, bookingRepo, roomRepo)

	// Handlers
	hotelHandler := hotel.NewHandler(hotelService)
	roomHandler := room.NewHandler(roomService)
	userHandler := user.NewHandler(userService)
	bookingHandler := booking.NewHandler(bookingService)

	mux := router.New(hotelHandler, roomHandler, userHandler, bookingHandler)

	s := server.New(mux, cfg.HTTP.Address)

	appLogger.Info("starting server", "address", cfg.HTTP.Address)

	if err := s.ListenAndServe(); err != nil {
		appLogger.Error("server stopped", "error", err)
	}
}
