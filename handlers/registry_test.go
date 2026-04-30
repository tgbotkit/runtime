package handlers_test

import (
	"context"
	"testing"

	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/events"
	"github.com/tgbotkit/runtime/handlers"
	"github.com/tgbotkit/runtime/logger"
	"github.com/tgbotkit/runtime/messagetype"
)

func TestRegistry(t *testing.T) {
	ee, err := eventemitter.NewSync(eventemitter.NewOptions())
	if err != nil {
		t.Fatalf("NewSync() unexpected error: %v", err)
	}

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

		if !called {
			t.Fatal("handler was not called")
		}
		if payload != expectedPayload {
			t.Fatalf("payload mismatch: got %p, want %p", payload, expectedPayload)
		}

		called = false
		payload = nil
		unsub()
		ee.Emit(context.Background(), events.OnUpdate, expectedPayload)

		if called {
			t.Fatal("handler was called after unsubscribe")
		}
		if payload != nil {
			t.Fatalf("payload=%v, want nil after unsubscribe", payload)
		}
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

		if !called {
			t.Fatal("handler was not called")
		}
		if payload != expectedPayload {
			t.Fatalf("payload mismatch: got %p, want %p", payload, expectedPayload)
		}

		called = false
		payload = nil
		unsub()
		ee.Emit(context.Background(), events.OnMessage, expectedPayload)

		if called {
			t.Fatal("handler was called after unsubscribe")
		}
		if payload != nil {
			t.Fatalf("payload=%v, want nil after unsubscribe", payload)
		}
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

		expectedPayloadPhoto := &events.MessageEvent{Type: messagetype.Photo}
		ee.Emit(context.Background(), events.OnMessage, expectedPayloadPhoto)
		if called {
			t.Fatal("handler should not be called for photo")
		}

		expectedPayloadText := &events.MessageEvent{Type: messagetype.Text}
		ee.Emit(context.Background(), events.OnMessage, expectedPayloadText)
		if !called {
			t.Fatal("handler was not called for text")
		}
		if payload != expectedPayloadText {
			t.Fatalf("payload mismatch: got %p, want %p", payload, expectedPayloadText)
		}

		called = false
		payload = nil
		unsub()
		ee.Emit(context.Background(), events.OnMessage, expectedPayloadText)

		if called {
			t.Fatal("handler was called after unsubscribe")
		}
		if payload != nil {
			t.Fatalf("payload=%v, want nil after unsubscribe", payload)
		}
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

		if !called {
			t.Fatal("handler was not called")
		}
		if payload != expectedPayload {
			t.Fatalf("payload mismatch: got %p, want %p", payload, expectedPayload)
		}

		called = false
		payload = nil
		unsub()
		ee.Emit(context.Background(), events.OnCommand, expectedPayload)

		if called {
			t.Fatal("handler was called after unsubscribe")
		}
		if payload != nil {
			t.Fatalf("payload=%v, want nil after unsubscribe", payload)
		}
	})
}
