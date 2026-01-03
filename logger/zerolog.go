package logger

import (
	"fmt"

	"github.com/rs/zerolog"
)

// Zerolog is a logger implementation that uses rs/zerolog.
type Zerolog struct {
	logger zerolog.Logger
}

// NewZerolog creates a new Zerolog logger.
func NewZerolog(l zerolog.Logger) Logger {
	return &Zerolog{logger: l}
}

func (z *Zerolog) Errorf(format string, args ...interface{}) {
	z.logger.Error().Msg(fmt.Sprintf(format, args...))
}

func (z *Zerolog) Fatalf(format string, args ...interface{}) {
	z.logger.Fatal().Msg(fmt.Sprintf(format, args...))
}

func (z *Zerolog) Fatal(args ...interface{}) {
	z.logger.Fatal().Msg(fmt.Sprint(args...))
}

func (z *Zerolog) Infof(format string, args ...interface{}) {
	z.logger.Info().Msg(fmt.Sprintf(format, args...))
}

func (z *Zerolog) Info(args ...interface{}) {
	z.logger.Info().Msg(fmt.Sprint(args...))
}

func (z *Zerolog) Warnf(format string, args ...interface{}) {
	z.logger.Warn().Msg(fmt.Sprintf(format, args...))
}

func (z *Zerolog) Debugf(format string, args ...interface{}) {
	z.logger.Debug().Msg(fmt.Sprintf(format, args...))
}

func (z *Zerolog) Debug(args ...interface{}) {
	z.logger.Debug().Msg(fmt.Sprint(args...))
}
