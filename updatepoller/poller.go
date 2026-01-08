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

	const updatesBufSize = 100

	return &Poller{
		opts:    opts,
		log:     opts.logger,
		updates: make(chan client.Update, updatesBufSize),
	}, nil
}

// Start starts the Poller. The context is used only for the startup timeout.
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

// Stop stops the Poller. The context is used only for the shutdown timeout.
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
			p.log.Errorf("load offset: %v", err)
		}

		return
	}

	resp, err := p.opts.client.GetUpdatesWithResponse(ctx, client.GetUpdatesJSONRequestBody{
		Offset: &offset,
	})
	if err != nil {
		if ctx.Err() == nil {
			p.log.Errorf("fetch updates: %v", err)
		}

		return
	}

	if resp.StatusCode() != http.StatusOK {
		p.log.Errorf("fetch updates: %s", resp.Status())

		return
	}

	if resp.JSON200 == nil || len(resp.JSON200.Result) == 0 {
		return
	}

	newOffset := p.processUpdates(ctx, resp.JSON200.Result)

	if saveErr := p.opts.offsetStore.Save(ctx, newOffset); saveErr != nil {
		if ctx.Err() == nil {
			p.log.Errorf("save offset: %v", saveErr)
		}
	}
}

func (p *Poller) processUpdates(ctx context.Context, updates []client.Update) int {
	var lastID int

	for _, update := range updates {
		select {
		case <-ctx.Done():
			return lastID + 1
		case p.updates <- update:
			lastID = update.UpdateId
		}
	}

	return lastID + 1
}
