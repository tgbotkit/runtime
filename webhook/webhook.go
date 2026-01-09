package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/metalagman/appkit/lifecycle"
	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime/botcontext"
)

const (
	// HeaderTelegramBotAPISecretToken is the header used by Telegram to send the secret token.
	HeaderTelegramBotAPISecretToken = "X-Telegram-Bot-Api-Secret-Token"
)

// Webhook implements http.Handler to receive incoming updates via an outgoing webhook.
type Webhook struct {
	opts    Options
	updates chan client.Update
}

var _ http.Handler = (*Webhook)(nil)
var _ lifecycle.Lifecycle = (*Webhook)(nil)

// New creates a new Webhook handler.
func New(opts Options) (*Webhook, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	return &Webhook{
		opts:    opts,
		updates: make(chan client.Update, opts.bufferSize),
	}, nil
}

// UpdateChan returns the updates channel.
func (h *Webhook) UpdateChan() <-chan client.Update {
	return h.updates
}

// Start satisfies the lifecycle.Lifecycle interface. The context is used only for the startup timeout.
func (h *Webhook) Start(ctx context.Context) error {
	return h.SetWebhook(ctx)
}

// Stop satisfies the lifecycle.Lifecycle interface. The context is used only for the shutdown timeout.
func (h *Webhook) Stop(_ context.Context) error {
	return nil
}

// SetWebhook sets the outgoing webhook for the bot.
func (h *Webhook) SetWebhook(ctx context.Context) error {
	if h.opts.client == nil || h.opts.url == "" {
		return nil
	}

	params := client.SetWebhookJSONRequestBody{
		Url: h.opts.url,
	}
	if h.opts.token != "" {
		params.SecretToken = &h.opts.token
	}

	resp, err := h.opts.client.SetWebhookWithResponse(ctx, params)
	if err != nil {
		return fmt.Errorf("set webhook: %w", err)
	}

	if resp.JSON200 == nil || !bool(resp.JSON200.Ok) {
		return fmt.Errorf("set webhook: unexpected response: %s", resp.Status())
	}

	return nil
}

// ServeHTTP implements http.Handler interface.
// It validates the request, decodes the update, and sends it to the updates channel.
func (h *Webhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// If a secret token is configured, validate the header.
	if len(h.opts.token) > 0 {
		if r.Header.Get(HeaderTelegramBotAPISecretToken) != h.opts.token {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	var update client.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if bc := botcontext.FromContext(r.Context()); bc != nil {
		bc.Logger().Debugf("got update: %v", update.UpdateId)
	}

	select {
	case h.updates <- update:
		w.WriteHeader(http.StatusOK)
	case <-r.Context().Done():
		w.WriteHeader(http.StatusRequestTimeout)
	}
}