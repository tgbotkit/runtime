package eventemitter

import (
	"context"
	"errors"
	"path"
	"sync"
)

// listenerEntry represents a registered listener.
type listenerEntry struct {
	Listener Listener
	Once     bool
	Event    string // Store the event pattern this listener was registered for
}

// SyncEventEmitter is a concrete implementation of EventEmitter.
type SyncEventEmitter struct {
	mu         sync.RWMutex
	listeners  map[string][]*listenerEntry
	middleware map[string][]Middleware
	opts       Options
}

var _ EventEmitter = (*SyncEventEmitter)(nil)

// NewSync creates a new SyncEventEmitter.
func NewSync(opts Options) (*SyncEventEmitter, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	return &SyncEventEmitter{
		listeners:  make(map[string][]*listenerEntry),
		middleware: make(map[string][]Middleware),
		opts:       opts,
	}, nil
}

// AddListener adds a listener for the given event.
func (e *SyncEventEmitter) AddListener(event string, listener Listener) UnsubscribeFunc {
	e.mu.Lock()
	defer e.mu.Unlock()

	entry := &listenerEntry{Listener: listener, Once: false, Event: event}
	e.listeners[event] = append(e.listeners[event], entry)
	return func() {
		e.removeListener(event, entry)
	}
}

// Once registers a listener that will be called only once.
func (e *SyncEventEmitter) Once(event string, listener Listener) UnsubscribeFunc {
	e.mu.Lock()
	defer e.mu.Unlock()

	entry := &listenerEntry{Listener: listener, Once: true, Event: event}
	e.listeners[event] = append(e.listeners[event], entry)
	return func() {
		e.removeListener(event, entry)
	}
}

// removeListener removes a specific listener for the given event.
func (e *SyncEventEmitter) removeListener(event string, entry *listenerEntry) {
	e.mu.Lock()
	defer e.mu.Unlock()

	entries := e.listeners[event]
	for i, existingEntry := range entries {
		if existingEntry == entry {
			e.listeners[event] = append(entries[:i], entries[i+1:]...)
			// Clean up empty slices to keep the map clean
			if len(e.listeners[event]) == 0 {
				delete(e.listeners, event)
			}
			return
		}
	}
}

// Emit notifies all listeners of the given event with the provided payload.
func (e *SyncEventEmitter) Emit(ctx context.Context, event string, payload any) {
	e.mu.RLock()
	
	var entries []*listenerEntry
	var middleware []Middleware

	// Find all matching listeners and middleware
	for pattern, listeners := range e.listeners {
		matched, err := path.Match(pattern, event)
		if err == nil && matched {
			entries = append(entries, listeners...)
		}
	}

	for pattern, mws := range e.middleware {
		matched, err := path.Match(pattern, event)
		if err == nil && matched {
			middleware = append(middleware, mws...)
		}
	}
	
	e.mu.RUnlock()

	var toRemove []*listenerEntry

	for _, entry := range entries {
		// Chain middleware and the listener.
		listener := entry.Listener
		for i := len(middleware) - 1; i >= 0; i-- {
			listener = middleware[i].Handle(listener)
		}

		if err := listener.Handle(ctx, payload); err != nil {
			// ErrBreak stops propagation without being an error
			if errors.Is(err, ErrBreak) {
				if entry.Once {
					toRemove = append(toRemove, entry)
				}
				break
			}

			if e.opts.errorHandler != nil {
				e.opts.errorHandler(event, err)
			}

			if e.opts.stopOnError {
				if entry.Once {
					toRemove = append(toRemove, entry)
				}
				break
			}
		}

		if entry.Once {
			toRemove = append(toRemove, entry)
		}
	}

	// Remove once listeners after execution
	for _, entry := range toRemove {
		e.removeListener(entry.Event, entry)
	}
}

// Use applies middleware to the given event.
func (e *SyncEventEmitter) Use(event string, middleware ...Middleware) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.middleware[event] = append(e.middleware[event], middleware...)
}

// ListenerCount returns the number of listeners for the given event.
func (e *SyncEventEmitter) ListenerCount(event string) int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	
	count := 0
	for pattern, listeners := range e.listeners {
		matched, err := path.Match(pattern, event)
		if err == nil && matched {
			count += len(listeners)
		}
	}
	return count
}

// RemoveAllListeners removes all listeners for the given event.
func (e *SyncEventEmitter) RemoveAllListeners(event string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	// If exact match, delete it
	delete(e.listeners, event)
}
