package botcontext

import (
	"context"

	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/logger"
)

// BotContext provides access to bot capabilities without coupling to internals.
// It is embedded in context.Context and passed to handlers.
type BotContext interface {
	Client() client.ClientWithResponsesInterface
	EventEmitter() eventemitter.EventEmitter
	Logger() logger.Logger
}

type contextKey struct{}

var botContextKey contextKey

// WithBotContext returns a new context with the BotContext embedded.
func WithBotContext(ctx context.Context, bot BotContext) context.Context {
	return context.WithValue(ctx, botContextKey, bot)
}

// FromContext retrieves the BotContext from the context.
// It returns nil if the BotContext is not found.
func FromContext(ctx context.Context) BotContext {
	val, _ := ctx.Value(botContextKey).(BotContext)
	return val
}
