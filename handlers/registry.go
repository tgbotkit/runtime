package handlers

import (
	"context"

	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/events"
)

// HandlerRegistry defines the interface for managing event subscriptions.
type HandlerRegistry interface {
	OnUpdate(handler UpdateHandler) eventemitter.UnsubscribeFunc
	OnMessage(handler MessageHandler) eventemitter.UnsubscribeFunc
	OnCommand(handler CommandHandler) eventemitter.UnsubscribeFunc
}

// Registry manages the subscription of handlers to events.
type Registry struct {
	em eventemitter.EventEmitter
}

// NewRegistry creates a new Registry.
func NewRegistry(em eventemitter.EventEmitter) *Registry {
	return &Registry{
		em: em,
	}
}

// OnUpdate registers a handler for the OnUpdateReceived event.
func (r *Registry) OnUpdate(handler UpdateHandler) eventemitter.UnsubscribeFunc {
	return eventemitter.On(r.em, events.OnUpdate, func(ctx context.Context, event *events.UpdateEvent) error {
		return handler(ctx, event)
	})
}

// OnMessage registers a handler for the OnMessageReceived event.
func (r *Registry) OnMessage(handler MessageHandler) eventemitter.UnsubscribeFunc {
	return eventemitter.On(r.em, events.OnMessage, func(ctx context.Context, event *events.MessageEvent) error {
		return handler(ctx, event)
	})
}

// OnCommand registers a handler for the OnCommand event.
func (r *Registry) OnCommand(handler CommandHandler) eventemitter.UnsubscribeFunc {
	return eventemitter.On(r.em, events.OnCommand, func(ctx context.Context, event *events.CommandEvent) error {
		return handler(ctx, event)
	})
}
