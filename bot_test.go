package runtime_test

import (
	"context"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tgbotkit/client"
	"github.com/tgbotkit/runtime"
	"github.com/tgbotkit/runtime/eventemitter"
	"github.com/tgbotkit/runtime/events"
)

// mockClient mocks the Telegram API client.
type mockClient struct {
	client.ClientWithResponsesInterface
	getMeFunc func(ctx context.Context, reqEditors ...client.RequestEditorFn) (*client.GetMeResponse, error)
}

func (m *mockClient) GetMeWithResponse(ctx context.Context, reqEditors ...client.RequestEditorFn) (*client.GetMeResponse, error) {
	if m.getMeFunc != nil {
		return m.getMeFunc(ctx, reqEditors...)
	}
	username := "TestBot"
	return &client.GetMeResponse{
		HTTPResponse: &http.Response{StatusCode: http.StatusOK},
		JSON200: &struct {
			Ok     client.GetMe200Ok `json:"ok"`
			Result client.User       `json:"result"`
		}{
			Ok: true,
			Result: client.User{
				Id:        12345,
				IsBot:     true,
				FirstName: "Test Bot",
				Username:  &username,
			},
		},
	}, nil
}

// mockUpdateSource mocks the UpdateSource interface.
type mockUpdateSource struct {
	ch chan client.Update
}

func (m *mockUpdateSource) UpdateChan() <-chan client.Update {
	return m.ch
}

func (m *mockUpdateSource) Start(_ context.Context) error {
	return nil
}

func (m *mockUpdateSource) Stop(_ context.Context) error {
	select {
	case <-m.ch:
		// already closed
	default:
		close(m.ch)
	}
	return nil
}

func TestNew(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cl := &mockClient{}
		opts := runtime.NewOptions("test-token", runtime.WithClient(cl))
		bot, err := runtime.New(opts)
		assert.NoError(t, err)
		assert.NotNil(t, bot)
		assert.Equal(t, cl, bot.Client())
		assert.NotNil(t, bot.EventEmitter())
		assert.NotNil(t, bot.Handlers())
	})

	t.Run("validation error", func(t *testing.T) {
		opts := runtime.NewOptions("") // Empty token, should fail validation
		bot, err := runtime.New(opts)
		assert.Error(t, err)
		assert.Nil(t, bot)
	})

	t.Run("getMe error", func(t *testing.T) {
		cl := &mockClient{}
		cl.getMeFunc = func(_ context.Context, _ ...client.RequestEditorFn) (*client.GetMeResponse, error) {
			return nil, assert.AnError
		}
		opts := runtime.NewOptions("test-token", runtime.WithClient(cl))
		bot, err := runtime.New(opts)
		assert.Error(t, err)
		assert.Nil(t, bot)
	})
}

func TestBot_Run(t *testing.T) {
	cl := &mockClient{}
	us := &mockUpdateSource{ch: make(chan client.Update, 1)}

	ee, _ := eventemitter.NewSync(eventemitter.NewOptions())

	// Track event emission
	var eventReceived atomic.Bool
	ee.AddListener(events.OnUpdate, eventemitter.ListenerFunc(func(_ context.Context, _ any) error {
		eventReceived.Store(true)
		return nil
	}))

	opts := runtime.NewOptions(
		"test-token",
		runtime.WithClient(cl),
		runtime.WithUpdateSource(us),
		runtime.WithEventEmitter(ee),
	)

	bot, err := runtime.New(opts)
	assert.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Push an update
	us.ch <- client.Update{UpdateId: 1}

	// Run in background
	errCh := make(chan error)
	go func() {
		errCh <- bot.Run(ctx)
	}()

	// Give it some time to process
	time.Sleep(50 * time.Millisecond)

	assert.True(t, eventReceived.Load(), "OnUpdate event should have been emitted")

	// Cancel context to stop Run
	cancel()
	err = <-errCh
	assert.ErrorIs(t, err, context.Canceled)
}

func TestBot_Run_SourceClose(t *testing.T) {
	cl := &mockClient{}
	us := &mockUpdateSource{ch: make(chan client.Update)}
	close(us.ch) // Immediately close channel

	opts := runtime.NewOptions(
		"test-token",
		runtime.WithClient(cl),
		runtime.WithUpdateSource(us),
	)

	bot, err := runtime.New(opts)
	assert.NoError(t, err)

	err = bot.Run(context.Background())
	assert.NoError(t, err, "Run should return nil when update channel is closed")
}
