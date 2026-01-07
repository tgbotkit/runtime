package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime"
	"github.com/tgbotkit/runtime/events"
	"github.com/tgbotkit/runtime/webhook"
)

func main() {
	_ = godotenv.Load()

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN is required")
	}

	wh, _ := webhook.New(webhook.NewOptions())

	bot, err := runtime.New(runtime.NewOptions(
		token,
		runtime.WithUpdateSource(wh),
	))
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}

	bot.Handlers().OnCommand(func(ctx context.Context, event *events.CommandEvent) error {
		if event.Command != "start" {
			return nil
		}

		text := fmt.Sprintf("You sent the command: /%s\n", event.Command)
		if event.Args != "" {
			text += fmt.Sprintf("With arguments: %s", event.Args)
		} else {
			text += "With no arguments."
		}

		_, _ = bot.Client().SendMessageWithResponse(ctx, client.SendMessageJSONRequestBody{
			ChatId: event.Message.Chat.Id,
			Text:   text,
		})
		return nil
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("Webhook server listening on :8080")
		if err := http.ListenAndServe(":8080", wh); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	if err := bot.Run(ctx); err != nil {
		log.Fatalf("bot error: %v", err)
	}
}
