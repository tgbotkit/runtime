package handlers

import (
	"context"

	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/events"
	"github.com/tgbotkit/runtime/logger"
	"github.com/tgbotkit/runtime/messagetype"
)

// RegistryInterface defines the interface for managing event subscriptions.
type RegistryInterface interface {
	OnUpdate(handler UpdateHandler) eventemitter.UnsubscribeFunc
	OnMessage(handler MessageHandler) eventemitter.UnsubscribeFunc
	OnMessageType(t messagetype.MessageType, handler MessageHandler) eventemitter.UnsubscribeFunc
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

// OnCommand registers a handler for the OnCommand event.
func (r *Registry) OnCommand(handler CommandHandler) eventemitter.UnsubscribeFunc {
	r.l.Debugf("adding OnCommand handler: %T", handler)
	return eventemitter.On(r.em, events.OnCommand, func(ctx context.Context, event *events.CommandEvent) error {
		return handler(ctx, event)
	})
}