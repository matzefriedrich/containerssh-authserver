package logging

import (
	"log/slog"
	"os"
	"strings"

	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"
)

// NewSlogLogger creates a new slog logger based on the provided application configuration.
func NewSlogLogger(settings *configuration.ApplicationConfiguration) *slog.Logger {
	var level slog.Level
	switch strings.ToLower(settings.LogLevel) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}
