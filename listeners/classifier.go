// Package listeners provides core event listeners for the bot, such as update classification and command parsing.
package listeners

import (
	"context"

	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/events"
	"github.com/tgbotkit/runtime/messagetype"
)

// Classifier returns a listener that analyzes incoming updates and emits more specific events
// based on the update content.
func Classifier(emitter eventemitter.EventEmitter) eventemitter.Listener {
	return eventemitter.ListenerFunc(func(ctx context.Context, payload any) error {
		if event, ok := payload.(*events.UpdateEvent); ok {
			if event != nil && event.Update != nil {
				classifyUpdate(ctx, emitter, event)
			}
		}

		return nil
	})
}

// classifyUpdate inspects the update and emits corresponding events.
func classifyUpdate(ctx context.Context, emitter eventemitter.EventEmitter, event *events.UpdateEvent) {
	update := event.Update

	classifyMessages(ctx, emitter, update)
	classifyQueries(ctx, emitter, update)
	classifyChatUpdates(ctx, emitter, update)
	classifyBusinessUpdates(ctx, emitter, update)
}

func classifyMessages(ctx context.Context, emitter eventemitter.EventEmitter, update *client.Update) {
	if update.Message != nil {
		emitMessage(ctx, emitter, events.OnMessage, update.Message)
	}

	if update.EditedMessage != nil {
		emitMessage(ctx, emitter, events.OnEditedMessage, update.EditedMessage)
	}

	if update.ChannelPost != nil {
		emitMessage(ctx, emitter, events.OnChannelPost, update.ChannelPost)
	}

	if update.EditedChannelPost != nil {
		emitMessage(ctx, emitter, events.OnEditedChannelPost, update.EditedChannelPost)
	}

	if update.BusinessMessage != nil {
		emitMessage(ctx, emitter, events.OnBusinessMessage, update.BusinessMessage)
	}

	if update.EditedBusinessMessage != nil {
		emitMessage(ctx, emitter, events.OnEditedBusinessMessage, update.EditedBusinessMessage)
	}

	if update.GuestMessage != nil {
		emitMessage(ctx, emitter, events.OnGuestMessage, update.GuestMessage)
	}
}

func classifyQueries(ctx context.Context, emitter eventemitter.EventEmitter, update *client.Update) {
	if update.CallbackQuery != nil {
		emitter.Emit(
			ctx,
			events.OnCallbackQuery,
			&events.CallbackQueryEvent{CallbackQuery: update.CallbackQuery},
		)
	}

	if update.InlineQuery != nil {
		emitter.Emit(
			ctx,
			events.OnInlineQuery,
			&events.InlineQueryEvent{InlineQuery: update.InlineQuery},
		)
	}

	if update.ChosenInlineResult != nil {
		emitter.Emit(ctx, events.OnChosenInlineResult, &events.ChosenInlineResultEvent{
			ChosenInlineResult: update.ChosenInlineResult,
		})
	}

	if update.ShippingQuery != nil {
		emitter.Emit(
			ctx,
			events.OnShippingQuery,
			&events.ShippingQueryEvent{ShippingQuery: update.ShippingQuery},
		)
	}

	if update.PreCheckoutQuery != nil {
		emitter.Emit(ctx, events.OnPreCheckoutQuery, &events.PreCheckoutQueryEvent{
			PreCheckoutQuery: update.PreCheckoutQuery,
		})
	}
}

func classifyChatUpdates(ctx context.Context, emitter eventemitter.EventEmitter, update *client.Update) {
	if update.Poll != nil {
		emitter.Emit(ctx, events.OnPoll, &events.PollEvent{Poll: update.Poll})
	}

	if update.PollAnswer != nil {
		emitter.Emit(ctx, events.OnPollAnswer, &events.PollAnswerEvent{PollAnswer: update.PollAnswer})
	}

	if update.ChatMember != nil {
		emitter.Emit(ctx, events.OnChatMember, &events.ChatMemberEvent{ChatMember: update.ChatMember})
	}

	if update.MyChatMember != nil {
		emitter.Emit(ctx, events.OnMyChatMember, &events.ChatMemberEvent{ChatMember: update.MyChatMember})
	}

	if update.ChatJoinRequest != nil {
		emitter.Emit(ctx, events.OnChatJoinRequest, &events.ChatJoinRequestEvent{
			ChatJoinRequest: update.ChatJoinRequest,
		})
	}

	if update.ChatBoost != nil {
		emitter.Emit(ctx, events.OnChatBoost, &events.ChatBoostEvent{ChatBoost: update.ChatBoost})
	}

	if update.RemovedChatBoost != nil {
		emitter.Emit(ctx, events.OnRemovedChatBoost, &events.RemovedChatBoostEvent{
			RemovedChatBoost: update.RemovedChatBoost,
		})
	}

	if update.MessageReaction != nil {
		emitter.Emit(ctx, events.OnMessageReaction, &events.MessageReactionEvent{
			MessageReaction: update.MessageReaction,
		})
	}

	if update.MessageReactionCount != nil {
		emitter.Emit(ctx, events.OnMessageReactionCount, &events.MessageReactionCountEvent{
			MessageReactionCount: update.MessageReactionCount,
		})
	}
}

func classifyBusinessUpdates(ctx context.Context, emitter eventemitter.EventEmitter, update *client.Update) {
	for _, candidate := range []struct {
		event   string
		payload any
	}{
		{
			event: events.OnBusinessConnection,
			payload: &events.BusinessConnectionEvent{
				BusinessConnection: update.BusinessConnection,
			},
		},
		{
			event: events.OnDeletedBusinessMessages,
			payload: &events.DeletedBusinessMessagesEvent{
				DeletedBusinessMessages: update.DeletedBusinessMessages,
			},
		},
		{
			event: events.OnPurchasedPaidMedia,
			payload: &events.PurchasedPaidMediaEvent{
				PurchasedPaidMedia: update.PurchasedPaidMedia,
			},
		},
		{
			event:   events.OnManagedBot,
			payload: &events.ManagedBotEvent{ManagedBot: update.ManagedBot},
		},
		{
			event: events.OnSubscription,
			payload: &events.SubscriptionEvent{
				Subscription: update.Subscription,
			},
		},
	} {
		if !isNilPayload(candidate.payload) {
			emitter.Emit(ctx, candidate.event, candidate.payload)
		}
	}
}

func emitMessage(ctx context.Context, emitter eventemitter.EventEmitter, event string, message *client.Message) {
	emitter.Emit(ctx, event, &events.MessageEvent{
		Message: message,
		Type:    messagetype.Detect(message),
	})
}

func isNilPayload(payload any) bool {
	switch v := payload.(type) {
	case *events.BusinessConnectionEvent:
		return v.BusinessConnection == nil
	case *events.DeletedBusinessMessagesEvent:
		return v.DeletedBusinessMessages == nil
	case *events.PurchasedPaidMediaEvent:
		return v.PurchasedPaidMedia == nil
	case *events.ManagedBotEvent:
		return v.ManagedBot == nil
	case *events.SubscriptionEvent:
		return v.Subscription == nil
	default:
		return payload == nil
	}
}
