package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime"
	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/events"
	"github.com/tgbotkit/runtime/logger"
	"github.com/tgbotkit/runtime/webhook"
)

// PingPongHandler is a simple handler that responds "pong" to "ping" messages.
type PingPongHandler struct{}

// Subscribe registers the handler to listen for text messages and commands.
func (h *PingPongHandler) Subscribe(ee eventemitter.EventEmitter) {
	eventemitter.On[events.MessageEvent](ee, events.OnMessageReceived, h.onTextMessage)
}

func (h *PingPongHandler) onTextMessage(ctx context.Context, event *events.MessageEvent) error {
	if event.Message.Text == nil {
		return nil
	}

	botCtx := event.Bot

	if *event.Message.Text == "ping" {
		botCtx.Logger().Infof("Received ping from chat %d", event.Message.Chat.Id)
		resp, err := botCtx.Client().SendMessageWithResponse(ctx, client.SendMessageJSONRequestBody{
			ChatId: event.Message.Chat.Id,
			Text:   "pong",
		})
		if err != nil {
			botCtx.Logger().Errorf("Failed to send pong: %v", err)
			return err
		}

		if resp.StatusCode() != http.StatusOK {
			botCtx.Logger().Errorf("Failed to send pong: %s, body: %s", resp.Status(), string(resp.Body))
			return errors.New("failed to send pong")
		}
	}
	return nil
}

func main() {
	_ = godotenv.Load()

	// 1. Initialize core components
	log := logger.NewSlog(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	// 2. Load configuration
	token := os.Getenv("TELEGRAM_TOKEN")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 3. Validate configuration
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN is required")
	}

	// 4. Initialize the webhook handler
	webhookHandler, err := webhook.New(webhook.NewOptions())
	if err != nil {
		log.Fatalf("failed to create webhook handler: %v", err)
	}

	// 5. Initialize handlers
	pingPongHandler := &PingPongHandler{}

	// 6. Initialize the bot
	// We pass the update stream as the UpdateSource.
	botInstance, err := runtime.New(runtime.NewOptions(
		token,
		runtime.WithLogger(log),
		runtime.WithUpdateSource(webhookHandler),
	))
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}

	// 7. Add handlers
	botInstance.AddHandler(pingPongHandler)

	// 8. Set up a graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// 9. Start the webhook server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: webhookHandler,
	}

	go func() {
		log.Infof("Webhook server is listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("webhook server error: %v", err)
		}
	}()

	// 10. Run the bot
	log.Info("Bot is running... Press Ctrl+C to stop.")
	if err := botInstance.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("bot runtime error: %v", err)
	}

	// 11. Shutdown the server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Errorf("failed to shutdown server: %v", err)
	}
}
