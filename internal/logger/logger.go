package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// logrusLogger implements the Logger interface using logrus.
type logrusLogger struct {
	logger *log.Logger
}

// New creates a new Logger instance configured from environment variables.
func New() Logger {
	logger := log.New()

	lvl, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		lvl = "error" // default
	}

	ll, err := log.ParseLevel(lvl)
	if err != nil {
		ll = log.DebugLevel
	}
	logger.SetLevel(ll)

	format, ok := os.LookupEnv("LOG_FORMAT")
	if ok {
		setLoggerFormat(logger, format)
	}

	return &logrusLogger{logger: logger}
}

func setLoggerFormat(logger *log.Logger, logFormat string) {
	switch logFormat {
	case "json":
		logger.SetFormatter(&log.JSONFormatter{})
	default:
		logger.SetFormatter(&log.TextFormatter{})
	}
}

func (l *logrusLogger) Trace(funcName, format string, args ...any) {
	l.logger.WithField("function", funcName).Tracef(format, args...)
}

func (l *logrusLogger) Debug(funcName, format string, args ...any) {
	l.logger.WithField("function", funcName).Debugf(format, args...)
}

func (l *logrusLogger) Info(funcName, format string, args ...any) {
	l.logger.WithField("function", funcName).Infof(format, args...)
}

func (l *logrusLogger) Warn(funcName, format string, args ...any) {
	l.logger.WithField("function", funcName).Warnf(format, args...)
}

func (l *logrusLogger) Error(funcName string, err error, format string, args ...any) {
	l.logger.WithField("function", funcName).WithError(err).Errorf(format, args...)
}

// Global logger functions for backward compatibility
// These will be deprecated once all code uses dependency injection

var defaultLogger Logger

func init() {
	defaultLogger = New()
}

func SetLoggerFormat(logFormat string) {
	// For backward compatibility with existing code
	switch logFormat {
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.SetFormatter(&log.TextFormatter{})
	}
}

// Wrap logrus with function name - global functions for backward compatibility
func Trace(funcName, format string, args ...any) {
	defaultLogger.Trace(funcName, format, args...)
}

func Debug(funcName, format string, args ...any) {
	defaultLogger.Debug(funcName, format, args...)
}

func Info(funcName, format string, args ...any) {
	defaultLogger.Info(funcName, format, args...)
}

func Warn(funcName, format string, args ...any) {
	defaultLogger.Warn(funcName, format, args...)
}

func Error(funcName string, err error, format string, args ...any) {
	defaultLogger.Error(funcName, err, format, args...)
}
