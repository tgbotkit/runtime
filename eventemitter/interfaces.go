package eventemitter

import "context"

// Listener is a function that handles an event.
type Listener func(ctx context.Context, payload any) error

// ErrorHandler is called when a listener returns an error.
type ErrorHandler func(event string, err error)

// EventEmitter defines the interface for event management.
type EventEmitter interface {
	// AddListener registers a listener for the given event.
	AddListener(event string, listener Listener)
	// Once registers a listener that will be called only once.
	Once(event string, listener Listener)
	// RemoveListener removes a specific listener for the given event.
	RemoveListener(event string, listener Listener)
	// Emit notifies all listeners of the given event with the provided payload.
	Emit(ctx context.Context, event string, payload any)
	// ListenerCount returns the number of listeners for the given event.
	ListenerCount(event string) int
	// RemoveAllListeners removes all listeners for the given event.
	RemoveAllListeners(event string)
}
