package updatepoller

import (
	"context"
	"errors"
	"net/http"
	"slices"
	"sync/atomic"
	"testing"
	"time"

	"github.com/tgbotkit/client"
)

type pollerMockClient struct {
	client.ClientWithResponsesInterface
	getUpdatesFunc func(ctx context.Context, body client.GetUpdatesJSONRequestBody) (*client.GetUpdatesResponse, error)
}

func (m *pollerMockClient) GetUpdatesWithResponse(
	ctx context.Context,
	body client.GetUpdatesJSONRequestBody,
	_ ...client.RequestEditorFn,
) (*client.GetUpdatesResponse, error) {
	return m.getUpdatesFunc(ctx, body)
}

type pollerOffsetStore struct {
	offset   atomic.Int64
	saves    atomic.Int64
	failSave atomic.Bool
}

var errSaveOffset = errors.New("save offset")

func newPollerOffsetStore(offset int) *pollerOffsetStore {
	store := &pollerOffsetStore{}
	store.offset.Store(int64(offset))

	return store
}

func (s *pollerOffsetStore) Load(_ context.Context) (int, error) {
	return int(s.offset.Load()), nil
}

func (s *pollerOffsetStore) Save(_ context.Context, offset int) error {
	s.saves.Add(1)
	if s.failSave.Load() {
		return errSaveOffset
	}

	s.offset.Store(int64(offset))

	return nil
}

func (s *pollerOffsetStore) offsetValue() int {
	return int(s.offset.Load())
}

func (s *pollerOffsetStore) saveCount() int {
	return int(s.saves.Load())
}

func TestPollerDoesNotSaveOffsetWhenEnqueueInterrupted(t *testing.T) {
	store := newPollerOffsetStore(7)
	tgClient := &pollerMockClient{
		getUpdatesFunc: func(_ context.Context, _ client.GetUpdatesJSONRequestBody) (*client.GetUpdatesResponse, error) {
			return getUpdatesResponse([]client.Update{
				{UpdateId: 10},
				{UpdateId: 11},
			}), nil
		},
	}

	p, err := NewPoller(NewOptions(tgClient, WithOffsetStore(store)))
	if err != nil {
		t.Fatalf("NewPoller() unexpected error: %v", err)
	}
	p.updates = make(chan client.Update, 1)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		p.poll(ctx)
	}()

	waitForUpdateBufferLen(t, p.updates, 1)
	cancel()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("poll did not return after context cancellation")
	}

	if got := store.saveCount(); got != 0 {
		t.Fatalf("save count=%d, want 0", got)
	}
	if got := store.offsetValue(); got != 7 {
		t.Fatalf("offset=%d, want 7", got)
	}
}

func TestPollerSavesLastUpdateIDAfterSuccessfulBatch(t *testing.T) {
	store := newPollerOffsetStore(7)
	updates := []client.Update{
		{UpdateId: 10},
		{UpdateId: 11},
		{UpdateId: 15},
	}
	tgClient := &pollerMockClient{
		getUpdatesFunc: func(_ context.Context, _ client.GetUpdatesJSONRequestBody) (*client.GetUpdatesResponse, error) {
			return getUpdatesResponse(updates), nil
		},
	}

	p, err := NewPoller(NewOptions(tgClient, WithOffsetStore(store)))
	if err != nil {
		t.Fatalf("NewPoller() unexpected error: %v", err)
	}
	p.updates = make(chan client.Update, len(updates))

	p.poll(context.Background())

	if got := store.saveCount(); got != 1 {
		t.Fatalf("save count=%d, want 1", got)
	}
	if got := store.offsetValue(); got != 16 {
		t.Fatalf("offset=%d, want 16", got)
	}
	if got := len(p.updates); got != len(updates) {
		t.Fatalf("queued updates=%d, want %d", got, len(updates))
	}
}

func TestPollerRetriesPendingOffsetBeforeFetchingMore(t *testing.T) {
	store := newPollerOffsetStore(7)
	store.failSave.Store(true)

	var fetches atomic.Int32
	tgClient := &pollerMockClient{
		getUpdatesFunc: func(_ context.Context, _ client.GetUpdatesJSONRequestBody) (*client.GetUpdatesResponse, error) {
			fetches.Add(1)

			return getUpdatesResponse([]client.Update{{UpdateId: 10}}), nil
		},
	}

	p, err := NewPoller(NewOptions(tgClient, WithOffsetStore(store)))
	if err != nil {
		t.Fatalf("NewPoller() unexpected error: %v", err)
	}

	if ok := p.poll(context.Background()); ok {
		t.Fatal("poll() success=true, want false when offset save fails")
	}
	if got := fetches.Load(); got != 1 {
		t.Fatalf("fetches=%d, want 1", got)
	}
	if got := store.offsetValue(); got != 7 {
		t.Fatalf("offset=%d, want 7", got)
	}
	if !p.hasPendingOffset || p.pendingOffset != 11 {
		t.Fatalf("pending offset=(%v, %d), want (true, 11)", p.hasPendingOffset, p.pendingOffset)
	}

	store.failSave.Store(false)

	if ok := p.poll(context.Background()); !ok {
		t.Fatal("poll() success=false, want true after pending offset save")
	}
	if got := fetches.Load(); got != 1 {
		t.Fatalf("fetches=%d, want still 1", got)
	}
	if got := store.offsetValue(); got != 11 {
		t.Fatalf("offset=%d, want 11", got)
	}
	if p.hasPendingOffset {
		t.Fatal("pending offset still set after successful retry")
	}
	if got := len(p.updates); got != 1 {
		t.Fatalf("queued updates=%d, want 1", got)
	}
}

