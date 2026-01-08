package runtime

import (
	"github.com/tgbotkit/client"
)

// UpdateChan is a channel that receives updates from the Telegram API.
type UpdateChan <-chan client.Update

// UpdateSource represents a source of updates.
type UpdateSource interface {
	UpdateChan() <-chan client.Update
}