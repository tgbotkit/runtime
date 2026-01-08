// Package main provides an example of a bot that handles commands.
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

	// Register a handler for the /start command
	bot.Handlers().OnCommand(func(ctx context.Context, event *events.CommandEvent) error {
		if event.Command == "start" {
			log.Printf("Received /start command from %d", event.Message.Chat.Id)
			_, _ = bot.Client().SendMessageWithResponse(ctx, client.SendMessageJSONRequestBody{
				ChatId: event.Message.Chat.Id,
				Text:   "Hello! I am a bot built with tgbotkit.",
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