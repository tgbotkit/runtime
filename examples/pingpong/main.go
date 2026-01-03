package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime"
	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/events"
	"github.com/tgbotkit/runtime/logger"
)

// PingPongHandler is a simple handler that responds "pong" to "ping" messages.
type PingPongHandler struct {
}

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
	// 3. Validate configuration
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN is required")
	}

	// 4. Initialize handlers
	pingPongHandler := &PingPongHandler{}

	// 5. Initialize the bot
	botInstance, err := runtime.New(runtime.NewOptions(
		token,
		runtime.WithLogger(log),
	))
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}

	botInstance.AddHandler(pingPongHandler)

	// 7. Set up a graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// 8. Log startup message
	log.Info("Ping-pong bot is running... Press Ctrl+C to stop.")

	// 9. Run the bot
	if err := botInstance.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("bot runtime error: %v", err)
	}
}
