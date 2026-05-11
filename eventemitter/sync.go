package eventemitter

import (
	"context"
	"errors"
	"path"
	"sort"
	"sync"
)

// listenerEntry represents a registered listener.
type listenerEntry struct {
	Listener Listener
	Once     bool
	Event    string // Store the event pattern this listener was registered for
	sequence uint64
}

type middlewareEntry struct {
	Middleware Middleware
	sequence   uint64
}

// SyncEventEmitter is a concrete implementation of EventEmitter.
type SyncEventEmitter struct {
	mu           sync.RWMutex
	listeners    map[string][]*listenerEntry
	middleware   map[string][]*middlewareEntry
	nextSequence uint64
	opts         Options
}

var _ EventEmitter = (*SyncEventEmitter)(nil)

// NewSync creates a new SyncEventEmitter.
func NewSync(opts Options) (*SyncEventEmitter, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	return &SyncEventEmitter{
		listeners:  make(map[string][]*listenerEntry),
		middleware: make(map[string][]*middlewareEntry),
		opts:       opts,
	}, nil
}

// AddListener adds a listener for the given event.
func (e *SyncEventEmitter) AddListener(event string, listener Listener) UnsubscribeFunc {
	e.mu.Lock()
	defer e.mu.Unlock()

	entry := &listenerEntry{Listener: listener, Once: false, Event: event, sequence: e.nextSequenceID()}
	e.listeners[event] = append(e.listeners[event], entry)

	return func() {
		e.removeListener(event, entry)
	}
}

// Once registers a listener that will be called only once.
func (e *SyncEventEmitter) Once(event string, listener Listener) UnsubscribeFunc {
	e.mu.Lock()
	defer e.mu.Unlock()

	entry := &listenerEntry{Listener: listener, Once: true, Event: event, sequence: e.nextSequenceID()}
	e.listeners[event] = append(e.listeners[event], entry)

	return func() {
		e.removeListener(event, entry)
	}
}

// Emit notifies all listeners of the given event with the provided payload.
func (e *SyncEventEmitter) Emit(ctx context.Context, event string, payload any) {
	entries, middleware := e.getListenersAndMiddleware(event)

	for _, entry := range entries {
		if entry.Once && !e.claimOnce(entry) {
			continue
		}

		if stop := e.handleEntry(ctx, event, payload, entry, middleware); stop {
			break
		}
	}
}

// Use applies middleware to the given event.
func (e *SyncEventEmitter) Use(event string, middleware ...Middleware) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, mw := range middleware {
		entry := &middlewareEntry{Middleware: mw, sequence: e.nextSequenceID()}
		e.middleware[event] = append(e.middleware[event], entry)
	}
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

func (e *SyncEventEmitter) getListenersAndMiddleware(event string) ([]*listenerEntry, []Middleware) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var entries []*listenerEntry

	var middlewareEntries []*middlewareEntry

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
			middlewareEntries = append(middlewareEntries, mws...)
		}
	}

	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].sequence < entries[j].sequence
	})
	sort.SliceStable(middlewareEntries, func(i, j int) bool {
		return middlewareEntries[i].sequence < middlewareEntries[j].sequence
	})

	middleware := make([]Middleware, 0, len(middlewareEntries))
	for _, entry := range middlewareEntries {
		middleware = append(middleware, entry.Middleware)
	}

	return entries, middleware
}

func (e *SyncEventEmitter) nextSequenceID() uint64 {
	e.nextSequence++

	return e.nextSequence
}

func (e *SyncEventEmitter) handleEntry(
	ctx context.Context,
	event string,
	payload any,
	entry *listenerEntry,
	middleware []Middleware,
) bool {
	// Chain middleware and the listener.
	listener := entry.Listener
	for i := len(middleware) - 1; i >= 0; i-- {
		listener = middleware[i].Handle(listener)
	}

	if err := listener.Handle(ctx, payload); err != nil {
		// ErrBreak stops propagation without being an error
		if errors.Is(err, ErrBreak) {
			return true
		}

		if e.opts.errorHandler != nil {
			e.opts.errorHandler(event, err)
		}

		if e.opts.stopOnError {
			return true
		}
	}

	return false
}

// removeListener removes a specific listener for the given event.
func (e *SyncEventEmitter) removeListener(event string, entry *listenerEntry) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.internalRemoveListener(event, entry)
}

func (e *SyncEventEmitter) internalRemoveListener(event string, entry *listenerEntry) {
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

func (e *SyncEventEmitter) claimOnce(entry *listenerEntry) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	entries := e.listeners[entry.Event]
	for _, existingEntry := range entries {
		if existingEntry == entry {
			e.internalRemoveListener(entry.Event, entry)

			return true
		}
	}

	return false
}
