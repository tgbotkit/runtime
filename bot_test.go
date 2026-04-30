package runtime_test

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

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
	ch   chan client.Update
	once sync.Once

	closed atomic.Bool
}

func (m *mockUpdateSource) UpdateChan() <-chan client.Update {
	return m.ch
}

func (m *mockUpdateSource) Start(_ context.Context) error {
	return nil
}

func (m *mockUpdateSource) Stop(_ context.Context) error {
	m.once.Do(func() {
		if m.closed.Load() {
			return
		}
		close(m.ch)
		m.closed.Store(true)
	})

	return nil
}

func TestNew(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cl := &mockClient{}
		opts := runtime.NewOptions("test-token", runtime.WithClient(cl))
		bot, err := runtime.New(opts)
		if err != nil {
			t.Fatalf("New() unexpected error: %v", err)
		}
		if bot == nil {
			t.Fatal("New() bot is nil")
		}
		if bot.Client() != cl {
			t.Fatalf("Client() mismatch: got %T, want %T", bot.Client(), cl)
		}
		if bot.EventEmitter() == nil {
			t.Fatal("EventEmitter() is nil")
		}
		if bot.Handlers() == nil {
			t.Fatal("Handlers() is nil")
		}
	})

	t.Run("validation error", func(t *testing.T) {
		opts := runtime.NewOptions("")
		bot, err := runtime.New(opts)
		if err == nil {
			t.Fatal("New() error is nil, want validation error")
		}
		if bot != nil {
			t.Fatalf("New() bot=%v, want nil", bot)
		}
	})

	t.Run("getMe error", func(t *testing.T) {
		wantErr := errors.New("getMe failed")
		cl := &mockClient{
			getMeFunc: func(_ context.Context, _ ...client.RequestEditorFn) (*client.GetMeResponse, error) {
				return nil, wantErr
			},
		}

		opts := runtime.NewOptions("test-token", runtime.WithClient(cl))
		bot, err := runtime.New(opts)
		if err == nil {
			t.Fatal("New() error is nil, want error")
		}
		if !errors.Is(err, wantErr) {
			t.Fatalf("New() error=%v, want wrapped %v", err, wantErr)
		}
		if bot != nil {
			t.Fatalf("New() bot=%v, want nil", bot)
		}
	})

	t.Run("startup timeout", func(t *testing.T) {
		cl := &mockClient{
			getMeFunc: func(ctx context.Context, _ ...client.RequestEditorFn) (*client.GetMeResponse, error) {
				<-ctx.Done()
				return nil, ctx.Err()
			},
		}

		opts := runtime.NewOptions(
			"test-token",
			runtime.WithClient(cl),
			runtime.WithStartupTimeout(10*time.Millisecond),
		)

		bot, err := runtime.New(opts)
		if err == nil {
			t.Fatal("New() error is nil, want timeout error")
		}
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("New() error=%v, want wrapped %v", err, context.DeadlineExceeded)
		}
		if bot != nil {
			t.Fatalf("New() bot=%v, want nil", bot)
		}
	})
}

func TestBot_Run(t *testing.T) {
	t.Run("processes updates and exits on context cancel", func(t *testing.T) {
		cl := &mockClient{}
		us := &mockUpdateSource{ch: make(chan client.Update, 1)}

		ee, err := eventemitter.NewSync(eventemitter.NewOptions())
		if err != nil {
			t.Fatalf("NewSync() unexpected error: %v", err)
		}

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
		if err != nil {
			t.Fatalf("New() unexpected error: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		us.ch <- client.Update{UpdateId: 1}

		errCh := make(chan error, 1)
		go func() {
			errCh <- bot.Run(ctx)
		}()

		time.Sleep(50 * time.Millisecond)

		if !eventReceived.Load() {
			t.Fatal("OnUpdate event was not emitted")
		}

		cancel()
		err = <-errCh
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Run() error=%v, want %v", err, context.Canceled)
		}
	})

	t.Run("returns ErrUpdateSourceClosed when source channel closes", func(t *testing.T) {
		cl := &mockClient{}
		closedCh := make(chan client.Update)
		close(closedCh)

		us := &mockUpdateSource{ch: closedCh}
		us.closed.Store(true)
		opts := runtime.NewOptions(
			"test-token",
			runtime.WithClient(cl),
			runtime.WithUpdateSource(us),
		)

		bot, err := runtime.New(opts)
		if err != nil {
			t.Fatalf("New() unexpected error: %v", err)
		}

		err = bot.Run(context.Background())
		if !errors.Is(err, runtime.ErrUpdateSourceClosed) {
			t.Fatalf("Run() error=%v, want %v", err, runtime.ErrUpdateSourceClosed)
		}
	})
}
