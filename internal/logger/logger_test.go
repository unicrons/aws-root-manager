package logger_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unicrons/aws-root-manager/internal/logger"
)

func TestConfigure_ValidLevels(t *testing.T) {
	tests := []struct {
		level     string
		slogLevel slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"error", slog.LevelError},
	}
	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			logger.Configure(tt.level, "text")
			assert.True(t, slog.Default().Enabled(context.Background(), tt.slogLevel),
				"Configure(%q) did not enable level %v", tt.level, tt.slogLevel)
		})
	}
}

func TestConfigure_ValidFormats(t *testing.T) {
	formats := []string{"text", "json"}
	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			logger.Configure("info", format)
		})
	}
}

func TestConfigure_UnknownLevelDefaultsToError(t *testing.T) {
	unknownLevels := []string{"trace", "", "verbose"}
	for _, level := range unknownLevels {
		logger.Configure(level, "text")
		assert.False(t, slog.Default().Enabled(context.Background(), slog.LevelWarn),
			"Configure(%q) should default to error level, but warn is enabled", level)
	}
}

func TestConfigure_UnknownFormatDefaultsToText(t *testing.T) {
	// unknown format should not panic — defaults to text handler
	logger.Configure("info", "logfmt")
	logger.Configure("info", "")
}
