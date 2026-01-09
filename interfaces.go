package runtime

import (
	"github.com/metalagman/appkit/lifecycle"
	"github.com/tgbotkit/client"
)

// UpdateChan is a channel that receives updates from the Telegram API.
type UpdateChan <-chan client.Update

// UpdateSource represents a source of updates.
type UpdateSource interface {
	lifecycle.Lifecycle

	// UpdateChan returns a channel that receives updates.
	UpdateChan() <-chan client.Update
}