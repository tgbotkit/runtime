package handlers_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/events"
	"github.com/tgbotkit/runtime/handlers"
	"github.com/tgbotkit/runtime/logger"
	"github.com/tgbotkit/runtime/messagetype"
)

func TestRegistry(t *testing.T) {
	ee, err := eventemitter.NewSync(eventemitter.NewOptions())
	assert.NoError(t, err)

	reg := handlers.NewRegistry(ee, logger.NewNop())

	t.Run("OnUpdate", func(t *testing.T) {
		var called bool
		var payload *events.UpdateEvent
		handler := func(_ context.Context, event *events.UpdateEvent) error {
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
		handler := func(_ context.Context, event *events.MessageEvent) error {
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

	t.Run("OnMessageType", func(t *testing.T) {
		var called bool
		var payload *events.MessageEvent
		handler := func(_ context.Context, event *events.MessageEvent) error {
			called = true
			payload = event
			return nil
		}

		unsub := reg.OnMessageType(messagetype.Text, handler)

		// Should not be called for Photo
		expectedPayloadPhoto := &events.MessageEvent{Type: messagetype.Photo}
		ee.Emit(context.Background(), events.OnMessage, expectedPayloadPhoto)
		assert.False(t, called, "handler should not be called for Photo")

		// Should be called for Text
		expectedPayloadText := &events.MessageEvent{Type: messagetype.Text}
		ee.Emit(context.Background(), events.OnMessage, expectedPayloadText)
		assert.True(t, called, "handler should be called for Text")
		assert.Same(t, expectedPayloadText, payload, "payload should match")

		// Test unsubscribe
		called = false
		payload = nil
		unsub()
		ee.Emit(context.Background(), events.OnMessage, expectedPayloadText)
		assert.False(t, called, "handler should not be called after unsubscribe")
		assert.Nil(t, payload, "payload should be nil after unsubscribe")
	})

	t.Run("OnCommand", func(t *testing.T) {
		var called bool
		var payload *events.CommandEvent
		handler := func(_ context.Context, event *events.CommandEvent) error {
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
