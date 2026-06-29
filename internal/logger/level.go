package logger

import (
	"fmt"
	"log/slog"
)

func parseLevel(level string) (slog.Level, error) {
	switch level {
	case "debug":
		return slog.LevelDebug, nil

	case "info":
		return slog.LevelInfo, nil

	case "warn":
		return slog.LevelWarn, nil

	case "error":
		return slog.LevelError, nil

	default:
		return 0, fmt.Errorf(
			"invalid log level %q (valid values: debug, info, warn, error)",
			level,
		)
	}
}
