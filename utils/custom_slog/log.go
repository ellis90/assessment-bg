package custom_slog

import (
	"golang.org/x/exp/slog"
	"os"
)

func StructuredLog() *slog.Logger {
	jsonHandler := slog.NewJSONHandler(os.Stdout).
		WithAttrs([]slog.Attr{slog.String("app-version", "v0.0.1-beta")})
	logger := slog.New(jsonHandler)
	return logger
}
