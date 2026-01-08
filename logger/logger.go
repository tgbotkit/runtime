// Package logger defines a common logging interface and provides implementations for various logging libraries.
package logger

// Logger is an interface for logging.
type Logger interface {
	// Errorf logs a message at the error level.
	Errorf(format string, args ...interface{})
	// Fatalf logs a message at the fatal level and calls os.Exit(1).
	Fatalf(format string, args ...interface{})
	// Fatal logs a message at the fatal level and calls os.Exit(1).
	Fatal(args ...interface{})
	// Infof logs a message at the info level.
	Infof(format string, args ...interface{})
	// Info logs a message at the info level.
	Info(args ...interface{})
	// Warnf logs a message at the warn level.
	Warnf(format string, args ...interface{})
	// Debugf logs a message at the debug level.
	Debugf(format string, args ...interface{})
	// Debug logs a message at the debug level.
	Debug(args ...interface{})
}
