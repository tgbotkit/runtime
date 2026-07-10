package listeners

import (
	"context"
	"strings"
	"unicode/utf16"

	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/events"
)

// CommandParser returns a listener that detects bot commands in messages and emits OnCommand events.
func CommandParser(emitter eventemitter.EventEmitter, botName string) eventemitter.Listener {
	return eventemitter.ListenerFunc(func(ctx context.Context, payload any) error {
		if event, ok := payload.(*events.MessageEvent); ok {
			if err := parseCommand(ctx, emitter, event, botName); err != nil {
				return err
			}
		}

		return nil
	})
}

// parseCommand checks a message for a bot command and emits a CommandEvent if one is found.
func parseCommand(
	ctx context.Context,
	emitter eventemitter.EventEmitter,
	event *events.MessageEvent,
	botName string,
) error {
	if event.Message == nil {
		return nil
	}

	text, entities, ok := commandSource(event.Message)
	if !ok {
		return nil
	}

	for _, entity := range entities {
		if entity.Type != "bot_command" {
			continue
		}
		if entity.Offset != 0 {
			continue
		}

		commandText := sliceText(text, entity.Offset, entity.Length)

		command, mention, hasMention := strings.Cut(commandText, "@")
		if hasMention {
			if botName == "" || !strings.EqualFold(mention, strings.TrimPrefix(botName, "@")) {
				continue // Command is for another bot
			}

			commandText = command
		}

		commandText = strings.TrimPrefix(commandText, "/")
		args := sliceTextFrom(text, entity.Offset+entity.Length)
		args = strings.TrimLeft(args, " ")

		emitter.Emit(ctx, events.OnCommand, &events.CommandEvent{
			Message: event.Message,
			Command: commandText,
			Args:    args,
		})

		return eventemitter.ErrBreak // Stop further processing of this message
	}

	return nil
}

func commandSource(message *client.Message) (string, []client.MessageEntity, bool) {
	if message.Text != nil && message.Entities != nil {
		return *message.Text, *message.Entities, true
	}

	if message.Caption != nil && message.CaptionEntities != nil {
		return *message.Caption, *message.CaptionEntities, true
	}

	return "", nil, false
}

func sliceText(text string, offset int, length int) string {
	runes := []rune(text)

	start := runeIndexFromUTF16Offset(runes, offset)

	end := runeIndexFromUTF16Offset(runes, offset+length)
	if end < start {
		end = start
	}

	return string(runes[start:end])
}

func sliceTextFrom(text string, offset int) string {
	runes := []rune(text)

	start := runeIndexFromUTF16Offset(runes, offset)

	return string(runes[start:])
}

func runeIndexFromUTF16Offset(runes []rune, offset int) int {
	if offset <= 0 {
		return 0
	}

	units := 0
	for i, r := range runes {
		if units >= offset {
			return i
		}

		units += utf16.RuneLen(r)
	}

	return len(runes)
}
