package updatepoller_test

import (
	"context"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime/updatepoller"
)

type mockClient struct {
	client.ClientWithResponsesInterface
	pollCount      atomic.Int32
	getUpdatesFunc func(ctx context.Context, body client.GetUpdatesJSONRequestBody) (*client.GetUpdatesResponse, error)
}

func (m *mockClient) GetUpdatesWithResponse(ctx context.Context, body client.GetUpdatesJSONRequestBody, _ ...client.RequestEditorFn) (*client.GetUpdatesResponse, error) {
	m.pollCount.Add(1)
	if m.getUpdatesFunc != nil {
		return m.getUpdatesFunc(ctx, body)
	}

	return &client.GetUpdatesResponse{
		HTTPResponse: &http.Response{StatusCode: http.StatusOK},
		JSON200: &struct {
			Ok     client.GetUpdates200Ok `json:"ok"`
			Result []client.Update        `json:"result"`
		}{
			Ok:     true,
			Result: []client.Update{},
		},
	}, nil
}

type mockOffsetStore struct {
	offset int
}

func (m *mockOffsetStore) Load(_ context.Context) (int, error) {
	return m.offset, nil
}

func (m *mockOffsetStore) Save(_ context.Context, offset int) error {
	m.offset = offset
	return nil
}

func TestPoller_LifecycleContext(t *testing.T) {
	tgClient := &mockClient{}
	store := &mockOffsetStore{}

	opts := updatepoller.NewOptions(
		tgClient,
		updatepoller.WithOffsetStore(store),
		updatepoller.WithPollingInterval(10*time.Millisecond),
	)

	p, err := updatepoller.NewPoller(opts)
	if err != nil {
		t.Fatalf("NewPoller() unexpected error: %v", err)
	}

	startCtx, cancelStart := context.WithCancel(context.Background())
	err = p.Start(startCtx)
	if err != nil {
		t.Fatalf("Start() unexpected error: %v", err)
	}

	time.Sleep(25 * time.Millisecond)
	if got := tgClient.pollCount.Load(); got <= 0 {
		t.Fatalf("poll count=%d, want > 0", got)
	}

	cancelStart()

	countAfterCancel := tgClient.pollCount.Load()
	time.Sleep(25 * time.Millisecond)
	if got := tgClient.pollCount.Load(); got <= countAfterCancel {
		t.Fatalf("poll count after cancel=%d, want > %d", got, countAfterCancel)
	}

	err = p.Stop(context.Background())
	if err != nil {
		t.Fatalf("Stop() unexpected error: %v", err)
	}

	countAfterStop := tgClient.pollCount.Load()
	time.Sleep(25 * time.Millisecond)
	if got := tgClient.pollCount.Load(); got != countAfterStop {
		t.Fatalf("poll count after stop=%d, want %d", got, countAfterStop)
	}

	err = p.Start(context.Background())
	if err != nil {
		t.Fatalf("restart Start() unexpected error: %v", err)
	}
	time.Sleep(25 * time.Millisecond)
	if got := tgClient.pollCount.Load(); got <= countAfterStop {
		t.Fatalf("poll count after restart=%d, want > %d", got, countAfterStop)
	}

	err = p.Stop(context.Background())
	if err != nil {
		t.Fatalf("final Stop() unexpected error: %v", err)
	}
}

func TestPoller_Restart(t *testing.T) {
	tgClient := &mockClient{}
	store := &mockOffsetStore{}

	opts := updatepoller.NewOptions(
		tgClient,
		updatepoller.WithOffsetStore(store),
		updatepoller.WithPollingInterval(10*time.Millisecond),
	)

	p, err := updatepoller.NewPoller(opts)
	if err != nil {
		t.Fatalf("NewPoller() unexpected error: %v", err)
	}

	err = p.Start(context.Background())
	if err != nil {
		t.Fatalf("Start() unexpected error: %v", err)
	}
	time.Sleep(20 * time.Millisecond)
	if got := tgClient.pollCount.Load(); got <= 0 {
		t.Fatalf("poll count=%d, want > 0", got)
	}

	err = p.Stop(context.Background())
	if err != nil {
		t.Fatalf("Stop() unexpected error: %v", err)
	}
	countAfterStop := tgClient.pollCount.Load()
	time.Sleep(20 * time.Millisecond)
	if got := tgClient.pollCount.Load(); got != countAfterStop {
		t.Fatalf("poll count after stop=%d, want %d", got, countAfterStop)
	}

	err = p.Start(context.Background())
	if err != nil {
		t.Fatalf("restart Start() unexpected error: %v", err)
	}
	time.Sleep(20 * time.Millisecond)
	if got := tgClient.pollCount.Load(); got <= countAfterStop {
		t.Fatalf("poll count after restart=%d, want > %d", got, countAfterStop)
	}

	err = p.Stop(context.Background())
	if err != nil {
		t.Fatalf("final Stop() unexpected error: %v", err)
	}
}

func TestPoller_StopTimeout(t *testing.T) {
	tgClient := &mockClient{}
	store := &mockOffsetStore{}

	tgClient.getUpdatesFunc = func(_ context.Context, _ client.GetUpdatesJSONRequestBody) (*client.GetUpdatesResponse, error) {
		time.Sleep(50 * time.Millisecond)
		return &client.GetUpdatesResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				Ok     client.GetUpdates200Ok `json:"ok"`
				Result []client.Update        `json:"result"`
			}{
				Ok:     true,
				Result: []client.Update{},
			},
		}, nil
	}

	opts := updatepoller.NewOptions(
		tgClient,
		updatepoller.WithOffsetStore(store),
		updatepoller.WithPollingInterval(10*time.Millisecond),
	)

	p, err := updatepoller.NewPoller(opts)
	if err != nil {
		t.Fatalf("NewPoller() unexpected error: %v", err)
	}

	err = p.Start(context.Background())
	if err != nil {
		t.Fatalf("Start() unexpected error: %v", err)
	}

	time.Sleep(20 * time.Millisecond)

	stopCtx, cancelStop := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancelStop()

	start := time.Now()
	err = p.Stop(stopCtx)
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("Stop() error is nil, want deadline error")
	}
	if !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Fatalf("Stop() error=%v, want context deadline exceeded", err)
	}
	if elapsed < 5*time.Millisecond {
		t.Fatalf("Stop() elapsed=%v, want >= %v", elapsed, 5*time.Millisecond)
	}
}

func TestPoller_StartCanceledContext(t *testing.T) {
	tgClient := &mockClient{}
	store := &mockOffsetStore{}

	opts := updatepoller.NewOptions(
		tgClient,
		updatepoller.WithOffsetStore(store),
		updatepoller.WithPollingInterval(10*time.Millisecond),
	)

	p, err := updatepoller.NewPoller(opts)
	if err != nil {
		t.Fatalf("NewPoller() unexpected error: %v", err)
	}

	startCtx, cancel := context.WithCancel(context.Background())
	cancel()

	err = p.Start(startCtx)
	if err == nil {
		t.Fatal("Start() error is nil, want context canceled")
	}
	if err != context.Canceled {
		t.Fatalf("Start() error=%v, want %v", err, context.Canceled)
	}

	time.Sleep(25 * time.Millisecond)
	if got := tgClient.pollCount.Load(); got != 0 {
		t.Fatalf("poll count=%d, want 0", got)
	}
}
