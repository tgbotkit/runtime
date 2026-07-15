package listeners_test

import (
	"context"
	"testing"

	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/events"
	"github.com/tgbotkit/runtime/listeners"
	"github.com/tgbotkit/runtime/messagetype"
)

func TestClassifier(t *testing.T) {
	ee, err := eventemitter.NewSync(eventemitter.NewOptions())
	if err != nil {
		t.Fatalf("NewSync() unexpected error: %v", err)
	}
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
		if err != nil {
			t.Fatalf("Handle() unexpected error: %v", err)
		}
		if receivedEvent == nil {
			t.Fatal("receivedEvent is nil")
		}
		if receivedEvent.Type != messagetype.Text {
			t.Fatalf("received type=%q, want %q", receivedEvent.Type, messagetype.Text)
		}
		if receivedEvent.Message != update.Message {
			t.Fatalf("message pointer mismatch: got %p, want %p", receivedEvent.Message, update.Message)
		}
	})

	t.Run("classifies photo message", func(t *testing.T) {
		ee, err := eventemitter.NewSync(eventemitter.NewOptions())
		if err != nil {
			t.Fatalf("NewSync() unexpected error: %v", err)
		}
		classifier := listeners.Classifier(ee)

		var receivedEvent *events.MessageEvent
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

		err = classifier.Handle(context.Background(), &events.UpdateEvent{Update: update})
		if err != nil {
			t.Fatalf("Handle() unexpected error: %v", err)
		}
		if receivedEvent == nil {
			t.Fatal("receivedEvent is nil")
		}
		if receivedEvent.Type != messagetype.Photo {
			t.Fatalf("received type=%q, want %q", receivedEvent.Type, messagetype.Photo)
		}
	})

	t.Run("ignores updates without message", func(t *testing.T) {
		ee, err := eventemitter.NewSync(eventemitter.NewOptions())
		if err != nil {
			t.Fatalf("NewSync() unexpected error: %v", err)
		}
		classifier := listeners.Classifier(ee)

		var called bool
		ee.AddListener(events.OnMessage, eventemitter.ListenerFunc(func(_ context.Context, _ any) error {
			called = true
			return nil
		}))

		update := &client.Update{UpdateId: 3}
		err = classifier.Handle(context.Background(), &events.UpdateEvent{Update: update})
		if err != nil {
			t.Fatalf("Handle() unexpected error: %v", err)
		}
		if called {
			t.Fatal("OnMessage listener should not be called")
		}
	})

	t.Run("ignores invalid payload", func(t *testing.T) {
		err := classifier.Handle(context.Background(), "invalid-payload")
		if err != nil {
			t.Fatalf("Handle() unexpected error: %v", err)
		}
	})

	t.Run("classifies message types", func(t *testing.T) {
		tests := []struct {
			name     string
			message  *client.Message
			expected messagetype.MessageType
		}{
			{name: "Audio", message: &client.Message{Audio: &client.Audio{}}, expected: messagetype.Audio},
			{name: "Video", message: &client.Message{Video: &client.Video{}}, expected: messagetype.Video},
			{name: "Sticker", message: &client.Message{Sticker: &client.Sticker{}}, expected: messagetype.Sticker},
			{name: "LivePhoto", message: &client.Message{LivePhoto: &client.LivePhoto{}}, expected: messagetype.LivePhoto},
			{name: "Voice", message: &client.Message{Voice: &client.Voice{}}, expected: messagetype.Voice},
			{name: "Document", message: &client.Message{Document: &client.Document{}}, expected: messagetype.Document},
			{name: "Contact", message: &client.Message{Contact: &client.Contact{}}, expected: messagetype.Contact},
			{name: "Location", message: &client.Message{Location: &client.Location{}}, expected: messagetype.Location},
			{name: "NewChatMembers", message: &client.Message{NewChatMembers: &[]client.User{}}, expected: messagetype.NewChatMembers},
			{name: "LeftChatMember", message: &client.Message{LeftChatMember: &client.User{}}, expected: messagetype.LeftChatMember},
			{name: "ChatOwnerChanged", message: &client.Message{ChatOwnerChanged: &client.ChatOwnerChanged{}}, expected: messagetype.ChatOwnerChanged},
			{name: "ChatOwnerLeft", message: &client.Message{ChatOwnerLeft: &client.ChatOwnerLeft{}}, expected: messagetype.ChatOwnerLeft},
			{name: "ManagedBotCreated", message: &client.Message{ManagedBotCreated: &client.ManagedBotCreated{}}, expected: messagetype.ManagedBotCreated},
			{name: "PollOptionAdded", message: &client.Message{PollOptionAdded: &client.PollOptionAdded{}}, expected: messagetype.PollOptionAdded},
			{name: "PollOptionDeleted", message: &client.Message{PollOptionDeleted: &client.PollOptionDeleted{}}, expected: messagetype.PollOptionDeleted},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ee, err := eventemitter.NewSync(eventemitter.NewOptions())
				if err != nil {
					t.Fatalf("NewSync() unexpected error: %v", err)
				}
				classifier := listeners.Classifier(ee)

				var receivedEvent *events.MessageEvent
				ee.AddListener(events.OnMessage, eventemitter.ListenerFunc(func(_ context.Context, payload any) error {
					e, ok := payload.(*events.MessageEvent)
					if !ok {
						t.Fatalf("payload type=%T, want *events.MessageEvent", payload)
					}
					receivedEvent = e
					return nil
				}))

				update := &client.Update{
					UpdateId: 1,
					Message:  tt.message,
				}
				err = classifier.Handle(context.Background(), &events.UpdateEvent{Update: update})
				if err != nil {
					t.Fatalf("Handle() unexpected error: %v", err)
				}
				if receivedEvent == nil {
					t.Fatal("receivedEvent is nil")
				}
				if receivedEvent.Type != tt.expected {
					t.Fatalf("received type=%q, want %q", receivedEvent.Type, tt.expected)
				}
			})
		}
	})

	t.Run("emits non-message update events", func(t *testing.T) {
		tests := []struct {
			name       string
			update     *client.Update
			event      string
			assertType func(t *testing.T, payload any)
		}{
			{
				name:   "CallbackQuery",
				update: &client.Update{CallbackQuery: &client.CallbackQuery{}},
				event:  events.OnCallbackQuery,
				assertType: func(t *testing.T, payload any) {
					t.Helper()
					if _, ok := payload.(*events.CallbackQueryEvent); !ok {
						t.Fatalf("payload type=%T, want *events.CallbackQueryEvent", payload)
					}
				},
			},
			{
				name:   "InlineQuery",
				update: &client.Update{InlineQuery: &client.InlineQuery{}},
				event:  events.OnInlineQuery,
				assertType: func(t *testing.T, payload any) {
					t.Helper()
					if _, ok := payload.(*events.InlineQueryEvent); !ok {
						t.Fatalf("payload type=%T, want *events.InlineQueryEvent", payload)
					}
				},
			},
			{
				name:   "PollAnswer",
				update: &client.Update{PollAnswer: &client.PollAnswer{}},
				event:  events.OnPollAnswer,
				assertType: func(t *testing.T, payload any) {
					t.Helper()
					if _, ok := payload.(*events.PollAnswerEvent); !ok {
						t.Fatalf("payload type=%T, want *events.PollAnswerEvent", payload)
					}
				},
			},
			{
				name:   "ChatMember",
				update: &client.Update{ChatMember: &client.ChatMemberUpdated{}},
				event:  events.OnChatMember,
				assertType: func(t *testing.T, payload any) {
					t.Helper()
					if _, ok := payload.(*events.ChatMemberEvent); !ok {
						t.Fatalf("payload type=%T, want *events.ChatMemberEvent", payload)
					}
				},
			},
			{
				name:   "MessageReaction",
				update: &client.Update{MessageReaction: &client.MessageReactionUpdated{}},
				event:  events.OnMessageReaction,
				assertType: func(t *testing.T, payload any) {
					t.Helper()
					if _, ok := payload.(*events.MessageReactionEvent); !ok {
						t.Fatalf("payload type=%T, want *events.MessageReactionEvent", payload)
					}
				},
			},
			{
				name:   "Subscription",
				update: &client.Update{Subscription: &client.BotSubscriptionUpdated{}},
				event:  events.OnSubscription,
				assertType: func(t *testing.T, payload any) {
					t.Helper()
					if _, ok := payload.(*events.SubscriptionEvent); !ok {
						t.Fatalf("payload type=%T, want *events.SubscriptionEvent", payload)
					}
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ee, err := eventemitter.NewSync(eventemitter.NewOptions())
				if err != nil {
					t.Fatalf("NewSync() unexpected error: %v", err)
				}
				classifier := listeners.Classifier(ee)

				var called bool
				ee.AddListener(tt.event, eventemitter.ListenerFunc(func(_ context.Context, payload any) error {
					called = true
					tt.assertType(t, payload)

					return nil
				}))

				err = classifier.Handle(context.Background(), &events.UpdateEvent{Update: tt.update})
				if err != nil {
					t.Fatalf("Handle() unexpected error: %v", err)
				}
				if !called {
					t.Fatalf("%s listener was not called", tt.event)
				}
			})
		}
	})

	t.Run("emits message-like update events", func(t *testing.T) {
		text := "edited"
		tests := []struct {
			name   string
			update *client.Update
			event  string
		}{
			{name: "EditedMessage", update: &client.Update{EditedMessage: &client.Message{Text: &text}}, event: events.OnEditedMessage},
			{name: "ChannelPost", update: &client.Update{ChannelPost: &client.Message{Text: &text}}, event: events.OnChannelPost},
			{name: "BusinessMessage", update: &client.Update{BusinessMessage: &client.Message{Text: &text}}, event: events.OnBusinessMessage},
			{name: "GuestMessage", update: &client.Update{GuestMessage: &client.Message{Text: &text}}, event: events.OnGuestMessage},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ee, err := eventemitter.NewSync(eventemitter.NewOptions())
				if err != nil {
					t.Fatalf("NewSync() unexpected error: %v", err)
				}
				classifier := listeners.Classifier(ee)

				var receivedEvent *events.MessageEvent
				ee.AddListener(tt.event, eventemitter.ListenerFunc(func(_ context.Context, payload any) error {
					event, ok := payload.(*events.MessageEvent)
					if !ok {
						t.Fatalf("payload type=%T, want *events.MessageEvent", payload)
					}
					receivedEvent = event

					return nil
				}))

				err = classifier.Handle(context.Background(), &events.UpdateEvent{Update: tt.update})
				if err != nil {
					t.Fatalf("Handle() unexpected error: %v", err)
				}
				if receivedEvent == nil {
					t.Fatal("receivedEvent is nil")
				}
				if receivedEvent.Type != messagetype.Text {
					t.Fatalf("received type=%q, want %q", receivedEvent.Type, messagetype.Text)
				}
			})
		}
	})
}
