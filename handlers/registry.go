package handlers

import (
	"context"

	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/events"
	"github.com/tgbotkit/runtime/logger"
	"github.com/tgbotkit/runtime/messagetype"
)

// Registry manages the subscription of handlers to events.
type Registry struct {
	em eventemitter.EventEmitter
	l  logger.Logger
}

// NewRegistry creates a new Registry.
func NewRegistry(em eventemitter.EventEmitter, l logger.Logger) *Registry {
	return &Registry{
		em: em,
		l:  l,
	}
}

// OnUpdate registers a handler for the OnUpdateReceived event.
func (r *Registry) OnUpdate(handler UpdateHandler) eventemitter.UnsubscribeFunc {
	r.l.Debugf("adding OnUpdate handler: %T", handler)

	return eventemitter.On(r.em, events.OnUpdate, func(ctx context.Context, event *events.UpdateEvent) error {
		return handler(ctx, event)
	})
}

// OnMessage registers a handler for the OnMessageReceived event.
func (r *Registry) OnMessage(handler MessageHandler) eventemitter.UnsubscribeFunc {
	r.l.Debugf("adding OnMessage handler: %T", handler)

	return eventemitter.On(r.em, events.OnMessage, func(ctx context.Context, event *events.MessageEvent) error {
		return handler(ctx, event)
	})
}

// OnMessageType registers a handler for the OnMessageReceived event with a specific message type.
func (r *Registry) OnMessageType(t messagetype.MessageType, handler MessageHandler) eventemitter.UnsubscribeFunc {
	r.l.Debugf("adding OnMessageType handler: %T for type %s", handler, t)

	return eventemitter.On(r.em, events.OnMessage, func(ctx context.Context, event *events.MessageEvent) error {
		if event.Type != t {
			return nil
		}

		return handler(ctx, event)
	})
}

// OnEditedMessage registers a handler for edited messages.
func (r *Registry) OnEditedMessage(handler MessageHandler) eventemitter.UnsubscribeFunc {
	return r.onMessageEvent(events.OnEditedMessage, "OnEditedMessage", handler)
}

// OnChannelPost registers a handler for channel posts.
func (r *Registry) OnChannelPost(handler MessageHandler) eventemitter.UnsubscribeFunc {
	return r.onMessageEvent(events.OnChannelPost, "OnChannelPost", handler)
}

// OnEditedChannelPost registers a handler for edited channel posts.
func (r *Registry) OnEditedChannelPost(handler MessageHandler) eventemitter.UnsubscribeFunc {
	return r.onMessageEvent(events.OnEditedChannelPost, "OnEditedChannelPost", handler)
}

// OnBusinessMessage registers a handler for business messages.
func (r *Registry) OnBusinessMessage(handler MessageHandler) eventemitter.UnsubscribeFunc {
	return r.onMessageEvent(events.OnBusinessMessage, "OnBusinessMessage", handler)
}

// OnEditedBusinessMessage registers a handler for edited business messages.
func (r *Registry) OnEditedBusinessMessage(handler MessageHandler) eventemitter.UnsubscribeFunc {
	return r.onMessageEvent(events.OnEditedBusinessMessage, "OnEditedBusinessMessage", handler)
}

// OnGuestMessage registers a handler for guest messages.
func (r *Registry) OnGuestMessage(handler MessageHandler) eventemitter.UnsubscribeFunc {
	return r.onMessageEvent(events.OnGuestMessage, "OnGuestMessage", handler)
}

// OnCommand registers a handler for the OnCommand event.
func (r *Registry) OnCommand(handler CommandHandler) eventemitter.UnsubscribeFunc {
	r.l.Debugf("adding OnCommand handler: %T", handler)

	return eventemitter.On(r.em, events.OnCommand, func(ctx context.Context, event *events.CommandEvent) error {
		return handler(ctx, event)
	})
}

// OnCallbackQuery registers a handler for callback query events.
func (r *Registry) OnCallbackQuery(handler CallbackQueryHandler) eventemitter.UnsubscribeFunc {
	return onEvent(r, events.OnCallbackQuery, "OnCallbackQuery", handler)
}

// OnInlineQuery registers a handler for inline query events.
func (r *Registry) OnInlineQuery(handler InlineQueryHandler) eventemitter.UnsubscribeFunc {
	return onEvent(r, events.OnInlineQuery, "OnInlineQuery", handler)
}

// OnChosenInlineResult registers a handler for chosen inline result events.
func (r *Registry) OnChosenInlineResult(handler ChosenInlineResultHandler) eventemitter.UnsubscribeFunc {
	return onEvent(r, events.OnChosenInlineResult, "OnChosenInlineResult", handler)
}

