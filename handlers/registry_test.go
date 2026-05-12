package handlers_test

import (
	"context"
	"testing"

	"github.com/tgbotkit/client"
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

	t.Run("OnEditedMessage", func(t *testing.T) {
		var called bool
		var payload *events.MessageEvent
		unsub := reg.OnEditedMessage(func(_ context.Context, event *events.MessageEvent) error {
			called = true
			payload = event

			return nil
		})

		expectedPayload := &events.MessageEvent{Type: messagetype.Text}
		ee.Emit(context.Background(), events.OnEditedMessage, expectedPayload)
		if !called {
			t.Fatal("handler was not called")
		}
		if payload != expectedPayload {
			t.Fatalf("payload mismatch: got %p, want %p", payload, expectedPayload)
		}

		called = false
		payload = nil
		unsub()
		ee.Emit(context.Background(), events.OnEditedMessage, expectedPayload)
		if called {
			t.Fatal("handler was called after unsubscribe")
		}
		if payload != nil {
			t.Fatalf("payload=%v, want nil after unsubscribe", payload)
		}
	})

	t.Run("OnCallbackQuery", func(t *testing.T) {
		var called bool
		var payload *events.CallbackQueryEvent
		unsub := reg.OnCallbackQuery(func(_ context.Context, event *events.CallbackQueryEvent) error {
			called = true
			payload = event

			return nil
		})

		expectedPayload := &events.CallbackQueryEvent{}
		ee.Emit(context.Background(), events.OnCallbackQuery, expectedPayload)
		if !called {
			t.Fatal("handler was not called")
		}
		if payload != expectedPayload {
			t.Fatalf("payload mismatch: got %p, want %p", payload, expectedPayload)
		}

		called = false
		payload = nil
		unsub()
		ee.Emit(context.Background(), events.OnCallbackQuery, expectedPayload)
		if called {
			t.Fatal("handler was called after unsubscribe")
		}
		if payload != nil {
			t.Fatalf("payload=%v, want nil after unsubscribe", payload)
		}
	})

	t.Run("OnCommandName", func(t *testing.T) {
		var called bool
		unsub := reg.OnCommandName("start", func(_ context.Context, event *events.CommandEvent) error {
			called = true
			if event.Command != "start" {
				t.Fatalf("command=%q, want start", event.Command)
			}

			return nil
		})

		ee.Emit(context.Background(), events.OnCommand, &events.CommandEvent{Command: "help"})
		if called {
			t.Fatal("handler should not be called for non-matching command")
		}

		ee.Emit(context.Background(), events.OnCommand, &events.CommandEvent{Command: "start"})
		if !called {
			t.Fatal("handler was not called for matching command")
		}

		called = false
		unsub()
		ee.Emit(context.Background(), events.OnCommand, &events.CommandEvent{Command: "start"})
		if called {
			t.Fatal("handler was called after unsubscribe")
		}
	})

	t.Run("OnCommandMatch", func(t *testing.T) {
		var called bool
		reg.OnCommandMatch(handlers.CommandAny("help", "about"), func(_ context.Context, _ *events.CommandEvent) error {
			called = true

			return nil
		})

		ee.Emit(context.Background(), events.OnCommand, &events.CommandEvent{Command: "start"})
		if called {
			t.Fatal("handler should not be called for non-matching command")
		}

		ee.Emit(context.Background(), events.OnCommand, &events.CommandEvent{Command: "about"})
		if !called {
			t.Fatal("handler was not called for matching command")
		}
	})

	t.Run("OnCallbackData", func(t *testing.T) {
		var called bool
		data := "settings:open"
		reg.OnCallbackData(data, func(_ context.Context, event *events.CallbackQueryEvent) error {
			called = true
			if event.CallbackQuery == nil || event.CallbackQuery.Data == nil || *event.CallbackQuery.Data != data {
				t.Fatalf("callback data mismatch: %#v", event.CallbackQuery)
			}

			return nil
		})

		other := "settings:close"
		ee.Emit(context.Background(), events.OnCallbackQuery, &events.CallbackQueryEvent{
			CallbackQuery: &client.CallbackQuery{Data: &other},
		})
		if called {
			t.Fatal("handler should not be called for non-matching data")
		}

		ee.Emit(context.Background(), events.OnCallbackQuery, &events.CallbackQueryEvent{
			CallbackQuery: &client.CallbackQuery{Data: &data},
		})
		if !called {
			t.Fatal("handler was not called for matching data")
		}
	})

	t.Run("OnCallbackDataPrefix", func(t *testing.T) {
		var called bool
		reg.OnCallbackDataPrefix("settings:", func(_ context.Context, _ *events.CallbackQueryEvent) error {
			called = true

			return nil
		})

		data := "profile:open"
		ee.Emit(context.Background(), events.OnCallbackQuery, &events.CallbackQueryEvent{
			CallbackQuery: &client.CallbackQuery{Data: &data},
		})
		if called {
			t.Fatal("handler should not be called for non-matching prefix")
		}

		data = "settings:open"
		ee.Emit(context.Background(), events.OnCallbackQuery, &events.CallbackQueryEvent{
			CallbackQuery: &client.CallbackQuery{Data: &data},
		})
		if !called {
			t.Fatal("handler was not called for matching prefix")
		}
	})

	t.Run("OnMessageMatch", func(t *testing.T) {
		var called bool
		text := "ping"
		reg.OnMessageMatch(handlers.MessageText(text), func(_ context.Context, event *events.MessageEvent) error {
			called = true
			if event.Message == nil || event.Message.Text == nil || *event.Message.Text != text {
				t.Fatalf("message text mismatch: %#v", event.Message)
			}

			return nil
		})

		other := "pong"
		ee.Emit(context.Background(), events.OnMessage, &events.MessageEvent{
			Message: &client.Message{Text: &other},
		})
		if called {
			t.Fatal("handler should not be called for non-matching text")
		}

		ee.Emit(context.Background(), events.OnMessage, &events.MessageEvent{
			Message: &client.Message{Text: &text},
		})
		if !called {
			t.Fatal("handler was not called for matching text")
		}
	})
}
