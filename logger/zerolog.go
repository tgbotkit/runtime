package logger

import (
	"fmt"

	"github.com/rs/zerolog"
)

// Zerolog is a logger that uses the zerolog package.
type Zerolog struct {
	l zerolog.Logger
}

var _ Logger = (*Zerolog)(nil)

// NewZerolog creates a new Zerolog.
func NewZerolog(l zerolog.Logger) *Zerolog {
	return &Zerolog{l: l}
}

// Errorf logs a message at the error level.
func (z *Zerolog) Errorf(format string, args ...interface{}) {
	z.l.Error().Msgf(format, args...)
}

// Fatalf logs a message at the fatal level and calls os.Exit(1).
func (z *Zerolog) Fatalf(format string, args ...interface{}) {
	z.l.Fatal().Msgf(format, args...)
}

// Fatal logs a message at the fatal level and calls os.Exit(1).
func (z *Zerolog) Fatal(args ...interface{}) {
	z.l.Fatal().Msg(fmt.Sprint(args...))
}

// Infof logs a message at the info level.
func (z *Zerolog) Infof(format string, args ...interface{}) {
	z.l.Info().Msgf(format, args...)
}

// Info logs a message at the info level.
func (z *Zerolog) Info(args ...interface{}) {
	z.l.Info().Msg(fmt.Sprint(args...))
}

// Warnf logs a message at the warn level.
func (z *Zerolog) Warnf(format string, args ...interface{}) {
	z.l.Warn().Msgf(format, args...)
}

// Debugf logs a message at the debug level.
func (z *Zerolog) Debugf(format string, args ...interface{}) {
	z.l.Debug().Msgf(format, args...)
}

// Debug logs a message at the debug level.
func (z *Zerolog) Debug(args ...interface{}) {
	z.l.Debug().Msg(fmt.Sprint(args...))
}
