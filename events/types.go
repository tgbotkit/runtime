package events

import (
	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/logger"
)

// BotContext provides access to bot capabilities without coupling to internals.
type BotContext interface {
	Client() client.ClientWithResponsesInterface
	EventEmitter() eventemitter.EventEmitter
	Logger() logger.Logger
}

// BotEvent is a generic event that includes the bot context.
// It is used for simple lifecycle events like OnBeforeStart.
type BotEvent struct {
	Bot BotContext
}

// UpdateEvent is emitted when a new update is received.
type UpdateEvent struct {
	Bot BotContext
	// Update is the received update.
	Update *client.Update
}

// MessageEvent is emitted when a new message is received.
type MessageEvent struct {
	Bot BotContext
	// Message is the received message.
	Message *client.Message
}

// CommandEvent is emitted when a command is received.
type CommandEvent struct {
	Bot BotContext
	// Message is the received message.
	Message *client.Message
	// Command is the received command name (without /).
	Command string
	// Args is the text following the command.
	Args string
}
