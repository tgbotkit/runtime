# Getting Started

This guide will help you build your first bot using `tgbotkit-runtime`.

## Prerequisites

- **Go:** You should have Go installed on your machine.
- **Telegram Token:** You'll need a token from [BotFather](https://t.me/BotFather) for your Telegram bot.

## 1. Create a New Project

```bash
mkdir my-bot
cd my-bot
go mod init my-bot
```

## 2. Add `tgbotkit-runtime`

```bash
go get github.com/tgbotkit/runtime
```

## 3. Create Your Bot

Create a `main.go` file with the following content:

```go
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
		log.Fatalf("failed to create bot: %v", err)
	}

	// Register a handler for text messages
	bot.Handlers().OnMessageType(messagetype.Text, func(ctx context.Context, event *events.MessageEvent) error {
		if event.Message.Text != nil {
			log.Printf("Received message: %s", *event.Message.Text)
			
			// Echo the message back
			_, err := bot.Client().SendMessageWithResponse(ctx, client.SendMessageJSONRequestBody{
				ChatId: event.Message.Chat.Id,
				Text:   "You said: " + *event.Message.Text,
			})
			return err
		}
		return nil
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("Bot is running...")
	if err := bot.Run(ctx); err != nil {
		log.Fatalf("bot error: %v", err)
	}
}
```

## 4. Run Your Bot

```bash
export TELEGRAM_TOKEN=your_bot_token_here
go run main.go
```

## Next Steps

- Explore more [Examples](../examples) in the repository.
- Learn about [Events & Handlers](events-and-handlers.md).
- Configure [Update Sources](update-sources.md) (e.g., using Webhooks).
- Add [Middleware](middleware.md) for logging or recovery.
