// Package main provides a simple ping-pong bot example.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tgbotkit/runtime"
	"github.com/tgbotkit/runtime/events"
	"github.com/tgbotkit/runtime/handlers"
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

	// Register a handler for ping text messages only.
	bot.Handlers().OnMessageMatch(
		handlers.MessageText("ping"),
		func(ctx context.Context, event *events.MessageEvent) error {
			_, err := bot.Responder().SendTextInChat(ctx, event.Message, "pong")

			return err
		},
	)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := bot.Run(ctx); err != nil {
		log.Fatalf("bot error: %v", err)
	}
}
