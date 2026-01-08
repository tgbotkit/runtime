// Package eventemitter provides a flexible event bus implementation for handling bot events.
package eventemitter

import (
	"errors"
)

// ErrBreak is a special error that can be returned by a listener to stop further event propagation.
var ErrBreak = errors.New("break")