package middleware

import (
	"context"
	"errors"
	"testing"

	"github.com/tgbotkit/runtime/eventemitter"
)

func TestRecoverer(t *testing.T) {
	t.Run("recovers from panic and returns error", func(t *testing.T) {
		logger := &mockLogger{}
		recoverer := Recoverer(logger)

		next := eventemitter.ListenerFunc(func(_ context.Context, _ any) error {
			panic("test panic")
		})

		err := recoverer.Handle(next).Handle(context.Background(), "test")

		if err == nil {
			t.Fatal("Handle() error is nil, want non-nil")
		}
		var recoveredErr recoveredPanicError
		if !errors.As(err, &recoveredErr) {
			t.Fatalf("Handle() error=%v, want recovered panic error", err)
		}
		if !logger.errorfCalled {
			t.Fatal("logger.Errorf was not called")
		}
	})

	t.Run("passes through without panic", func(t *testing.T) {
		logger := &mockLogger{}
		recoverer := Recoverer(logger)

		next := eventemitter.ListenerFunc(func(_ context.Context, _ any) error {
			return nil
		})

		err := recoverer.Handle(next).Handle(context.Background(), "test")
		if err != nil {
			t.Fatalf("Handle() unexpected error: %v", err)
		}
		if logger.errorfCalled {
			t.Fatal("logger.Errorf should not be called")
		}
	})
}
