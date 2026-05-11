// Package updatepoller provides a long-polling mechanism for receiving updates from the Telegram Bot API.
package updatepoller

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/metalagman/appkit/lifecycle"
	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime/logger"
)

// Poller polls the Telegram API for updates.
type Poller struct {
	opts Options
	log  logger.Logger

	mu     sync.Mutex
	cancel context.CancelFunc
	done   chan struct{}

	updates chan client.Update
}

const (
	minRetryBackoff = 100 * time.Millisecond
	maxRetryBackoff = 5 * time.Second
)

var _ lifecycle.Lifecycle = (*Poller)(nil)

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
		updates: make(chan client.Update, opts.bufferSize),
	}, nil
}

// Start starts the Poller. The context is used only for the startup timeout.
func (p *Poller) Start(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	p.mu.Lock()

	if p.cancel != nil {
		p.mu.Unlock()

		return nil
	}

	runCtx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	p.done = make(chan struct{})
	done := p.done
	p.mu.Unlock()

	go func() {
		defer p.finish(done)

		backoff := pollRetryBackoff{}
		delay := p.opts.pollingInterval

		for {
			if !sleepContext(runCtx, delay) {
				return
			}

			if p.poll(runCtx) {
				backoff.Reset()

				delay = p.opts.pollingInterval

				continue
			}

			delay = backoff.Next()
		}
	}()

	return nil
}

// Stop stops the Poller. The context is used only for the shutdown timeout.
func (p *Poller) Stop(ctx context.Context) error {
	p.mu.Lock()
	cancel := p.cancel
	done := p.done

	if cancel == nil {
		p.mu.Unlock()

		return nil
	}

	cancel()
	p.mu.Unlock()

	select {
	case <-done:
		p.mu.Lock()

		if p.done == done {
			p.cancel = nil
			p.done = nil
		}

		p.mu.Unlock()

		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// UpdateChan returns the updates channel.
func (p *Poller) UpdateChan() <-chan client.Update {
	return p.updates
}

func (p *Poller) finish(done chan struct{}) {
	p.mu.Lock()

	if p.done == done {
		p.cancel = nil
		p.done = nil
	}

	p.mu.Unlock()

	close(done)
}

func (p *Poller) poll(ctx context.Context) bool {
	offset, err := p.opts.offsetStore.Load(ctx)
	if err != nil {
		if ctx.Err() == nil {
			p.log.Errorf("load offset: %v", err)
		}

		return false
	}

	requestCtx, cancel := p.getUpdatesContext(ctx)
	defer cancel()

	resp, err := p.opts.client.GetUpdatesWithResponse(requestCtx, p.getUpdatesRequest(offset))
	if err != nil {
		if ctx.Err() == nil {
			p.log.Errorf("fetch updates: %v", err)
		}

		return false
	}

	if resp.StatusCode() != http.StatusOK {
		p.log.Errorf("fetch updates: %s", resp.Status())

		return false
	}

	if resp.JSON200 == nil || len(resp.JSON200.Result) == 0 {
		return true
	}

	newOffset, ok := p.processUpdates(ctx, resp.JSON200.Result)
	if !ok {
		return true
	}

	p.saveOffset(ctx, newOffset)

	return true
}

func (p *Poller) getUpdatesContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if p.opts.requestTimeout <= 0 {
		return ctx, func() {}
	}

	return context.WithTimeout(ctx, p.opts.requestTimeout)
}

func (p *Poller) getUpdatesRequest(offset int) client.GetUpdatesJSONRequestBody {
	body := client.GetUpdatesJSONRequestBody{
		Offset: &offset,
	}

	if p.opts.timeout > 0 {
		timeout := int(p.opts.timeout / time.Second)
		if timeout == 0 {
			timeout = 1
		}

		body.Timeout = &timeout
	}

	if p.opts.limit > 0 {
		body.Limit = &p.opts.limit
	}

	if p.opts.allowedUpdates != nil {
		body.AllowedUpdates = &p.opts.allowedUpdates
	}

	return body
}

func (p *Poller) saveOffset(ctx context.Context, offset int) {
	if saveErr := p.opts.offsetStore.Save(ctx, offset); saveErr != nil {
		if ctx.Err() == nil {
			p.log.Errorf("save offset: %v", saveErr)
		}
	}
}

func (p *Poller) processUpdates(ctx context.Context, updates []client.Update) (int, bool) {
	var newOffset int

	for _, update := range updates {
		select {
		case <-ctx.Done():
			return newOffset, false
		case p.updates <- update:
			newOffset = update.UpdateId + 1
		}
	}

	return newOffset, true
}

type pollRetryBackoff struct {
	current time.Duration
}

func (b *pollRetryBackoff) Next() time.Duration {
	if b.current == 0 {
		b.current = minRetryBackoff

		return b.current
	}

	b.current *= 2
	if b.current > maxRetryBackoff {
		b.current = maxRetryBackoff
	}

	return b.current
}

func (b *pollRetryBackoff) Reset() {
	b.current = 0
}

func sleepContext(ctx context.Context, delay time.Duration) bool {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}
