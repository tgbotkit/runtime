package updatepoller_test

import (
	"context"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)

	// Create a context and cancel it
	startCtx, cancelStart := context.WithCancel(context.Background())
	err = p.Start(startCtx)
	assert.NoError(t, err)

	// Wait for some polls
	time.Sleep(25 * time.Millisecond)
	assert.Greater(t, tgClient.pollCount.Load(), int32(0))

	// Canceling startCtx should NOT stop the poller now (as it uses internal background context)
	cancelStart()

	// Wait and check it's still polling
	countAfterCancel := tgClient.pollCount.Load()
	time.Sleep(25 * time.Millisecond)
	assert.Greater(t, tgClient.pollCount.Load(), countAfterCancel, "Poller should still be running after startCtx was canceled")

	// Stop explicitly
	err = p.Stop(context.Background())
	assert.NoError(t, err)

	// Verify it stopped
	countAfterStop := tgClient.pollCount.Load()
	time.Sleep(25 * time.Millisecond)
	assert.Equal(t, countAfterStop, tgClient.pollCount.Load(), "Poller should have stopped after Stop() call")

	// Verify we can restart
	err = p.Start(context.Background())
	assert.NoError(t, err)
	time.Sleep(25 * time.Millisecond)
	assert.Greater(t, tgClient.pollCount.Load(), countAfterStop)

	// Final stop
	err = p.Stop(context.Background())
	assert.NoError(t, err)
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
	assert.NoError(t, err)

	// Start
	err = p.Start(context.Background())
	assert.NoError(t, err)
	time.Sleep(20 * time.Millisecond)
	assert.Greater(t, tgClient.pollCount.Load(), int32(0))

	// Stop
	err = p.Stop(context.Background())
	assert.NoError(t, err)
	countAfterStop := tgClient.pollCount.Load()
	time.Sleep(20 * time.Millisecond)
	assert.Equal(t, countAfterStop, tgClient.pollCount.Load())

	// Restart
	err = p.Start(context.Background())
	assert.NoError(t, err)
	time.Sleep(20 * time.Millisecond)
	assert.Greater(t, tgClient.pollCount.Load(), countAfterStop)

	// Final stop
	err = p.Stop(context.Background())
	assert.NoError(t, err)
}

func TestPoller_StopTimeout(t *testing.T) {
	tgClient := &mockClient{}
	store := &mockOffsetStore{}

	// Create an API that takes longer than the stop timeout
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
	assert.NoError(t, err)

	err = p.Start(context.Background())
	assert.NoError(t, err)

	// Give it a moment to start the first poll
	time.Sleep(20 * time.Millisecond)

	// Stop with a very short timeout
	stopCtx, cancelStop := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancelStop()

	start := time.Now()
	err = p.Stop(stopCtx)
	elapsed := time.Since(start)

	assert.Error(t, err)
	if err != nil {
		assert.Contains(t, err.Error(), "context deadline exceeded")
	}
	assert.GreaterOrEqual(t, elapsed, 5*time.Millisecond)
}
