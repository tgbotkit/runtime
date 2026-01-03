package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

// Slog is a logger implementation that uses log/slog.
type Slog struct {
	logger *slog.Logger
}

// NewSlog creates a new Slog logger.
func NewSlog(l *slog.Logger) Logger {
	if l == nil {
		return NewNop()
	}
	return &Slog{logger: l}
}

func (s *Slog) Errorf(format string, args ...interface{}) {
	s.logger.Log(context.Background(), slog.LevelError, fmt.Sprintf(format, args...))
}

func (s *Slog) Fatalf(format string, args ...interface{}) {
	s.logger.Log(context.Background(), slog.LevelError, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (s *Slog) Fatal(args ...interface{}) {
	s.logger.Log(context.Background(), slog.LevelError, fmt.Sprint(args...))
	os.Exit(1)
}

func (s *Slog) Infof(format string, args ...interface{}) {
	s.logger.Log(context.Background(), slog.LevelInfo, fmt.Sprintf(format, args...))
}

func (s *Slog) Info(args ...interface{}) {
	s.logger.Log(context.Background(), slog.LevelInfo, fmt.Sprint(args...))
}

func (s *Slog) Warnf(format string, args ...interface{}) {
	s.logger.Log(context.Background(), slog.LevelWarn, fmt.Sprintf(format, args...))
}

func (s *Slog) Debugf(format string, args ...interface{}) {
	s.logger.Log(context.Background(), slog.LevelDebug, fmt.Sprintf(format, args...))
}

func (s *Slog) Debug(args ...interface{}) {
	s.logger.Log(context.Background(), slog.LevelDebug, fmt.Sprint(args...))
}
