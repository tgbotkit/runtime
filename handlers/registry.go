package handlers

import (
	"context"
	"slices"

	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/events"
	"github.com/tgbotkit/runtime/logger"
	"github.com/tgbotkit/runtime/messagetype"
)

// HandlerRegistry defines the interface for managing event subscriptions.
type HandlerRegistry interface {
	OnUpdate(handler UpdateHandler) eventemitter.UnsubscribeFunc
	OnMessage(handler MessageHandler) eventemitter.UnsubscribeFunc
	OnMessageWithType(handler MessageHandler, types ...messagetype.MessageType) eventemitter.UnsubscribeFunc
	OnCommand(handler CommandHandler) eventemitter.UnsubscribeFunc
}

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
	r.l.Debug("adding OnUpdate handler: %T", handler)
	return eventemitter.On(r.em, events.OnUpdate, func(ctx context.Context, event *events.UpdateEvent) error {
		return handler(ctx, event)
	})
}

// OnMessage registers a handler for the OnMessageReceived event.
func (r *Registry) OnMessage(handler MessageHandler) eventemitter.UnsubscribeFunc {
	r.l.Debug("adding OnMessage handler: %T", handler)
	return eventemitter.On(r.em, events.OnMessage, func(ctx context.Context, event *events.MessageEvent) error {
		return handler(ctx, event)
	})
}

// OnMessageWithType registers a handler for the OnMessageReceived event with specific message types.
func (r *Registry) OnMessageWithType(handler MessageHandler, types ...messagetype.MessageType) eventemitter.UnsubscribeFunc {
	r.l.Debug("adding OnMessageWithType handler: %T", handler)
	return eventemitter.On(r.em, events.OnMessage, func(ctx context.Context, event *events.MessageEvent) error {
		if len(types) > 0 && !slices.Contains(types, event.Type) {
			return nil
		}

		return handler(ctx, event)
	})
}

// OnCommand registers a handler for the OnCommand event.
func (r *Registry) OnCommand(handler CommandHandler) eventemitter.UnsubscribeFunc {
	r.l.Debug("adding OnCommand handler: %T", handler)
	return eventemitter.On(r.em, events.OnCommand, func(ctx context.Context, event *events.CommandEvent) error {
		return handler(ctx, event)
	})
}
