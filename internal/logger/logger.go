package logger

import (
	"log/slog"

	"github.com/aliasadiwastaken/hotel-booking-system/internal/config"
)

func New(cfg config.LoggerConfig) (*slog.Logger, error) {
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	handlerOptions := &slog.HandlerOptions{
		Level:     level,
		AddSource: cfg.AddSource,
	}

	handler, err := createHandler(cfg.Format, handlerOptions)
	if err != nil {
		return nil, err
	}

	return slog.New(handler), nil
}
