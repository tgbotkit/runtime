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
	pollingInterval time.Duration `default:"1s"`
	// offsetStore is the store to use for storing the update offset.
	offsetStore OffsetStore `validate:"required"`
	// logger is the logger to use.
	logger logger.Logger
}
