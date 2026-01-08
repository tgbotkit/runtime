// Package listeners provides core event listeners for the bot, such as update classification and command parsing.
package listeners

import (
	"context"

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

	if update.Message != nil {
		msgType := messagetype.Detect(update.Message)
		messageEvent := &events.MessageEvent{
			Message: update.Message,
			Type:    msgType,
		}
		emitter.Emit(ctx, events.OnMessage, messageEvent)
	}
}
