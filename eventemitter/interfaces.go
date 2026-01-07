package eventemitter

import "context"

// Listener is an interface that handles an event.
type Listener interface {
	Handle(ctx context.Context, payload any) error
}

// ListenerFunc is an adapter to allow the use of ordinary functions as Listener.
type ListenerFunc func(ctx context.Context, payload any) error

// Handle calls f(ctx, payload).
func (f ListenerFunc) Handle(ctx context.Context, payload any) error {
	return f(ctx, payload)
}

// Middleware is an interface that wraps a listener.
type Middleware interface {
	Handle(next Listener) Listener
}

// MiddlewareFunc is an adapter to allow the use of ordinary functions as Middleware.
type MiddlewareFunc func(next Listener) Listener

// Handle calls f(next).
func (f MiddlewareFunc) Handle(next Listener) Listener {
	return f(next)
}

// ErrorHandler is called when a listener returns an error.
type ErrorHandler func(event string, err error)

// UnsubscribeFunc is a function that unregisters a listener.
type UnsubscribeFunc func()

// EventEmitter defines the interface for event management.
type EventEmitter interface {
	// AddListener registers a listener for the given event.
	AddListener(event string, listener Listener) UnsubscribeFunc
	// Once registers a listener that will be called only once.
	Once(event string, listener Listener) UnsubscribeFunc
	// Emit notifies all listeners of the given event with the provided payload.
	Emit(ctx context.Context, event string, payload any)
	// Use applies middleware to the given event.
	Use(event string, middleware ...Middleware)
	// ListenerCount returns the number of listeners for the given event.
	ListenerCount(event string) int
	// RemoveAllListeners removes all listeners for the given event.
	RemoveAllListeners(event string)
}
