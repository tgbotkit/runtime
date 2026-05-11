package runtime

import (
	"time"

	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/logger"
)

//go:generate go tool options-gen -out-filename=options.gen.go -from-struct=Options

// Options is the options for the Bot.
type Options struct {
	// botToken is the Telegram bot token.
	botToken string `option:"mandatory"`
	// botUsername is the bot username used for command mention filtering.
	botUsername string

	// client is the Telegram API client.
	client client.ClientWithResponsesInterface
	// eventEmitter is the event emitter to use.
	eventEmitter eventemitter.EventEmitter
	// updateSource is the update source to use.
	updateSource UpdateSource
	// logger is the logger to use.
	logger logger.Logger
	// startupTimeout bounds blocking startup API calls.
	startupTimeout time.Duration `default:"10s" validate:"gt=0"`
	// defaultMiddlewareEnabled controls registration of runtime middleware.
	defaultMiddlewareEnabled bool `default:"true" option:"optional"`
	// defaultListenersEnabled controls registration of runtime listeners.
	defaultListenersEnabled bool `default:"true" option:"optional"`
}
