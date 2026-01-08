package logger

import (
	"log"
	"os"
)

// Slog is a logger that uses the standard library's log package.
type Slog struct{}

// NewSlog creates a new Slog.
func NewSlog() *Slog {
	return &Slog{}
}

// Errorf logs a message at the error level.
func (s *Slog) Errorf(format string, args ...interface{}) {
	log.Printf("ERROR: "+format, args...)
}

// Fatalf logs a message at the fatal level and calls os.Exit(1).
func (s *Slog) Fatalf(format string, args ...interface{}) {
	log.Printf("FATAL: "+format, args...)
	os.Exit(1)
}

// Fatal logs a message at the fatal level and calls os.Exit(1).
func (s *Slog) Fatal(args ...interface{}) {
	log.Print(append([]interface{}{"FATAL: "}, args...)...)
	os.Exit(1)
}

// Infof logs a message at the info level.
func (s *Slog) Infof(format string, args ...interface{}) {
	log.Printf("INFO: "+format, args...)
}

// Info logs a message at the info level.
func (s *Slog) Info(args ...interface{}) {
	log.Print(append([]interface{}{"INFO: "}, args...)...)
}

// Warnf logs a message at the warn level.
func (s *Slog) Warnf(format string, args ...interface{}) {
	log.Printf("WARN: "+format, args...)
}

// Debugf logs a message at the debug level.
func (s *Slog) Debugf(format string, args ...interface{}) {
	log.Printf("DEBUG: "+format, args...)
}

// Debug logs a message at the debug level.
func (s *Slog) Debug(args ...interface{}) {
	log.Print(append([]interface{}{"DEBUG: "}, args...)...)
}