package listeners_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/events"
	"github.com/tgbotkit/runtime/listeners"
	"github.com/tgbotkit/runtime/messagetype"
)

func TestCommandParser(t *testing.T) {
	ee, _ := eventemitter.NewSync(eventemitter.NewOptions())
	botName := "mybot"
	parser := listeners.CommandParser(ee, botName)

	t.Run("parses simple command", func(t *testing.T) {
		var receivedEvent *events.CommandEvent
		ee.AddListener(events.OnCommand, eventemitter.ListenerFunc(func(ctx context.Context, payload any) error {
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
		assert.NoError(t, err)
		assert.NotNil(t, receivedEvent)
		assert.Equal(t, "start", receivedEvent.Command)
		assert.Equal(t, "", receivedEvent.Args)
	})

	t.Run("parses command with args", func(t *testing.T) {
		// Reset listener
		ee, _ := eventemitter.NewSync(eventemitter.NewOptions())
		parser := listeners.CommandParser(ee, botName)
		var receivedEvent *events.CommandEvent
		ee.AddListener(events.OnCommand, eventemitter.ListenerFunc(func(ctx context.Context, payload any) error {
			receivedEvent = payload.(*events.CommandEvent)
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

		err := parser.Handle(context.Background(), event)
		assert.NoError(t, err)
		assert.NotNil(t, receivedEvent)
		assert.Equal(t, "echo", receivedEvent.Command)
		assert.Equal(t, "hello world", receivedEvent.Args)
	})

	t.Run("parses command with bot mention", func(t *testing.T) {
		ee, _ := eventemitter.NewSync(eventemitter.NewOptions())
		parser := listeners.CommandParser(ee, botName)
		var receivedEvent *events.CommandEvent
		ee.AddListener(events.OnCommand, eventemitter.ListenerFunc(func(ctx context.Context, payload any) error {
			receivedEvent = payload.(*events.CommandEvent)
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

		err := parser.Handle(context.Background(), event)
		assert.NoError(t, err)
		assert.NotNil(t, receivedEvent)
		assert.Equal(t, "start", receivedEvent.Command)
		assert.Equal(t, "args", receivedEvent.Args)
	})

	t.Run("ignores command for other bot", func(t *testing.T) {
		ee, _ := eventemitter.NewSync(eventemitter.NewOptions())
		parser := listeners.CommandParser(ee, botName)
		var called bool
		ee.AddListener(events.OnCommand, eventemitter.ListenerFunc(func(ctx context.Context, payload any) error {
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

		err := parser.Handle(context.Background(), event)
		assert.NoError(t, err)
		assert.False(t, called)
	})

	t.Run("ignores non-text message", func(t *testing.T) {
		ee, _ := eventemitter.NewSync(eventemitter.NewOptions())
		parser := listeners.CommandParser(ee, botName)
		
		event := &events.MessageEvent{Type: messagetype.Photo}
		err := parser.Handle(context.Background(), event)
		assert.NoError(t, err)
	})
	
	t.Run("ignores invalid payload", func(t *testing.T) {
		err := parser.Handle(context.Background(), "string")
		assert.NoError(t, err)
	})
}
