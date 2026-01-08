// Package handlers provides a type-safe way to register and manage event handlers.
package handlers

import (
	"context"

	"github.com/tgbotkit/runtime/events"
)

// UpdateHandler is a function that handles an update event.
type UpdateHandler func(ctx context.Context, event *events.UpdateEvent) error

// MessageHandler is a function that handles a message event.
type MessageHandler func(ctx context.Context, event *events.MessageEvent) error

// CommandHandler is a function that handles a command event.
type CommandHandler func(ctx context.Context, event *events.CommandEvent) error
