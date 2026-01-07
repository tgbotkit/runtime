package listeners

import (
	"context"
	"strings"

	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/events"
	"github.com/tgbotkit/runtime/messagetype"
)

func CommandParser(emitter eventemitter.EventEmitter, botName string) eventemitter.Listener {
	return eventemitter.ListenerFunc(func(ctx context.Context, payload any) error {
		// We only care about MessageEvent with text.
		if event, ok := payload.(*events.MessageEvent); ok && event.Type == messagetype.Text {
			if err := parseCommand(ctx, emitter, event, botName); err != nil {
				// Stop propagation if a command was handled.
				if err == eventemitter.ErrBreak {
					return nil // Don't propagate further, but don't treat as an error.
				}
				return err
			}
		}
		return nil
	})
}

// parseCommand checks a message for a bot command and emits a CommandEvent if one is found.
func parseCommand(ctx context.Context, emitter eventemitter.EventEmitter, event *events.MessageEvent, botName string) error {
	if event.Message == nil || event.Message.Text == nil || event.Message.Entities == nil {
		return nil
	}
	text := *event.Message.Text

	for _, entity := range *event.Message.Entities {
		if entity.Type != "bot_command" {
			continue
		}

		commandText := sliceText(text, entity.Offset, entity.Length)
		parts := strings.Split(commandText, "@")
		if len(parts) > 1 {
			if botName == "" || parts[1] != botName {
				continue // Command is for another bot
			}
			commandText = parts[0]
		}

		commandText = strings.TrimPrefix(commandText, "/")
		args := sliceTextFrom(text, entity.Offset+entity.Length)
		args = strings.TrimLeft(args, " ")

		emitter.Emit(ctx, events.OnCommand, &events.CommandEvent{
			Bot:     event.Bot,
			Message: event.Message,
			Command: commandText,
			Args:    args,
		})
		return eventemitter.ErrBreak // Stop further processing of this message
	}

	return nil
}

func sliceText(text string, offset int, length int) string {
	runes := []rune(text)
	if offset < 0 {
		offset = 0
	}
	if offset >= len(runes) {
		return ""
	}
	end := offset + length
	if end > len(runes) {
		end = len(runes)
	}
	return string(runes[offset:end])
}

func sliceTextFrom(text string, offset int) string {
	runes := []rune(text)
	if offset < 0 {
		offset = 0
	}
	if offset >= len(runes) {
		return ""
	}
	return string(runes[offset:])
}
