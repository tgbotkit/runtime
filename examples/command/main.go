// Package main provides an example of a bot that handles commands.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	// Register a handler for the /start command.
	bot.Handlers().OnCommandName("start", func(ctx context.Context, event *events.CommandEvent) error {
		_, err := bot.Responder().SendTextInChat(ctx, event.Message, "Hello! I am a bot built with tgbotkit-runtime.")

		return err
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := bot.Run(ctx); err != nil {
		log.Fatalf("bot error: %v", err)
	}
}
