package runtime

import (
	"context"
	"fmt"
	"time"

	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/events"
	"github.com/tgbotkit/runtime/internal/classifier"
	"github.com/tgbotkit/runtime/internal/commandparser"
	"github.com/tgbotkit/runtime/logger"
	"github.com/tgbotkit/runtime/updatepoller"
	"github.com/tgbotkit/runtime/updatepoller/offsetstore"
)

// Bot is the main bot structure.
type Bot struct {
	opts Options
}

// New creates a new Bot instance with the given options.
func New(opts Options) (*Bot, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	if opts.eventEmitter == nil {
		var err error
		opts.eventEmitter, err = eventemitter.NewSync(eventemitter.NewOptions())
		if err != nil {
			return nil, fmt.Errorf("failed to create default event emitter: %w", err)
		}
	}

	if opts.logger == nil {
		opts.logger = logger.NewNop()
	}

	if opts.client == nil {
		serverURL, err := client.NewServerUrlTelegramBotAPIEndpointSubstituteBotTokenWithYourBotToken(
			client.ServerUrlTelegramBotAPIEndpointSubstituteBotTokenWithYourBotTokenBotTokenVariable(opts.botToken),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create server URL: %w", err)
		}

		opts.client, err = client.NewClientWithResponses(serverURL)
		if err != nil {
			return nil, fmt.Errorf("failed to create API client: %w", err)
		}
	}

	bot := &Bot{
		opts: opts,
	}

	// Subscribe internal components
	classifier.Subscribe(opts.eventEmitter)

	cp, err := commandparser.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create command parser: %w", err)
	}
	cp.Subscribe(opts.eventEmitter)

	return bot, nil
}

// Run starts the bot's update processing loop and blocks until the context is canceled.
// It initializes a default update poller if no update source is configured.
func (b *Bot) Run(ctx context.Context) error {
	b.opts.eventEmitter.Emit(ctx, events.OnBeforeStart, &events.BotEvent{Bot: b})

	// init default update source if not provided
	if b.opts.updateSource == nil {
		poller, err := updatepoller.NewPoller(updatepoller.NewOptions(
			b.opts.client,
			updatepoller.WithOffsetStore(offsetstore.NewInMemoryOffsetStore(0)),
			updatepoller.WithPollingInterval(time.Second),
			updatepoller.WithLogger(b.opts.logger),
		))
		if err != nil {
			return fmt.Errorf("failed to create default poller: %w", err)
		}
		b.opts.updateSource = poller
	}

	return b.receiveLoop(ctx)
}

// AddHandler registers a new event handler with the bot's event emitter.
func (b *Bot) AddHandler(h Handler) *Bot {
	h.Subscribe(b.opts.eventEmitter)
	return b
}

// Client returns the underlying Telegram Bot API client.
func (b *Bot) Client() client.ClientWithResponsesInterface {
	return b.opts.client
}

// EventEmitter returns the bot's event emitter.
func (b *Bot) EventEmitter() eventemitter.EventEmitter {
	return b.opts.eventEmitter
}

// Logger returns the bot's logger.
func (b *Bot) Logger() logger.Logger {
	return b.opts.logger
}

// receiveLoop is the main event loop that consumes updates from the update source
// and emits them to the event emitter.
func (b *Bot) receiveLoop(ctx context.Context) error {
	ch := b.opts.updateSource.UpdateChan()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case update, ok := <-ch:
			if !ok {
				// Channel closed, graceful shutdown.
				return nil
			}
			b.opts.eventEmitter.Emit(ctx, events.OnUpdateReceived, &events.UpdateEvent{Bot: b, Update: &update})
		}
	}
}

// Handler is an interface for components that subscribe to events.
type Handler interface {
	Subscribe(eventemitter.EventEmitter)
}

// HandlerFunc is an adapter to allow the use of ordinary functions as Handlers.
type HandlerFunc func(eventemitter.EventEmitter)

// Subscribe calls f(ee), allowing HandlerFunc to implement the Handler interface.
func (f HandlerFunc) Subscribe(ee eventemitter.EventEmitter) {
f(ee)
}
