// Package middleware provides various middleware implementations for the bot's event emitter.
package middleware

import (
	"context"

	"github.com/tgbotkit/runtime/botcontext"
	"github.com/tgbotkit/runtime/eventemitter"
)

// ContextInjector returns a middleware that injects BotContext into the request context.
func ContextInjector(bot botcontext.BotContext) eventemitter.Middleware {
	return eventemitter.MiddlewareFunc(func(next eventemitter.Listener) eventemitter.Listener {
		return eventemitter.ListenerFunc(func(ctx context.Context, payload any) error {
			ctx = botcontext.WithBotContext(ctx, bot)

			return next.Handle(ctx, payload)
		})
	})
}
