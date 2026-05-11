// Package webhook provides an implementation of UpdateSource using Telegram webhooks.
package webhook

import "github.com/tgbotkit/client"

//go:generate go tool options-gen -out-filename=options.gen.go -from-struct=Options

// Options is the options for the Webhook handler.
type Options struct {
	token              string                              `option:"optional"`
	url                string                              `option:"optional"`
	client             client.ClientWithResponsesInterface `option:"optional"`
	bufferSize         int                                 `default:"100"     option:"optional"                  validate:"min=1"` //nolint:lll
	maxBodyBytes       int64                               `default:"1048576" option:"optional"                  validate:"min=1"` //nolint:lll
	allowedUpdates     []string                            `option:"optional"`
	dropPendingUpdates bool                                `option:"optional"`
	maxConnections     int                                 `option:"optional" validate:"omitempty,min=1,max=100"`

	webhookRegistrationEnabled    bool `option:"-"`
	webhookRegistrationConfigured bool `option:"-"`
}

// WithWebhookRegistrationEnabled controls whether Start configures the Telegram webhook.
func WithWebhookRegistrationEnabled(enabled bool) OptOptionsSetter {
	return func(o *Options) {
		o.webhookRegistrationEnabled = enabled
		o.webhookRegistrationConfigured = true
	}
}
