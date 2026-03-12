package logger

// Logger defines the interface for logging operations.
// This interface enables dependency injection and allows for alternative implementations
// such as no-op loggers for testing or custom loggers for different environments.
type Logger interface {
	// Trace logs a trace-level message with function name and formatted string
	Trace(funcName, format string, args ...any)

	// Debug logs a debug-level message with function name and formatted string
	Debug(funcName, format string, args ...any)

	// Info logs an info-level message with function name and formatted string
	Info(funcName, format string, args ...any)

	// Warn logs a warning-level message with function name and formatted string
	Warn(funcName, format string, args ...any)

	// Error logs an error-level message with function name, error, and formatted string
	Error(funcName string, err error, format string, args ...any)
}
