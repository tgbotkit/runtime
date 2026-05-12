// Package runtime provides the main bot logic and orchestration for the Telegram bot.
package runtime

import (
	"context"
	"fmt"
	"time"

	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime/botcontext"
	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/events"
	"github.com/tgbotkit/runtime/handlers"
	"github.com/tgbotkit/runtime/listeners"
	"github.com/tgbotkit/runtime/logger"
	"github.com/tgbotkit/runtime/middleware"
	"github.com/tgbotkit/runtime/respond"
	"github.com/tgbotkit/runtime/updatepoller"
	"github.com/tgbotkit/runtime/updatepoller/offsetstore"
	"golang.org/x/sync/errgroup"
)

const (
	defaultPollingInterval = time.Second
	defaultPollTimeout     = 30 * time.Second
	defaultRequestTimeout  = defaultPollTimeout + 5*time.Second
)

// Bot is the main bot structure.
type Bot struct {
	opts      Options
	registry  *handlers.Registry
	responder *respond.Responder
}

var _ botcontext.BotContext = (*Bot)(nil)

// New creates a new Bot instance with the given options.
func New(opts Options) (*Bot, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	if opts.client == nil && opts.botToken == "" {
		return nil, fmt.Errorf("bot token is required when client is not provided")
	}

	var err error
	if opts.eventEmitter == nil {
		opts.eventEmitter, err = eventemitter.NewSync(eventemitter.NewOptions(
			eventemitter.WithStopOnError(false),
		))
		if err != nil {
			return nil, fmt.Errorf("create default event emitter: %w", err)
		}
	}

	if opts.logger == nil {
		opts.logger = logger.NewNop()
	}

	if opts.client == nil {
		opts.client, err = newDefaultClient(opts.botToken)
		if err != nil {
			return nil, err
		}
	}

	botName, err := resolveBotName(opts)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		opts:      opts,
		registry:  handlers.NewRegistry(opts.eventEmitter, opts.logger),
		responder: respond.New(opts.client),
	}

	registerDefaults(opts, bot, botName)

	return bot, nil
}

// Run starts the bot's update processing loop and blocks until the context is canceled.
// It initializes a default update poller if no update source is configured.
func (b *Bot) Run(ctx context.Context) error {
	ctx = botcontext.WithBotContext(ctx, b)

	if b.opts.updateSource == nil {
		if err := b.initDefaultPoller(); err != nil {
			return err
		}
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return ToRunnable(b.opts.updateSource).Run(ctx)
	})

	g.Go(func() error {
		return b.receiveLoop(ctx)
	})

	return g.Wait()
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
func (b *Bot) Handlers() *handlers.Registry {
	return b.registry
}

// Responder returns helper methods for common Telegram responses.
func (b *Bot) Responder() *respond.Responder {
	return b.responder
}

func (b *Bot) initDefaultPoller() error {
	poller, err := updatepoller.NewPoller(updatepoller.NewOptions(
		b.opts.client,
		updatepoller.WithOffsetStore(offsetstore.NewInMemoryOffsetStore(0)),
		updatepoller.WithPollingInterval(defaultPollingInterval),
		updatepoller.WithTimeout(defaultPollTimeout),
		updatepoller.WithRequestTimeout(defaultRequestTimeout),
		updatepoller.WithLogger(b.opts.logger),
	))
	if err != nil {
		return fmt.Errorf("create default poller: %w", err)
	}

	b.opts.updateSource = poller

	return nil
}

func newDefaultClient(botToken string) (client.ClientWithResponsesInterface, error) {
	serverURL, err := client.NewServerUrlTelegramBotAPIEndpointSubstituteBotTokenWithYourBotToken(
		client.ServerUrlTelegramBotAPIEndpointSubstituteBotTokenWithYourBotTokenBotTokenVariable(botToken),
	)
	if err != nil {
		return nil, fmt.Errorf("create server URL: %w", err)
	}

	api, err := client.NewClientWithResponses(serverURL)
	if err != nil {
		return nil, fmt.Errorf("create API client: %w", err)
	}

	return api, nil
}

func resolveBotName(opts Options) (string, error) {
	if !opts.defaultListenersEnabled || opts.botUsername != "" {
		return opts.botUsername, nil
	}

	startupCtx, cancelStartup := context.WithTimeout(context.Background(), opts.startupTimeout)
	defer cancelStartup()

	botName, err := loadBotName(startupCtx, opts.client)
	if err != nil {
		return "", fmt.Errorf("load bot name: %w", err)
	}

	return botName, nil
}

func registerDefaults(opts Options, bot *Bot, botName string) {
	if opts.defaultMiddlewareEnabled {
		opts.eventEmitter.Use("*", middleware.ContextInjector(bot))
		opts.eventEmitter.Use("*", middleware.Logger(bot.opts.logger))
		opts.eventEmitter.Use("*", middleware.Recoverer(bot.opts.logger))
	}

	if opts.defaultListenersEnabled {
		opts.eventEmitter.AddListener(events.OnUpdate, listeners.Classifier(opts.eventEmitter))
		opts.eventEmitter.AddListener(events.OnMessage, listeners.CommandParser(opts.eventEmitter, botName))
	}
}

func loadBotName(ctx context.Context, api client.ClientWithResponsesInterface) (string, error) {
	// Fetch bot's own info to get the username
	// It is important to do it once at startup
	me, err := api.GetMeWithResponse(ctx)
	if err != nil {
		return "", fmt.Errorf("get bot info: %w", err)
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
			return nil
		case update, ok := <-ch:
			if !ok {
				if ctx.Err() != nil {
					return nil
				}

				return ErrUpdateSourceClosed
			}

			b.Logger().Debugf("got update: %v", update.UpdateId)
			b.opts.eventEmitter.Emit(ctx, events.OnUpdate, &events.UpdateEvent{Update: &update})
		}
	}
}
