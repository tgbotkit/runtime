package updatepoller

import (
	"time"

	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime/logger"
)

//go:generate go tool options-gen -out-filename=options.gen.go -from-struct=Options

// Options is the options for the Poller.
type Options struct {
	// client is the Telegram API client.
	client client.ClientWithResponsesInterface `option:"mandatory" validate:"required"`
	// pollingInterval is the interval between polling requests.
	pollingInterval time.Duration `default:"1s" validate:"gt=0"`
	// timeout is the long-polling timeout sent to Telegram.
	timeout time.Duration `default:"30s" validate:"gte=0"`
	// requestTimeout is the per-GetUpdates transport deadline. Zero disables it.
	requestTimeout time.Duration `validate:"gte=0"`
	// limit is the maximum number of updates to fetch per request.
	limit int `validate:"omitempty,min=1,max=100"`
	// allowedUpdates restricts the Telegram update types returned by polling.
	allowedUpdates []string
	// bufferSize is the capacity of the update channel.
	bufferSize int `default:"100" validate:"gt=0"`
	// offsetStore is the store to use for storing the update offset.
	offsetStore OffsetStore `validate:"required"`
	// logger is the logger to use.
	logger logger.Logger
}
