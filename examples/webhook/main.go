// Package main provides an example of a bot using webhooks.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tgbotkit/runtime"
	"github.com/tgbotkit/runtime/events"
	"github.com/tgbotkit/runtime/handlers"
	"github.com/tgbotkit/runtime/webhook"
)

const shutdownTimeout = 5 * time.Second

func main() {
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN is required")
	}

	// Initialize webhook update source. Registration is disabled here because
	// this example assumes the public Telegram webhook is managed externally.
	wh, err := webhook.New(webhook.NewOptions(
		webhook.WithWebhookRegistrationEnabled(false),
	))
	if err != nil {
		log.Fatalf("create webhook: %v", err)
	}

	bot, err := runtime.New(runtime.NewOptions(
		token,
		runtime.WithUpdateSource(wh),
	))
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

	server := &http.Server{
		Addr:    ":8080",
		Handler: wh,
	}

	go func() {
		log.Printf("Webhook server listening on :8080")

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	runErr := bot.Run(ctx)

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancelShutdown()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	if runErr != nil {
		log.Fatalf("bot error: %v", runErr)
	}
}
