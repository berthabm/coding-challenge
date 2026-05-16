package utils

import (
	"log/slog"
	"os"
	"strings"
)

// Logger is a thin alias so services/controllers depend on slog interface patterns.
type Logger = *slog.Logger

// NewLogger builds a structured JSON/text logger based on LOG_LEVEL.
func NewLogger(level string) Logger {
	var lvl slog.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
	return slog.New(handler)
}
