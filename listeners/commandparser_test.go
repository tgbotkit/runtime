package listeners_test

import (
	"context"
	"errors"
	"testing"

	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/events"
	"github.com/tgbotkit/runtime/listeners"
	"github.com/tgbotkit/runtime/messagetype"
)

func TestCommandParser(t *testing.T) {
	ee, err := eventemitter.NewSync(eventemitter.NewOptions())
	if err != nil {
		t.Fatalf("NewSync() unexpected error: %v", err)
	}
	botName := "mybot"
	parser := listeners.CommandParser(ee, botName)

	t.Run("parses simple command", func(t *testing.T) {
		var receivedEvent *events.CommandEvent
		ee.AddListener(events.OnCommand, eventemitter.ListenerFunc(func(_ context.Context, payload any) error {
			if e, ok := payload.(*events.CommandEvent); ok {
				receivedEvent = e
			}
			return nil
		}))

		text := "/start"
		msg := &client.Message{
			MessageId: 1,
			Text:      &text,
			Entities: &[]client.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 6},
			},
		}

		event := &events.MessageEvent{
			Message: msg,
			Type:    messagetype.Text,
		}

		err := parser.Handle(context.Background(), event)
		if !errors.Is(err, eventemitter.ErrBreak) {
			t.Fatalf("Handle() error=%v, want %v", err, eventemitter.ErrBreak)
		}
		if receivedEvent == nil {
			t.Fatal("receivedEvent is nil")
		}
		if receivedEvent.Command != "start" {
			t.Fatalf("command=%q, want %q", receivedEvent.Command, "start")
		}
		if receivedEvent.Args != "" {
			t.Fatalf("args=%q, want empty string", receivedEvent.Args)
		}
	})

	t.Run("parses command with args", func(t *testing.T) {
		ee, err := eventemitter.NewSync(eventemitter.NewOptions())
		if err != nil {
			t.Fatalf("NewSync() unexpected error: %v", err)
		}
		parser := listeners.CommandParser(ee, botName)

		var receivedEvent *events.CommandEvent
		ee.AddListener(events.OnCommand, eventemitter.ListenerFunc(func(_ context.Context, payload any) error {
			e, ok := payload.(*events.CommandEvent)
			if !ok {
				t.Fatalf("payload type=%T, want *events.CommandEvent", payload)
			}
			receivedEvent = e
			return nil
		}))

		text := "/echo hello world"
		msg := &client.Message{
			MessageId: 2,
			Text:      &text,
			Entities: &[]client.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 5},
			},
		}
		event := &events.MessageEvent{Message: msg, Type: messagetype.Text}

		err = parser.Handle(context.Background(), event)
		if !errors.Is(err, eventemitter.ErrBreak) {
			t.Fatalf("Handle() error=%v, want %v", err, eventemitter.ErrBreak)
		}
		if receivedEvent == nil {
			t.Fatal("receivedEvent is nil")
		}
		if receivedEvent.Command != "echo" {
			t.Fatalf("command=%q, want %q", receivedEvent.Command, "echo")
		}
		if receivedEvent.Args != "hello world" {
			t.Fatalf("args=%q, want %q", receivedEvent.Args, "hello world")
		}
	})

	t.Run("parses command with bot mention", func(t *testing.T) {
		ee, err := eventemitter.NewSync(eventemitter.NewOptions())
		if err != nil {
			t.Fatalf("NewSync() unexpected error: %v", err)
		}
		parser := listeners.CommandParser(ee, botName)

		var receivedEvent *events.CommandEvent
		ee.AddListener(events.OnCommand, eventemitter.ListenerFunc(func(_ context.Context, payload any) error {
			e, ok := payload.(*events.CommandEvent)
			if !ok {
				t.Fatalf("payload type=%T, want *events.CommandEvent", payload)
			}
			receivedEvent = e
			return nil
		}))

		text := "/start@mybot args"
		msg := &client.Message{
			MessageId: 3,
			Text:      &text,
			Entities: &[]client.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 12},
			},
		}
		event := &events.MessageEvent{Message: msg, Type: messagetype.Text}

		err = parser.Handle(context.Background(), event)
		if !errors.Is(err, eventemitter.ErrBreak) {
			t.Fatalf("Handle() error=%v, want %v", err, eventemitter.ErrBreak)
		}
		if receivedEvent == nil {
			t.Fatal("receivedEvent is nil")
		}
		if receivedEvent.Command != "start" {
			t.Fatalf("command=%q, want %q", receivedEvent.Command, "start")
		}
		if receivedEvent.Args != "args" {
			t.Fatalf("args=%q, want %q", receivedEvent.Args, "args")
		}
	})

	t.Run("ignores command for other bot", func(t *testing.T) {
		ee, err := eventemitter.NewSync(eventemitter.NewOptions())
		if err != nil {
			t.Fatalf("NewSync() unexpected error: %v", err)
		}
		parser := listeners.CommandParser(ee, botName)

		var called bool
		ee.AddListener(events.OnCommand, eventemitter.ListenerFunc(func(_ context.Context, _ any) error {
			called = true
			return nil
		}))

		text := "/start@otherbot"
		msg := &client.Message{
			MessageId: 4,
			Text:      &text,
			Entities: &[]client.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 15},
			},
		}
		event := &events.MessageEvent{Message: msg, Type: messagetype.Text}

		err = parser.Handle(context.Background(), event)
		if err != nil {
			t.Fatalf("Handle() unexpected error: %v", err)
		}
		if called {
			t.Fatal("OnCommand listener should not be called")
		}
	})

	t.Run("ignores non-text message", func(t *testing.T) {
		ee, err := eventemitter.NewSync(eventemitter.NewOptions())
		if err != nil {
			t.Fatalf("NewSync() unexpected error: %v", err)
		}
		parser := listeners.CommandParser(ee, botName)

		event := &events.MessageEvent{Type: messagetype.Photo}
		err = parser.Handle(context.Background(), event)
		if err != nil {
			t.Fatalf("Handle() unexpected error: %v", err)
		}
	})

	t.Run("ignores invalid payload", func(t *testing.T) {
		err := parser.Handle(context.Background(), "string")
		if err != nil {
			t.Fatalf("Handle() unexpected error: %v", err)
		}
	})
}

func TestCommandParser_StopsOnErrBreakInEmitterChain(t *testing.T) {
	ee, err := eventemitter.NewSync(eventemitter.NewOptions())
	if err != nil {
		t.Fatalf("NewSync() unexpected error: %v", err)
	}
	botName := "mybot"

	ee.AddListener(events.OnMessage, listeners.CommandParser(ee, botName))

	var commandHandled bool
	ee.AddListener(events.OnCommand, eventemitter.ListenerFunc(func(_ context.Context, payload any) error {
		_, ok := payload.(*events.CommandEvent)
		commandHandled = ok
		return nil
	}))

	var afterParserCalled bool
	ee.AddListener(events.OnMessage, eventemitter.ListenerFunc(func(_ context.Context, _ any) error {
		afterParserCalled = true
		return nil
	}))

	text := "/start"
	msg := &client.Message{
		MessageId: 1,
		Text:      &text,
		Entities: &[]client.MessageEntity{
			{Type: "bot_command", Offset: 0, Length: 6},
		},
	}

	ee.Emit(context.Background(), events.OnMessage, &events.MessageEvent{
		Message: msg,
		Type:    messagetype.Text,
	})

	if !commandHandled {
		t.Fatal("command was not handled")
	}
	if afterParserCalled {
		t.Fatal("listener after parser was called, want skipped")
	}
}
