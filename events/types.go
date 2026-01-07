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