func TestPollerGetUpdatesRequestOptions(t *testing.T) {
	store := newPollerOffsetStore(42)
	allowedUpdates := []string{"message", "callback_query"}
	var gotBody client.GetUpdatesJSONRequestBody
	tgClient := &pollerMockClient{
		getUpdatesFunc: func(_ context.Context, body client.GetUpdatesJSONRequestBody) (*client.GetUpdatesResponse, error) {
			gotBody = body

			return getUpdatesResponse(nil), nil
		},
	}

	p, err := NewPoller(NewOptions(
		tgClient,
		WithOffsetStore(store),
		WithTimeout(45*time.Second),
		WithLimit(25),
		WithAllowedUpdates(allowedUpdates),
	))
	if err != nil {
		t.Fatalf("NewPoller() unexpected error: %v", err)
	}

	p.poll(context.Background())

	if gotBody.Offset == nil || *gotBody.Offset != 42 {
		t.Fatalf("offset=%v, want 42", gotBody.Offset)
	}
	if gotBody.Timeout == nil || *gotBody.Timeout != 45 {
		t.Fatalf("timeout=%v, want 45", gotBody.Timeout)
	}
	if gotBody.Limit == nil || *gotBody.Limit != 25 {
		t.Fatalf("limit=%v, want 25", gotBody.Limit)
	}
	if gotBody.AllowedUpdates == nil {
		t.Fatal("allowed updates is nil")
	}
	if got, want := *gotBody.AllowedUpdates, allowedUpdates; !slices.Equal(got, want) {
		t.Fatalf("allowed updates=%v, want %v", got, want)
	}
}

func TestPollerRequestTimeout(t *testing.T) {
	store := newPollerOffsetStore(0)
	requestTimeout := 20 * time.Millisecond
	tgClient := &pollerMockClient{
		getUpdatesFunc: func(ctx context.Context, _ client.GetUpdatesJSONRequestBody) (*client.GetUpdatesResponse, error) {
			deadline, ok := ctx.Deadline()
			if !ok {
				t.Fatal("GetUpdates context has no deadline")
			}

			remaining := time.Until(deadline)
			if remaining <= 0 || remaining > requestTimeout {
				t.Fatalf("deadline remaining=%v, want within %v", remaining, requestTimeout)
			}

			<-ctx.Done()

			return nil, ctx.Err()
		},
	}

	p, err := NewPoller(NewOptions(
		tgClient,
		WithOffsetStore(store),
		WithRequestTimeout(requestTimeout),
	))
	if err != nil {
		t.Fatalf("NewPoller() unexpected error: %v", err)
	}

	start := time.Now()
	if ok := p.poll(context.Background()); ok {
		t.Fatal("poll() success=true, want false after request timeout")
	}
	if elapsed := time.Since(start); elapsed < requestTimeout || elapsed > time.Second {
		t.Fatalf("poll elapsed=%v, want at least %v and less than 1s", elapsed, requestTimeout)
	}
}

func TestPollerBufferSizeOption(t *testing.T) {
	tgClient := &pollerMockClient{
		getUpdatesFunc: func(_ context.Context, _ client.GetUpdatesJSONRequestBody) (*client.GetUpdatesResponse, error) {
			return getUpdatesResponse(nil), nil
		},
	}

	p, err := NewPoller(NewOptions(
		tgClient,
		WithOffsetStore(newPollerOffsetStore(0)),
		WithBufferSize(7),
	))
	if err != nil {
		t.Fatalf("NewPoller() unexpected error: %v", err)
	}
	if got := cap(p.updates); got != 7 {
		t.Fatalf("update buffer size=%d, want 7", got)
	}
}

func TestPollerRetryBackoffResetAfterSuccess(t *testing.T) {
	backoff := pollRetryBackoff{}

	if got := backoff.Next(); got != minRetryBackoff {
		t.Fatalf("first backoff=%v, want %v", got, minRetryBackoff)
	}
	if got := backoff.Next(); got != 2*minRetryBackoff {
		t.Fatalf("second backoff=%v, want %v", got, 2*minRetryBackoff)
	}

	backoff.Reset()

	if got := backoff.Next(); got != minRetryBackoff {
		t.Fatalf("backoff after reset=%v, want %v", got, minRetryBackoff)
	}
}

func getUpdatesResponse(updates []client.Update) *client.GetUpdatesResponse {
	return &client.GetUpdatesResponse{
		HTTPResponse: &http.Response{StatusCode: http.StatusOK},
		JSON200: &struct {
			Ok     client.GetUpdates200Ok `json:"ok"`
			Result []client.Update        `json:"result"`
		}{
			Ok:     true,
			Result: updates,
		},
	}
}

func waitForUpdateBufferLen(t *testing.T, ch <-chan client.Update, want int) {
	t.Helper()

	deadline := time.NewTimer(time.Second)
	defer deadline.Stop()

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	for {
		if got := len(ch); got == want {
			return
		}

		select {
		case <-deadline.C:
			t.Fatalf("update channel length=%d, want %d", len(ch), want)
		case <-ticker.C:
		}
	}
}
