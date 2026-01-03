package webhook

import (
	"encoding/json"
	"net/http"

	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime"
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
func (h *Webhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
