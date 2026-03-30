package logger_test

import (
	"testing"

	"github.com/unicrons/aws-root-manager/internal/logger"
)

func TestConfigure_ValidLevels(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error"}
	for _, level := range levels {
		t.Run(level, func(t *testing.T) {
			logger.Configure(level, "text")
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
	// unknown level should not panic — defaults to error
	logger.Configure("trace", "text")
	logger.Configure("", "text")
	logger.Configure("verbose", "text")
}

func TestConfigure_UnknownFormatDefaultsToText(t *testing.T) {
	// unknown format should not panic — defaults to text handler
	logger.Configure("info", "logfmt")
	logger.Configure("info", "")
}
