package eventemitter

import "context"

type TypedListener[T any] func(ctx context.Context, payload *T) error

// On registers a typed handler for a specific event.
func On[E any](ee EventEmitter, event string, handler TypedListener[E]) UnsubscribeFunc {
	listener := ListenerFunc(func(ctx context.Context, payload any) error {
		if e, ok := payload.(*E); ok {
			return handler(ctx, e)
		}
		return nil
	})
	return ee.AddListener(event, listener)
}
