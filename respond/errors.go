package respond

import "errors"

// ErrNilClient is returned when a responder has no Telegram API client.
var ErrNilClient = errors.New("nil telegram client")

// ErrNilMessage is returned when a message-based target cannot be built.
var ErrNilMessage = errors.New("nil message")

// ErrNilCallbackQuery is returned when a callback answer has no callback query.
var ErrNilCallbackQuery = errors.New("nil callback query")

// ErrNoMessageTarget is returned when a callback query has no usable message target.
var ErrNoMessageTarget = errors.New("message target unavailable")
