package logger

import (
	"log/slog"
	"os"
)

func init() {
	lvl := os.Getenv("LOG_LEVEL")
	format := os.Getenv("LOG_FORMAT")
	Configure(lvl, format)
}

// Configure sets up the global slog logger based on the given level and format.
// This is used by the CLI at startup; external consumers control slog via slog.SetDefault.
func Configure(level, format string) {
	slogLevel := parseLevel(level)

	opts := &slog.HandlerOptions{Level: slogLevel, AddSource: true}

	var handler slog.Handler
	switch format {
	case "json":
		handler = slog.NewJSONHandler(os.Stderr, opts)
	default:
		handler = slog.NewTextHandler(os.Stderr, opts)
	}

	slog.SetDefault(slog.New(handler))
}

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelError
	}
}
