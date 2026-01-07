package handlers

import (
	"context"

	"github.com/tgbotkit/runtime/events"
)

type UpdateHandler func(ctx context.Context, event *events.UpdateEvent) error

type MessageHandler func(ctx context.Context, event *events.MessageEvent) error

type CommandHandler func(ctx context.Context, event *events.CommandEvent) error
