package commandparser

import (
	"context"
	"fmt"
	"strings"

	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/events"
)

// CommandParser is a component that handles bot commands.
type CommandParser struct {
	botName string
}

func New() (*CommandParser, error) {
	return &CommandParser{}, nil
}

func (h *CommandParser) Subscribe(bot events.BotContext) {
	ee := bot.EventEmitter()
	eventemitter.On[events.BotEvent](ee, events.OnBeforeStart, h.onBeforeStart)
	eventemitter.On[events.MessageEvent](ee, events.OnTextMessageReceived, h.onTextMessageReceived)
}

func (h *CommandParser) onBeforeStart(ctx context.Context, event *events.BotEvent) error {
	botName, err := getBotName(ctx, event.Bot.Client())
	if err != nil {
		return fmt.Errorf("unable to load bot name: %w", err)
	}
	h.botName = botName
	return nil
}

func (h *CommandParser) onTextMessageReceived(ctx context.Context, event *events.MessageEvent) error {
	if event == nil || event.Message == nil || event.Message.Text == nil {
		return nil
	}
	message := event.Message

	if message.Entities == nil {
		return nil
	}
	text := *message.Text

	for _, entity := range *message.Entities {
		if entity.Type != "bot_command" {
			continue
		}
		commandText := sliceText(text, entity.Offset, entity.Length)
		parts := strings.Split(commandText, "@")
		if len(parts) > 1 {
			if h.botName == "" || parts[1] != h.botName {
				return nil
			}
			commandText = parts[0]
		}
		commandText = strings.TrimPrefix(commandText, "/")
		args := sliceTextFrom(text, entity.Offset+entity.Length)
		args = strings.TrimLeft(args, " ")

		event.Bot.EventEmitter().Emit(ctx, events.OnCommand, &events.CommandEvent{Bot: event.Bot, Message: event.Message, Command: commandText, Args: args})
		return eventemitter.ErrBreak // stop propagation
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

func getBotName(ctx context.Context, client client.ClientWithResponsesInterface) (string, error) {
	resp, err := client.GetMeWithResponse(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get bot info: %w", err)
	}

	if resp.JSON200 == nil || resp.JSON200.Result.Username == nil {
		return "", fmt.Errorf("unexpected response from GetMe: %d", resp.StatusCode())
	}

	return *resp.JSON200.Result.Username, nil
}
