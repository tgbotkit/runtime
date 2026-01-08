// Package main provides a simple ping-pong bot example.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime"
	"github.com/tgbotkit/runtime/events"
	"github.com/tgbotkit/runtime/messagetype"
)

func main() {
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN is required")
	}

	bot, err := runtime.New(runtime.NewOptions(token))
	if err != nil {
		log.Fatalf("create bot: %v", err)
	}

	// Register a handler for text message events only
	bot.Handlers().OnMessageType(messagetype.Text, func(ctx context.Context, event *events.MessageEvent) error {
		if event.Message.Text != nil && *event.Message.Text == "ping" {
			_, _ = bot.Client().SendMessageWithResponse(ctx, client.SendMessageJSONRequestBody{
				ChatId: event.Message.Chat.Id,
				Text:   "pong",
			})
		}
		return nil
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := bot.Run(ctx); err != nil {
		log.Fatalf("bot error: %v", err)
	}
}