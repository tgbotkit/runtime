package middleware

import (
	"context"

	"github.com/tgbotkit/runtime/botcontext"
	"github.com/tgbotkit/runtime/eventemitter"
)

func ContextInjector(bot botcontext.BotContext) eventemitter.Middleware {
	return eventemitter.MiddlewareFunc(func(next eventemitter.Listener) eventemitter.Listener {
		return eventemitter.ListenerFunc(func(ctx context.Context, payload any) error {
			ctx = botcontext.WithBotContext(ctx, bot)
			return next.Handle(ctx, payload)
		})
	})
}
