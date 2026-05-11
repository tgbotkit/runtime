package webhook

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	if !h.registrationEnabled() {
		return nil
	}

	if err := h.validateWebhookRegistration(); err != nil {
		return err
	}

	resp, err := h.opts.client.SetWebhookWithResponse(ctx, h.setWebhookRequest())
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
	if !h.authorize(w, r) {
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	update, ok := h.decodeUpdate(w, r)
	if !ok {
		return
	}

	h.enqueueUpdate(w, r, update)
}

func (h *Webhook) registrationEnabled() bool {
	if h.opts.webhookRegistrationConfigured {
		return h.opts.webhookRegistrationEnabled
	}

	return h.opts.client != nil || h.opts.url != ""
}

func (h *Webhook) validateWebhookRegistration() error {
	if h.opts.client == nil {
		return fmt.Errorf("set webhook: client is required")
	}

	if h.opts.url == "" {
		return fmt.Errorf("set webhook: url is required")
	}

	return nil
}

func (h *Webhook) setWebhookRequest() client.SetWebhookJSONRequestBody {
	params := client.SetWebhookJSONRequestBody{
		Url: h.opts.url,
	}

	if h.opts.token != "" {
		params.SecretToken = &h.opts.token
	}

	if h.opts.allowedUpdates != nil {
		params.AllowedUpdates = &h.opts.allowedUpdates
	}

	if h.opts.dropPendingUpdates {
		params.DropPendingUpdates = &h.opts.dropPendingUpdates
	}

	if h.opts.maxConnections > 0 {
		params.MaxConnections = &h.opts.maxConnections
	}

	return params
}

func (h *Webhook) authorize(w http.ResponseWriter, r *http.Request) bool {
	if len(h.opts.token) == 0 || r.Header.Get(HeaderTelegramBotAPISecretToken) == h.opts.token {
		return true
	}

	w.WriteHeader(http.StatusUnauthorized)

	return false
}

func (h *Webhook) decodeUpdate(w http.ResponseWriter, r *http.Request) (client.Update, bool) {
	var update client.Update

	r.Body = http.MaxBytesReader(w, r.Body, h.opts.maxBodyBytes)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&update); err != nil {
		writeDecodeError(w, err)

		return client.Update{}, false
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		w.WriteHeader(http.StatusBadRequest)

		return client.Update{}, false
	}

	return update, true
}

func writeDecodeError(w http.ResponseWriter, err error) {
	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		w.WriteHeader(http.StatusRequestEntityTooLarge)

		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

func (h *Webhook) enqueueUpdate(w http.ResponseWriter, r *http.Request, update client.Update) {
	if bc := botcontext.FromContext(r.Context()); bc != nil {
		bc.Logger().Debugf("got update: %v", update.UpdateId)
	}

	select {
	case h.updates <- update:
		w.WriteHeader(http.StatusOK)
	case <-r.Context().Done():
		w.WriteHeader(http.StatusRequestTimeout)
	default:
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}
