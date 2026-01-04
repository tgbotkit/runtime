package webhook

import (
	"encoding/json"
	"net/http"

	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime"
)

const (
	// HeaderTelegramBotAPISecretToken is the header used by Telegram to send the secret token.
	HeaderTelegramBotAPISecretToken = "X-Telegram-Bot-Api-Secret-Token"
)

var _ runtime.UpdateSource = (*Webhook)(nil)

// Webhook implements http.Handler to receive incoming updates via an outgoing webhook.
type Webhook struct {
	opts    Options
	updates chan client.Update
}

// New creates a new Webhook handler.
func New(opts Options) (*Webhook, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	return &Webhook{
		opts:    opts,
		updates: make(chan client.Update, 100),
	}, nil
}

// UpdateChan returns the updates channel.
func (h *Webhook) UpdateChan() <-chan client.Update {
	return h.updates
}

// ServeHTTP implements http.Handler interface.
// It validates the request, decodes the update, and sends it to the updates channel.
func (h *Webhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// If a secret token is configured, validate the header.
	if len(h.opts.Token) > 0 {
		if r.Header.Get(HeaderTelegramBotAPISecretToken) != h.opts.Token {
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

	select {
	case h.updates <- update:
		w.WriteHeader(http.StatusOK)
	case <-r.Context().Done():
		w.WriteHeader(http.StatusRequestTimeout)
	}
}
