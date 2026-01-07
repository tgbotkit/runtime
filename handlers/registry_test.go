package handlers_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/events"
	"github.com/tgbotkit/runtime/handlers"
)

func TestRegistry(t *testing.T) {
	ee, err := eventemitter.NewSync(eventemitter.NewOptions())
	assert.NoError(t, err)

	reg := handlers.NewRegistry(ee)

	t.Run("OnUpdate", func(t *testing.T) {
		var called bool
		var payload *events.UpdateEvent
		handler := func(ctx context.Context, event *events.UpdateEvent) error {
			called = true
			payload = event
			return nil
		}

		unsub := reg.OnUpdate(handler)
		expectedPayload := &events.UpdateEvent{}
		ee.Emit(context.Background(), events.OnUpdate, expectedPayload)
		assert.True(t, called, "handler should be called")
		assert.Same(t, expectedPayload, payload, "payload should match")

		// Test unsubscribe
		called = false
		payload = nil
		unsub()
		ee.Emit(context.Background(), events.OnUpdate, expectedPayload)
		assert.False(t, called, "handler should not be called after unsubscribe")
		assert.Nil(t, payload, "payload should be nil after unsubscribe")
	})

	t.Run("OnMessage", func(t *testing.T) {
		var called bool
		var payload *events.MessageEvent
		handler := func(ctx context.Context, event *events.MessageEvent) error {
			called = true
			payload = event
			return nil
		}

		unsub := reg.OnMessage(handler)
		expectedPayload := &events.MessageEvent{}
		ee.Emit(context.Background(), events.OnMessage, expectedPayload)
		assert.True(t, called, "handler should be called")
		assert.Same(t, expectedPayload, payload, "payload should match")

		// Test unsubscribe
		called = false
		payload = nil
		unsub()
		ee.Emit(context.Background(), events.OnMessage, expectedPayload)
		assert.False(t, called, "handler should not be called after unsubscribe")
		assert.Nil(t, payload, "payload should be nil after unsubscribe")
	})

	t.Run("OnCommand", func(t *testing.T) {
		var called bool
		var payload *events.CommandEvent
		handler := func(ctx context.Context, event *events.CommandEvent) error {
			called = true
			payload = event
			return nil
		}

		unsub := reg.OnCommand(handler)
		expectedPayload := &events.CommandEvent{}
		ee.Emit(context.Background(), events.OnCommand, expectedPayload)
		assert.True(t, called, "handler should be called")
		assert.Same(t, expectedPayload, payload, "payload should match")

		// Test unsubscribe
		called = false
		payload = nil
		unsub()
		ee.Emit(context.Background(), events.OnCommand, expectedPayload)
		assert.False(t, called, "handler should not be called after unsubscribe")
		assert.Nil(t, payload, "payload should be nil after unsubscribe")
	})
}
