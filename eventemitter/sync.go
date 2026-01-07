package eventemitter

import (
	"context"
	"errors"
	"reflect"
	"sync"
)

type listenerEntry struct {
	listener Listener
	once     bool
}

// SyncEventEmitter is a concrete implementation of EventEmitter.
type SyncEventEmitter struct {
	mu        sync.RWMutex
	listeners map[string][]listenerEntry
	opts      Options
}

var _ EventEmitter = (*SyncEventEmitter)(nil)

// NewSync creates a new SyncEventEmitter.
func NewSync(opts Options) (*SyncEventEmitter, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	return &SyncEventEmitter{
		listeners: make(map[string][]listenerEntry),
		opts:      opts,
	}, nil
}

// AddListener adds a listener for the given event.
func (e *SyncEventEmitter) AddListener(event string, listener Listener) {
	e.mu.Lock()
	defer e.mu.Unlock()

	entry := listenerEntry{listener: listener, once: false}
	e.listeners[event] = append(e.listeners[event], entry)
}

// Once registers a listener that will be called only once.
func (e *SyncEventEmitter) Once(event string, listener Listener) {
	e.mu.Lock()
	defer e.mu.Unlock()

	entry := listenerEntry{listener: listener, once: true}
	e.listeners[event] = append(e.listeners[event], entry)
}

// RemoveListener removes a specific listener for the given event.
func (e *SyncEventEmitter) RemoveListener(event string, listener Listener) {
	e.mu.Lock()
	defer e.mu.Unlock()

	entries := e.listeners[event]
	for i, entry := range entries {
		// Compare function pointers
		if reflect.ValueOf(entry.listener).Pointer() == reflect.ValueOf(listener).Pointer() {
			e.listeners[event] = append(entries[:i], entries[i+1:]...)
			return
		}
	}
}

// Emit notifies all listeners of the given event with the provided payload.
func (e *SyncEventEmitter) Emit(ctx context.Context, event string, payload any) {
	e.mu.RLock()
	entries := append([]listenerEntry(nil), e.listeners[event]...)
	e.mu.RUnlock()

	var toRemove []Listener

	for _, entry := range entries {
		if err := entry.listener(ctx, payload); err != nil {
			// ErrBreak stops propagation without being an error
			if errors.Is(err, ErrBreak) {
				if entry.once {
					toRemove = append(toRemove, entry.listener)
				}
				break
			}

			if e.opts.errorHandler != nil {
				e.opts.errorHandler(event, err)
			}

			if e.opts.stopOnError {
				if entry.once {
					toRemove = append(toRemove, entry.listener)
				}
				break
			}
		}

		if entry.once {
			toRemove = append(toRemove, entry.listener)
		}
	}

	// Remove once listeners after execution
	for _, listener := range toRemove {
		e.RemoveListener(event, listener)
	}
}

// ListenerCount returns the number of listeners for the given event.
func (e *SyncEventEmitter) ListenerCount(event string) int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.listeners[event])
}

// RemoveAllListeners removes all listeners for the given event.
func (e *SyncEventEmitter) RemoveAllListeners(event string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.listeners, event)
}
