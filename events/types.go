package events

import (
	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime/messagetype"
)

// UpdateEvent is emitted when a new update is received.
type UpdateEvent struct {
	// Update is the received update.
	Update *client.Update
}

// MessageEvent is emitted when a new message is received.
type MessageEvent struct {
	// Message is the received message.
	Message *client.Message
	// Type is the type of the message.
	Type messagetype.MessageType
}

// CommandEvent is emitted when a command is received.
type CommandEvent struct {
	// Message is the received message.
	Message *client.Message
	// Command is the received command name (without /).
	Command string
	// Args is the text following the command.
	Args string
}

// CallbackQueryEvent is emitted when a callback query is received.
type CallbackQueryEvent struct {
	CallbackQuery *client.CallbackQuery
}

// InlineQueryEvent is emitted when an inline query is received.
type InlineQueryEvent struct {
	InlineQuery *client.InlineQuery
}

// ChosenInlineResultEvent is emitted when an inline result is chosen.
type ChosenInlineResultEvent struct {
	ChosenInlineResult *client.ChosenInlineResult
}

// ShippingQueryEvent is emitted when a shipping query is received.
type ShippingQueryEvent struct {
	ShippingQuery *client.ShippingQuery
}

// PreCheckoutQueryEvent is emitted when a pre-checkout query is received.
type PreCheckoutQueryEvent struct {
	PreCheckoutQuery *client.PreCheckoutQuery
}

// PollEvent is emitted when a poll update is received.
type PollEvent struct {
	Poll *client.Poll
}

// PollAnswerEvent is emitted when a poll answer update is received.
type PollAnswerEvent struct {
	PollAnswer *client.PollAnswer
}

// ChatMemberEvent is emitted when a chat member update is received.
type ChatMemberEvent struct {
	ChatMember *client.ChatMemberUpdated
}

// ChatJoinRequestEvent is emitted when a chat join request is received.
type ChatJoinRequestEvent struct {
	ChatJoinRequest *client.ChatJoinRequest
}

// ChatBoostEvent is emitted when a chat boost update is received.
type ChatBoostEvent struct {
	ChatBoost *client.ChatBoostUpdated
}

// RemovedChatBoostEvent is emitted when a chat boost is removed.
type RemovedChatBoostEvent struct {
	RemovedChatBoost *client.ChatBoostRemoved
}

// MessageReactionEvent is emitted when a message reaction update is received.
type MessageReactionEvent struct {
	MessageReaction *client.MessageReactionUpdated
}

// MessageReactionCountEvent is emitted when a reaction-count update is received.
type MessageReactionCountEvent struct {
	MessageReactionCount *client.MessageReactionCountUpdated
}

// BusinessConnectionEvent is emitted when a business connection changes.
type BusinessConnectionEvent struct {
	BusinessConnection *client.BusinessConnection
}

// DeletedBusinessMessagesEvent is emitted when business messages are deleted.
type DeletedBusinessMessagesEvent struct {
	DeletedBusinessMessages *client.BusinessMessagesDeleted
}

// PurchasedPaidMediaEvent is emitted when paid media is purchased.
type PurchasedPaidMediaEvent struct {
	PurchasedPaidMedia *client.PaidMediaPurchased
}

// ManagedBotEvent is emitted when a managed bot update is received.
type ManagedBotEvent struct {
	ManagedBot *client.ManagedBotUpdated
}

// SubscriptionEvent is emitted when a bot subscription update is received.
type SubscriptionEvent struct {
	Subscription *client.BotSubscriptionUpdated
}
