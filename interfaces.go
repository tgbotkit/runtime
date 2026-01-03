package runtime

import (
	"github.com/tgbotkit/client"
)

type UpdateChan <-chan client.Update

// UpdateSource represents a source of updates.
type UpdateSource interface {
	UpdateChan() <-chan client.Update
}
