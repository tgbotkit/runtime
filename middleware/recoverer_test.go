package middleware

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tgbotkit/runtime/eventemitter"
)

func TestRecoverer(t *testing.T) {
	t.Run("recovers from panic", func(t *testing.T) {
		logger := &mockLogger{}
		recoverer := Recoverer(logger)

		next := eventemitter.ListenerFunc(func(_ context.Context, _ any) error {
			panic("test panic")
		})

		// Should not panic
		err := recoverer.Handle(next).Handle(context.Background(), "test")
		
		assert.NoError(t, err)
		assert.True(t, logger.errorfCalled, "Errorf should be called on panic")
	})

	t.Run("passes through without panic", func(t *testing.T) {
		logger := &mockLogger{}
		recoverer := Recoverer(logger)

		next := eventemitter.ListenerFunc(func(_ context.Context, _ any) error {
			return nil
		})

		err := recoverer.Handle(next).Handle(context.Background(), "test")
		
		assert.NoError(t, err)
		assert.False(t, logger.errorfCalled, "Errorf should not be called without panic")
	})
}