// OnShippingQuery registers a handler for shipping query events.
func (r *Registry) OnShippingQuery(handler ShippingQueryHandler) eventemitter.UnsubscribeFunc {
	return onEvent(r, events.OnShippingQuery, "OnShippingQuery", handler)
}

// OnPreCheckoutQuery registers a handler for pre-checkout query events.
func (r *Registry) OnPreCheckoutQuery(handler PreCheckoutQueryHandler) eventemitter.UnsubscribeFunc {
	return onEvent(r, events.OnPreCheckoutQuery, "OnPreCheckoutQuery", handler)
}

// OnPoll registers a handler for poll events.
func (r *Registry) OnPoll(handler PollHandler) eventemitter.UnsubscribeFunc {
	return onEvent(r, events.OnPoll, "OnPoll", handler)
}

// OnPollAnswer registers a handler for poll answer events.
func (r *Registry) OnPollAnswer(handler PollAnswerHandler) eventemitter.UnsubscribeFunc {
	return onEvent(r, events.OnPollAnswer, "OnPollAnswer", handler)
}

// OnChatMember registers a handler for chat member events.
func (r *Registry) OnChatMember(handler ChatMemberHandler) eventemitter.UnsubscribeFunc {
	return onEvent(r, events.OnChatMember, "OnChatMember", handler)
}

// OnMyChatMember registers a handler for the bot's chat member events.
func (r *Registry) OnMyChatMember(handler ChatMemberHandler) eventemitter.UnsubscribeFunc {
	return onEvent(r, events.OnMyChatMember, "OnMyChatMember", handler)
}

// OnChatJoinRequest registers a handler for chat join request events.
func (r *Registry) OnChatJoinRequest(handler ChatJoinRequestHandler) eventemitter.UnsubscribeFunc {
	return onEvent(r, events.OnChatJoinRequest, "OnChatJoinRequest", handler)
}

// OnChatBoost registers a handler for chat boost events.
func (r *Registry) OnChatBoost(handler ChatBoostHandler) eventemitter.UnsubscribeFunc {
	return onEvent(r, events.OnChatBoost, "OnChatBoost", handler)
}

// OnRemovedChatBoost registers a handler for removed chat boost events.
func (r *Registry) OnRemovedChatBoost(handler RemovedChatBoostHandler) eventemitter.UnsubscribeFunc {
	return onEvent(r, events.OnRemovedChatBoost, "OnRemovedChatBoost", handler)
}

// OnMessageReaction registers a handler for message reaction events.
func (r *Registry) OnMessageReaction(handler MessageReactionHandler) eventemitter.UnsubscribeFunc {
	return onEvent(r, events.OnMessageReaction, "OnMessageReaction", handler)
}

// OnMessageReactionCount registers a handler for message reaction count events.
func (r *Registry) OnMessageReactionCount(handler MessageReactionCountHandler) eventemitter.UnsubscribeFunc {
	return onEvent(r, events.OnMessageReactionCount, "OnMessageReactionCount", handler)
}

// OnBusinessConnection registers a handler for business connection events.
func (r *Registry) OnBusinessConnection(handler BusinessConnectionHandler) eventemitter.UnsubscribeFunc {
	return onEvent(r, events.OnBusinessConnection, "OnBusinessConnection", handler)
}

// OnDeletedBusinessMessages registers a handler for deleted business messages events.
func (r *Registry) OnDeletedBusinessMessages(handler DeletedBusinessMessagesHandler) eventemitter.UnsubscribeFunc {
	return onEvent(r, events.OnDeletedBusinessMessages, "OnDeletedBusinessMessages", handler)
}

// OnPurchasedPaidMedia registers a handler for purchased paid media events.
func (r *Registry) OnPurchasedPaidMedia(handler PurchasedPaidMediaHandler) eventemitter.UnsubscribeFunc {
	return onEvent(r, events.OnPurchasedPaidMedia, "OnPurchasedPaidMedia", handler)
}

// OnManagedBot registers a handler for managed bot events.
func (r *Registry) OnManagedBot(handler ManagedBotHandler) eventemitter.UnsubscribeFunc {
	return onEvent(r, events.OnManagedBot, "OnManagedBot", handler)
}

func (r *Registry) onMessageEvent(
	event string,
	name string,
	handler MessageHandler,
) eventemitter.UnsubscribeFunc {
	r.l.Debugf("adding %s handler: %T", name, handler)

	return eventemitter.On(r.em, event, func(ctx context.Context, event *events.MessageEvent) error {
		return handler(ctx, event)
	})
}

func onEvent[E any, H ~func(context.Context, *E) error](
	r *Registry,
	event string,
	name string,
	handler H,
) eventemitter.UnsubscribeFunc {
	r.l.Debugf("adding %s handler: %T", name, handler)

	return eventemitter.On(r.em, event, func(ctx context.Context, event *E) error {
		return handler(ctx, event)
	})
}
