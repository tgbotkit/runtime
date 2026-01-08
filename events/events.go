// Package events defines the events emitted by the bot and their associated payload types.
package events

// Constants for event names.
const (
	// OnUpdate is emitted when a new update is received from Telegram.
	OnUpdate = "onUpdate"
	// OnMessage is emitted when a new message is received, regardless of its type.
	// The specific type is available in the MessageEvent.Type field.
	OnMessage = "onMessage"
	// OnCommand is emitted when a command is received from a text message.
	OnCommand = "onCommand"
)
