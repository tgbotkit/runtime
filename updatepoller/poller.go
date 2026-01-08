// Package updatepoller provides a long-polling mechanism for receiving updates from the Telegram Bot API.
package updatepoller

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime/logger"
)

// Poller polls the Telegram API for updates.
type Poller struct {
	opts   Options
	log    logger.Logger
	cancel context.CancelFunc
	wg     sync.WaitGroup

	updates chan client.Update
}

// NewPoller creates a new Poller instance with the given options.
func NewPoller(opts Options) (*Poller, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("invalid poller options: %w", err)
	}
	if opts.logger == nil {
		opts.logger = logger.NewNop()
	}
	return &Poller{
		opts:    opts,
		log:     opts.logger,
		updates: make(chan client.Update, 100),
	}, nil
}

// Start starts the Poller.
func (p *Poller) Start(_ context.Context) error {
	if p.cancel != nil {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	p.wg.Add(1)

	go func() {
		defer p.wg.Done()
		ticker := time.NewTicker(p.opts.pollingInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				p.poll(ctx)
			}
		}
	}()

	return nil
}

// Stop stops the Poller.
func (p *Poller) Stop(ctx context.Context) error {
	if p.cancel != nil {
		p.cancel()
		p.cancel = nil
	}

	c := make(chan struct{})
	go func() {
		defer close(c)
		p.wg.Wait()
	}()

	select {
	case <-c:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// UpdateChan returns the updates channel.
func (p *Poller) UpdateChan() <-chan client.Update {
	return p.updates
}

func (p *Poller) poll(ctx context.Context) {
	offset, err := p.opts.offsetStore.Load(ctx)
	if err != nil {
		if ctx.Err() == nil {
			p.log.Errorf("failed to load offset: %v", err)
		}
		return
	}

	resp, err := p.opts.client.GetUpdatesWithResponse(ctx, client.GetUpdatesJSONRequestBody{
		Offset: &offset,
	})
	if err != nil || resp.StatusCode() != http.StatusOK {
		if err != nil {
			if ctx.Err() == nil {
				p.log.Errorf("failed to fetch updates: %v", err)
			}
		} else if resp.StatusCode() != http.StatusOK {
			p.log.Errorf("failed to fetch updates: %s", resp.Status())
		}
		return
	}

	if resp.JSON200 == nil {
		return
	}

	for _, update := range resp.JSON200.Result {
		select {
		case <-ctx.Done():
			return
		case p.updates <- update:
			offset = update.UpdateId + 1
		}
	}

	if len(resp.JSON200.Result) > 0 {
		if saveErr := p.opts.offsetStore.Save(ctx, offset); saveErr != nil {
			if ctx.Err() == nil {
				p.log.Errorf("failed to save offset: %v", saveErr)
			}
		}
	}
}