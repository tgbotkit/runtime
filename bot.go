package runtime

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime/botcontext"
	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/events"
	"github.com/tgbotkit/runtime/handlers"
	"github.com/tgbotkit/runtime/listeners"
	"github.com/tgbotkit/runtime/logger"
	"github.com/tgbotkit/runtime/middleware"
	"github.com/tgbotkit/runtime/updatepoller"
	"github.com/tgbotkit/runtime/updatepoller/offsetstore"
)

// Bot is the main bot structure.
type Bot struct {
	opts     Options
	registry *handlers.Registry
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

	botName, err := loadBotName(opts.client)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		opts: opts,
	}
	bot.registry = handlers.NewRegistry(opts.eventEmitter, opts.logger)

	// Register internal middlewares
	opts.eventEmitter.Use("*", middleware.ContextInjector(bot))
	opts.eventEmitter.Use("*", middleware.Logger(bot.opts.logger))
	opts.eventEmitter.Use("*", middleware.Recoverer(bot.opts.logger))

	opts.eventEmitter.AddListener(events.OnUpdate, listeners.Classifier(opts.eventEmitter))
	opts.eventEmitter.AddListener(events.OnMessage, listeners.CommandParser(opts.eventEmitter, botName))

	return bot, nil
}

// Run starts the bot's update processing loop and blocks until the context is canceled.
// It initializes a default update poller if no update source is configured.
func (b *Bot) Run(ctx context.Context) error {
	ctx = botcontext.WithBotContext(ctx, b)

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

// Handlers returns the bot's handler registry.
func (b *Bot) Handlers() handlers.HandlerRegistry {
	return b.registry
}

func loadBotName(api client.ClientWithResponsesInterface) (string, error) {
	// Fetch bot's own info to get the username
	// It is important to do it once at startup
	me, err := api.GetMeWithResponse(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get bot info: %w", err)
	}
	if me.JSON200 == nil || me.JSON200.Result.Username == nil {
		return "", fmt.Errorf("could not retrieve bot username from GetMe response")
	}
	return *me.JSON200.Result.Username, nil
}

// receiveLoop is the main event loop that consumes updates from the update source
// and emits them to the event emitter.
func (b *Bot) receiveLoop(ctx context.Context) error {
	ch := b.opts.updateSource.UpdateChan()
	b.Logger().Info("receive loop started")
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case update, ok := <-ch:
			if !ok {
				// Channel closed, graceful shutdown.
				return nil
			}
			zerolog.Ctx(ctx).Debug().Interface("update", update).Msg("got update")
			b.opts.eventEmitter.Emit(ctx, events.OnUpdate, &events.UpdateEvent{Bot: b, Update: &update})
		}
	}
}
