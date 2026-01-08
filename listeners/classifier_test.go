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

func TestClassifier(t *testing.T) {
	ee, _ := eventemitter.NewSync(eventemitter.NewOptions())
	classifier := listeners.Classifier(ee)

	t.Run("classifies text message", func(t *testing.T) {
		var receivedEvent *events.MessageEvent
		ee.AddListener(events.OnMessage, eventemitter.ListenerFunc(func(_ context.Context, payload any) error {
			if e, ok := payload.(*events.MessageEvent); ok {
				receivedEvent = e
			}
			return nil
		}))

		text := "hello"
		update := &client.Update{
			UpdateId: 1,
			Message: &client.Message{
				MessageId: 100,
				Text:      &text,
			},
		}

		err := classifier.Handle(context.Background(), &events.UpdateEvent{Update: update})
		assert.NoError(t, err)
		assert.NotNil(t, receivedEvent)
		assert.Equal(t, messagetype.Text, receivedEvent.Type)
		assert.Equal(t, update.Message, receivedEvent.Message)
	})

	t.Run("classifies photo message", func(t *testing.T) {
		var receivedEvent *events.MessageEvent
		// Clear previous listeners or use new emitter? Reusing ee is fine if we add new listener or reset.
		// Easier to just use checks.
		
		// Let's create a new emitter for clean state
		ee, _ := eventemitter.NewSync(eventemitter.NewOptions())
		classifier := listeners.Classifier(ee)
		
		ee.AddListener(events.OnMessage, eventemitter.ListenerFunc(func(_ context.Context, payload any) error {
			if e, ok := payload.(*events.MessageEvent); ok {
				receivedEvent = e
			}
			return nil
		}))

		update := &client.Update{
			UpdateId: 2,
			Message: &client.Message{
				MessageId: 101,
				Photo:     &[]client.PhotoSize{},
			},
		}

		err := classifier.Handle(context.Background(), &events.UpdateEvent{Update: update})
		assert.NoError(t, err)
		assert.NotNil(t, receivedEvent)
		assert.Equal(t, messagetype.Photo, receivedEvent.Type)
	})

	t.Run("ignores updates without message", func(t *testing.T) {
		ee, _ := eventemitter.NewSync(eventemitter.NewOptions())
		classifier := listeners.Classifier(ee)
		
		var called bool
		ee.AddListener(events.OnMessage, eventemitter.ListenerFunc(func(_ context.Context, _ any) error {
			called = true
			return nil
		}))

		update := &client.Update{
			UpdateId: 3,
			// No message
		}

		err := classifier.Handle(context.Background(), &events.UpdateEvent{Update: update})
		assert.NoError(t, err)
		assert.False(t, called)
	})
	
	t.Run("ignores invalid payload", func(t *testing.T) {
		err := classifier.Handle(context.Background(), "invalid-payload")
		assert.NoError(t, err)
	})

	t.Run("classifies message types", func(t *testing.T) {
		tests := []struct {
			name     string
			message  *client.Message
			expected messagetype.MessageType
		}{
			{
				name:     "Audio",
				message:  &client.Message{Audio: &client.Audio{}},
				expected: messagetype.Audio,
			},
			{
				name:     "Video",
				message:  &client.Message{Video: &client.Video{}},
				expected: messagetype.Video,
			},
			{
				name:     "Sticker",
				message:  &client.Message{Sticker: &client.Sticker{}},
				expected: messagetype.Sticker,
			},
			{
				name:     "Voice",
				message:  &client.Message{Voice: &client.Voice{}},
				expected: messagetype.Voice,
			},
			{
				name:     "Document",
				message:  &client.Message{Document: &client.Document{}},
				expected: messagetype.Document,
			},
			{
				name:     "Contact",
				message:  &client.Message{Contact: &client.Contact{}},
				expected: messagetype.Contact,
			},
			{
				name:     "Location",
				message:  &client.Message{Location: &client.Location{}},
				expected: messagetype.Location,
			},
			{
				name:     "NewChatMembers",
				message:  &client.Message{NewChatMembers: &[]client.User{}},
				expected: messagetype.NewChatMembers,
			},
			{
				name:     "LeftChatMember",
				message:  &client.Message{LeftChatMember: &client.User{}},
				expected: messagetype.LeftChatMember,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ee, _ := eventemitter.NewSync(eventemitter.NewOptions())
				classifier := listeners.Classifier(ee)
				var receivedEvent *events.MessageEvent
				ee.AddListener(events.OnMessage, eventemitter.ListenerFunc(func(_ context.Context, payload any) error {
					receivedEvent = payload.(*events.MessageEvent)
					return nil
				}))

				update := &client.Update{
					UpdateId: 1,
					Message:  tt.message,
				}
				err := classifier.Handle(context.Background(), &events.UpdateEvent{Update: update})
				assert.NoError(t, err)
				assert.NotNil(t, receivedEvent)
				assert.Equal(t, tt.expected, receivedEvent.Type)
			})
		}
	})
}