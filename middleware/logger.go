package middleware

import (
	"context"

	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/logger"
)

// Logger returns a middleware that logs event processing.
func Logger(l logger.Logger) eventemitter.Middleware {
	return eventemitter.MiddlewareFunc(func(next eventemitter.Listener) eventemitter.Listener {
		return eventemitter.ListenerFunc(func(ctx context.Context, payload any) error {
			l.Debugf("handling event: %T", payload)

			err := next.Handle(ctx, payload)
			if err != nil {
				l.Errorf("error handling event %T: %v", payload, err)
			}

			return err
		})
	})
}
