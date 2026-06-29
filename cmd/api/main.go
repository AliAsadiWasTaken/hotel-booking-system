package main

import (
	"log/slog"
	"os"

	"github.com/aliasadiwastaken/hotel-booking-system/internal/config"
	"github.com/aliasadiwastaken/hotel-booking-system/internal/database"
	"github.com/aliasadiwastaken/hotel-booking-system/internal/logger"
	"github.com/aliasadiwastaken/hotel-booking-system/internal/router"
	"github.com/aliasadiwastaken/hotel-booking-system/internal/server"
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

	mux := router.New()

	s := server.New(
		mux,
		cfg.HTTP.Address,
	)

	appLogger.Info(
		"starting server",
		"address", cfg.HTTP.Address,
	)

	err = s.ListenAndServe()
	if err != nil {
		appLogger.Error(
			"server stopped",
			"error", err,
		)
	}

}
