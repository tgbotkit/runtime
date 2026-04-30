package middleware

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/logger"
)

type recoveredPanicError struct {
	value any
}

func (e recoveredPanicError) Error() string {
	return fmt.Sprintf("panic recovered: %v", e.value)
}

// Recoverer returns a middleware that recovers from panics in listeners.
func Recoverer(log logger.Logger) eventemitter.Middleware {
	return eventemitter.MiddlewareFunc(func(next eventemitter.Listener) eventemitter.Listener {
		return eventemitter.ListenerFunc(func(ctx context.Context, payload any) (err error) {
			defer func() {
				if r := recover(); r != nil {
					log.Errorf("panic recovered: %v\n%s", r, debug.Stack())
					err = recoveredPanicError{value: r}
				}
			}()

			return next.Handle(ctx, payload)
		})
	})
}
