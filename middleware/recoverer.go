package middleware

import (
	"context"
	"runtime/debug"

	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/logger"
)

// Recoverer returns a middleware that recovers from panics in listeners.
func Recoverer(log logger.Logger) eventemitter.Middleware {
	return eventemitter.MiddlewareFunc(func(next eventemitter.Listener) eventemitter.Listener {
		return eventemitter.ListenerFunc(func(ctx context.Context, payload any) (err error) {
			defer func() {
				if r := recover(); r != nil {
					log.Errorf("panic recovered: %v\n%s", r, debug.Stack())
				}
			}()

			return next.Handle(ctx, payload)
		})
	})
}
