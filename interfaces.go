package runtime

import (
	"context"

	"github.com/tgbotkit/client"
)

// UpdateChan is a channel that receives updates from the Telegram API.
type UpdateChan <-chan client.Update

// UpdateSource represents a source of updates.
type UpdateSource interface {
	// UpdateChan returns a channel that receives updates.
	UpdateChan() <-chan client.Update
	// Start starts the update source. The context is used only for the startup timeout
	// and is not the application lifecycle context.
	Start(ctx context.Context) error
	// Stop stops the update source. The context is used only for the shutdown timeout
	// and is not the application lifecycle context.
	Stop(ctx context.Context) error
}