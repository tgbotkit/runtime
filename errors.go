package runtime

import "errors"

// ErrUpdateSourceClosed is returned when the configured update source closes
// its update channel before the run context is canceled.
var ErrUpdateSourceClosed = errors.New("update source closed")
