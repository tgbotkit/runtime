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

// CallbackQueryHandler is a function that handles a callback query event.
type CallbackQueryHandler func(ctx context.Context, event *events.CallbackQueryEvent) error

// InlineQueryHandler is a function that handles an inline query event.
type InlineQueryHandler func(ctx context.Context, event *events.InlineQueryEvent) error

// ChosenInlineResultHandler is a function that handles a chosen inline result event.
type ChosenInlineResultHandler func(ctx context.Context, event *events.ChosenInlineResultEvent) error

// ShippingQueryHandler is a function that handles a shipping query event.
type ShippingQueryHandler func(ctx context.Context, event *events.ShippingQueryEvent) error

// PreCheckoutQueryHandler is a function that handles a pre-checkout query event.
type PreCheckoutQueryHandler func(ctx context.Context, event *events.PreCheckoutQueryEvent) error

// PollHandler is a function that handles a poll event.
type PollHandler func(ctx context.Context, event *events.PollEvent) error

// PollAnswerHandler is a function that handles a poll answer event.
type PollAnswerHandler func(ctx context.Context, event *events.PollAnswerEvent) error

// ChatMemberHandler is a function that handles a chat member event.
type ChatMemberHandler func(ctx context.Context, event *events.ChatMemberEvent) error

// ChatJoinRequestHandler is a function that handles a chat join request event.
type ChatJoinRequestHandler func(ctx context.Context, event *events.ChatJoinRequestEvent) error

// ChatBoostHandler is a function that handles a chat boost event.
type ChatBoostHandler func(ctx context.Context, event *events.ChatBoostEvent) error

// RemovedChatBoostHandler is a function that handles a removed chat boost event.
type RemovedChatBoostHandler func(ctx context.Context, event *events.RemovedChatBoostEvent) error

// MessageReactionHandler is a function that handles a message reaction event.
type MessageReactionHandler func(ctx context.Context, event *events.MessageReactionEvent) error

// MessageReactionCountHandler is a function that handles a message reaction count event.
type MessageReactionCountHandler func(ctx context.Context, event *events.MessageReactionCountEvent) error

// BusinessConnectionHandler is a function that handles a business connection event.
type BusinessConnectionHandler func(ctx context.Context, event *events.BusinessConnectionEvent) error

// DeletedBusinessMessagesHandler is a function that handles a deleted business messages event.
type DeletedBusinessMessagesHandler func(ctx context.Context, event *events.DeletedBusinessMessagesEvent) error

// PurchasedPaidMediaHandler is a function that handles a purchased paid media event.
type PurchasedPaidMediaHandler func(ctx context.Context, event *events.PurchasedPaidMediaEvent) error

// ManagedBotHandler is a function that handles a managed bot event.
type ManagedBotHandler func(ctx context.Context, event *events.ManagedBotEvent) error
