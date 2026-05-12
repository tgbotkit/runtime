package handlers

import (
	"strings"

	"github.com/tgbotkit/runtime/events"
	"github.com/tgbotkit/runtime/messagetype"
)

// CommandName matches a command by exact name.
func CommandName(name string) CommandMatcher {
	return func(event *events.CommandEvent) bool {
		return event != nil && event.Command == name
	}
}

// CommandAny matches any of the provided command names.
func CommandAny(names ...string) CommandMatcher {
	return func(event *events.CommandEvent) bool {
		if event == nil {
			return false
		}

		for _, name := range names {
			if event.Command == name {
				return true
			}
		}

		return false
	}
}

// CallbackData matches callback query data by exact value.
func CallbackData(data string) CallbackQueryMatcher {
	return func(event *events.CallbackQueryEvent) bool {
		if event == nil || event.CallbackQuery == nil || event.CallbackQuery.Data == nil {
			return false
		}

		return *event.CallbackQuery.Data == data
	}
}

// CallbackDataPrefix matches callback query data by prefix.
func CallbackDataPrefix(prefix string) CallbackQueryMatcher {
	return func(event *events.CallbackQueryEvent) bool {
		if event == nil || event.CallbackQuery == nil || event.CallbackQuery.Data == nil {
			return false
		}

		return strings.HasPrefix(*event.CallbackQuery.Data, prefix)
	}
}

// MessageText matches text messages by exact text.
func MessageText(text string) MessageMatcher {
	return func(event *events.MessageEvent) bool {
		if event == nil || event.Message == nil || event.Message.Text == nil {
			return false
		}

		return *event.Message.Text == text
	}
}

// MessageTextPrefix matches text messages by prefix.
func MessageTextPrefix(prefix string) MessageMatcher {
	return func(event *events.MessageEvent) bool {
		if event == nil || event.Message == nil || event.Message.Text == nil {
			return false
		}

		return strings.HasPrefix(*event.Message.Text, prefix)
	}
}

// MessageType matches messages by classified message type.
func MessageType(t messagetype.MessageType) MessageMatcher {
	return func(event *events.MessageEvent) bool {
		return event != nil && event.Type == t
	}
}
