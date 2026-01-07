package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime"
	"github.com/tgbotkit/runtime/events"
	"github.com/tgbotkit/runtime/messagetype"
)

func main() {
	_ = godotenv.Load()

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN is required")
	}

	bot, err := runtime.New(runtime.NewOptions(token))
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}

	bot.Handlers().OnMessageWithType(func(ctx context.Context, event *events.MessageEvent) error {
		if event.Message.Text != nil && *event.Message.Text == "ping" {
			_, _ = bot.Client().SendMessageWithResponse(ctx, client.SendMessageJSONRequestBody{
				ChatId: event.Message.Chat.Id,
				Text:   "pong",
			})
		}
		return nil
	}, messagetype.Text)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("Bot is running... Press Ctrl+C to stop.")
	if err := bot.Run(ctx); err != nil {
		log.Fatalf("bot error: %v", err)
	}
}
