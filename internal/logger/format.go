package logger

import (
	"fmt"
	"log/slog"
	"os"
)

func createHandler(
	format string,
	opts *slog.HandlerOptions,
) (slog.Handler, error) {

	switch format {

	case "text":
		return slog.NewTextHandler(os.Stdout, opts), nil

	case "json":
		return slog.NewJSONHandler(os.Stdout, opts), nil

	default:
		return nil, fmt.Errorf(
			"invalid log format %q (valid values: text, json)",
			format,
		)
	}
}
