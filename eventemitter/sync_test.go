package eventemitter

import (
	"context"
	"errors"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestEventEmitter_Emit_BreakOnError(t *testing.T) {
	ee, err := NewSync(NewOptions(WithStopOnError(true)))
	if err != nil {
		t.Fatalf("failed to create event emitter: %v", err)
	}
	ctx := context.Background()

	var callCount int
	errDummy := errors.New("dummy error")

	ee.AddListener("test", ListenerFunc(func(_ context.Context, _ any) error {
		callCount++
		return errDummy
	}))

	ee.AddListener("test", ListenerFunc(func(_ context.Context, _ any) error {
		callCount++
		return nil
	}))

	ee.Emit(ctx, "test", nil)

	if callCount != 1 {
		t.Errorf("expected 1 listener to be called, got %d", callCount)
	}
}

func TestEventEmitter_Middleware(t *testing.T) {
	ee, err := NewSync(NewOptions())
	if err != nil {
		t.Fatalf("failed to create event emitter: %v", err)
	}
	ctx := context.Background()

	var middlewareCalled bool
	ee.Use("test", MiddlewareFunc(func(next Listener) Listener {
		return ListenerFunc(func(ctx context.Context, payload any) error {
			middlewareCalled = true
			return next.Handle(ctx, payload)
		})
	}))

	var handlerCalled bool
	ee.AddListener("test", ListenerFunc(func(_ context.Context, _ any) error {
		handlerCalled = true
		return nil
	}))

	ee.Emit(ctx, "test", nil)

	if !middlewareCalled {
		t.Error("expected middleware to be called")
	}
	if !handlerCalled {
		t.Error("expected handler to be called")
	}
}

func TestEventEmitter_RemoveListener(t *testing.T) {
	ee, err := NewSync(NewOptions())
	if err != nil {
		t.Fatalf("failed to create event emitter: %v", err)
	}
	ctx := context.Background()

	var callCount int
	listener := ListenerFunc(func(_ context.Context, _ any) error {
		callCount++
		return nil
	})

	unsubscribe := ee.AddListener("test", listener)
	if ee.ListenerCount("test") != 1 {
		t.Errorf("expected 1 listener, got %d", ee.ListenerCount("test"))
	}

	ee.Emit(ctx, "test", nil)

	if callCount != 1 {
		t.Errorf("expected 1 listener to be called, got %d", callCount)
	}

	unsubscribe()
	if ee.ListenerCount("test") != 0 {
		t.Errorf("expected 0 listeners, got %d", ee.ListenerCount("test"))
	}

	ee.Emit(ctx, "test", nil)

	if callCount != 1 {
		t.Errorf("expected 1 listener to be called, got %d", callCount)
	}
}

func TestEventEmitter_GlobMatching(t *testing.T) {
	ee, err := NewSync(NewOptions())
	if err != nil {
		t.Fatalf("failed to create event emitter: %v", err)
	}
	ctx := context.Background()

	var specificCalled, globCalled bool

	ee.AddListener("test.event", ListenerFunc(func(_ context.Context, _ any) error {
		specificCalled = true
		return nil
	}))

	ee.AddListener("test.*", ListenerFunc(func(_ context.Context, _ any) error {
		globCalled = true
		return nil
	}))

	ee.Emit(ctx, "test.event", nil)

	if !specificCalled {
		t.Error("expected specific listener to be called")
	}
	if !globCalled {
		t.Error("expected glob listener to be called")
	}

	if ee.ListenerCount("test.event") != 2 {
		t.Errorf("expected 2 listeners for test.event (1 specific + 1 glob), got %d", ee.ListenerCount("test.event"))
	}
}

func TestEventEmitter_DeterministicMixedPatternListeners(t *testing.T) {
	ee, err := NewSync(NewOptions())
	if err != nil {
		t.Fatalf("failed to create event emitter: %v", err)
	}
	ctx := context.Background()

	var calls []string

	ee.AddListener("test.*", ListenerFunc(func(_ context.Context, _ any) error {
		calls = append(calls, "wildcard-first")
		return nil
	}))
	ee.AddListener("test.event", ListenerFunc(func(_ context.Context, _ any) error {
		calls = append(calls, "exact-second")
		return nil
	}))
	ee.AddListener("test.*", ListenerFunc(func(_ context.Context, _ any) error {
		calls = append(calls, "wildcard-third")
		return nil
	}))
	ee.AddListener("test.event", ListenerFunc(func(_ context.Context, _ any) error {
		calls = append(calls, "exact-fourth")
		return nil
	}))

	want := []string{"wildcard-first", "exact-second", "wildcard-third", "exact-fourth"}
	for range 100 {
		calls = calls[:0]

		ee.Emit(ctx, "test.event", nil)

		assertCallOrder(t, calls, want)
	}
}

func TestEventEmitter_DeterministicMixedPatternMiddleware(t *testing.T) {
	ee, err := NewSync(NewOptions())
	if err != nil {
		t.Fatalf("failed to create event emitter: %v", err)
	}
	ctx := context.Background()

	var calls []string
	recordMiddleware := func(name string) Middleware {
		return MiddlewareFunc(func(next Listener) Listener {
			return ListenerFunc(func(ctx context.Context, payload any) error {
				calls = append(calls, name)
				return next.Handle(ctx, payload)
			})
		})
	}

	ee.Use("test.*", recordMiddleware("wildcard-first"))
	ee.Use("test.event", recordMiddleware("exact-second"))
	ee.Use("test.*", recordMiddleware("wildcard-third"))
	ee.Use("test.event", recordMiddleware("exact-fourth"))
	ee.AddListener("test.event", ListenerFunc(func(_ context.Context, _ any) error {
		calls = append(calls, "listener")
		return nil
	}))

	want := []string{"wildcard-first", "exact-second", "wildcard-third", "exact-fourth", "listener"}
	for range 100 {
		calls = calls[:0]

		ee.Emit(ctx, "test.event", nil)

		assertCallOrder(t, calls, want)
	}
}

func TestEventEmitter_Once(t *testing.T) {
	ee, err := NewSync(NewOptions())
	if err != nil {
		t.Fatalf("failed to create event emitter: %v", err)
	}
	ctx := context.Background()

	var callCount int
	ee.Once("test", ListenerFunc(func(_ context.Context, _ any) error {
		callCount++
		return nil
	}))

	ee.Emit(ctx, "test", nil)
	ee.Emit(ctx, "test", nil)

	if callCount != 1 {
		t.Errorf("expected listener to be called once, got %d", callCount)
	}
}

func TestEventEmitter_Once_ConcurrentEmit(t *testing.T) {
	ee, err := NewSync(NewOptions())
	if err != nil {
		t.Fatalf("failed to create event emitter: %v", err)
	}

	var callCount atomic.Int32
	releaseListener := make(chan struct{})
	ee.Once("test", ListenerFunc(func(_ context.Context, _ any) error {
		callCount.Add(1)
		<-releaseListener

		return nil
	}))

	const goroutines = 100
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for range goroutines {
		go func() {
			defer wg.Done()
			<-start
			ee.Emit(context.Background(), "test", nil)
		}()
	}

	close(start)
	waitForCallCount(t, &callCount, 1)
	close(releaseListener)
	wg.Wait()

	if got := callCount.Load(); got != 1 {
		t.Fatalf("call count=%d, want 1", got)
	}
}

func TestEventEmitter_ErrBreak(t *testing.T) {
	ee, err := NewSync(NewOptions())
	if err != nil {
		t.Fatalf("failed to create event emitter: %v", err)
	}
	ctx := context.Background()

	var firstCalled, secondCalled bool

	ee.AddListener("test", ListenerFunc(func(_ context.Context, _ any) error {
		firstCalled = true
		return ErrBreak
	}))

	ee.AddListener("test", ListenerFunc(func(_ context.Context, _ any) error {
		secondCalled = true
		return nil
	}))

	ee.Emit(ctx, "test", nil)

	if !firstCalled {
		t.Error("expected first listener to be called")
	}
	if secondCalled {
		t.Error("expected second listener NOT to be called due to ErrBreak")
	}
}

func TestEventEmitter_ErrBreakMixedPatternRegistrationOrder(t *testing.T) {
	ee, err := NewSync(NewOptions())
	if err != nil {
		t.Fatalf("failed to create event emitter: %v", err)
	}
	ctx := context.Background()

	var calls []string

	ee.AddListener("test.event", ListenerFunc(func(_ context.Context, _ any) error {
		calls = append(calls, "exact-break")
		return ErrBreak
	}))
	ee.AddListener("test.*", ListenerFunc(func(_ context.Context, _ any) error {
		calls = append(calls, "wildcard")
		return nil
	}))

	ee.Emit(ctx, "test.event", nil)

	assertCallOrder(t, calls, []string{"exact-break"})
}

func TestEventEmitter_ErrorHandler(t *testing.T) {
	var errorCaught error
	errorHandler := func(_ string, err error) {
		errorCaught = err
	}

	ee, err := NewSync(NewOptions(WithErrorHandler(errorHandler)))
	if err != nil {
		t.Fatalf("failed to create event emitter: %v", err)
	}
	ctx := context.Background()

	dummyErr := errors.New("dummy error")
	ee.AddListener("test", ListenerFunc(func(_ context.Context, _ any) error {
		return dummyErr
	}))

	ee.Emit(ctx, "test", nil)

	if errorCaught != dummyErr {
		t.Errorf("expected error handler to catch error, got %v", errorCaught)
	}
}

func TestEventEmitter_RemoveAllListeners(t *testing.T) {
	ee, err := NewSync(NewOptions())
	if err != nil {
		t.Fatalf("failed to create event emitter: %v", err)
	}

	ee.AddListener("test", ListenerFunc(func(_ context.Context, _ any) error { return nil }))
	ee.AddListener("test", ListenerFunc(func(_ context.Context, _ any) error { return nil }))

	if ee.ListenerCount("test") != 2 {
		t.Errorf("expected 2 listeners, got %d", ee.ListenerCount("test"))
	}

	ee.RemoveAllListeners("test")

	if ee.ListenerCount("test") != 0 {
		t.Errorf("expected 0 listeners, got %d", ee.ListenerCount("test"))
	}
}

func TestGenericOn(t *testing.T) {
	ee, err := NewSync(NewOptions())
	if err != nil {
		t.Fatalf("failed to create event emitter: %v", err)
	}
	ctx := context.Background()

	type MyEvent struct {
		Data string
	}

	var receivedData string
	On[MyEvent](ee, "test", func(_ context.Context, payload *MyEvent) error {
		receivedData = payload.Data
		return nil
	})

	ee.Emit(ctx, "test", &MyEvent{Data: "hello"})

	if receivedData != "hello" {
		t.Errorf("expected 'hello', got '%s'", receivedData)
	}

	// Test with wrong type
	receivedData = ""
	ee.Emit(ctx, "test", "wrong type")
	if receivedData != "" {
		t.Error("expected listener not to be called with wrong type")
	}
}

func TestEventEmitter_StopOnError_False(t *testing.T) {
	ee, err := NewSync(NewOptions(WithStopOnError(false)))
	if err != nil {
		t.Fatalf("failed to create event emitter: %v", err)
	}
	ctx := context.Background()

	var callCount int
	errDummy := errors.New("dummy error")

	ee.AddListener("test", ListenerFunc(func(_ context.Context, _ any) error {
		callCount++
		return errDummy
	}))

	ee.AddListener("test", ListenerFunc(func(_ context.Context, _ any) error {
		callCount++
		return nil
	}))

	ee.Emit(ctx, "test", nil)

	if callCount != 2 {
		t.Errorf("expected both listeners to be called, got %d", callCount)
	}
}

func TestEventEmitter_Once_WithErrors(t *testing.T) {
	ctx := context.Background()

	t.Run("once with ErrBreak", func(t *testing.T) {
		ee, _ := NewSync(NewOptions())
		var callCount int
		ee.Once("test", ListenerFunc(func(_ context.Context, _ any) error {
			callCount++
			return ErrBreak
		}))
		ee.Emit(ctx, "test", nil)
		ee.Emit(ctx, "test", nil)
		if callCount != 1 {
			t.Errorf("expected 1 call, got %d", callCount)
		}
	})

	t.Run("once with error and stopOnError true", func(t *testing.T) {
		ee, _ := NewSync(NewOptions(WithStopOnError(true)))
		var callCount int
		ee.Once("test", ListenerFunc(func(_ context.Context, _ any) error {
			callCount++
			return errors.New("fail")
		}))
		ee.Emit(ctx, "test", nil)
		ee.Emit(ctx, "test", nil)
		if callCount != 1 {
			t.Errorf("expected 1 call, got %d", callCount)
		}
	})

	t.Run("once with error and stopOnError false", func(t *testing.T) {
		ee, _ := NewSync(NewOptions(WithStopOnError(false)))
		var callCount int
		ee.Once("test", ListenerFunc(func(_ context.Context, _ any) error {
			callCount++
			return errors.New("fail")
		}))
		ee.Emit(ctx, "test", nil)
		ee.Emit(ctx, "test", nil)
		if callCount != 1 {
			t.Errorf("expected 1 call, got %d", callCount)
		}
	})
}

func TestEventEmitter_InvalidGlob(t *testing.T) {
	ee, _ := NewSync(NewOptions())
	ctx := context.Background()

	var called bool
	// Invalid pattern according to path.Match: a [ with no closing ]
	ee.AddListener("[", ListenerFunc(func(_ context.Context, _ any) error {
		called = true
		return nil
	}))

	ee.Emit(ctx, "test", nil)
	if called {
		t.Error("expected invalid glob listener not to be called")
	}

	// ListenerCount returns 0 because matching fails
	if ee.ListenerCount("[") != 0 {
		t.Errorf("expected 0 matching listeners for invalid pattern, got %d", ee.ListenerCount("["))
	}
}

func TestEventEmitter_Once_Unsubscribe(t *testing.T) {
	ee, _ := NewSync(NewOptions())
	ctx := context.Background()

	var callCount int
	unsubscribe := ee.Once("test", ListenerFunc(func(_ context.Context, _ any) error {
		callCount++
		return nil
	}))

	unsubscribe()
	ee.Emit(ctx, "test", nil)

	if callCount != 0 {
		t.Errorf("expected 0 calls after unsubscribe, got %d", callCount)
	}
}

func TestEventEmitter_UnsubscribeTwice(t *testing.T) {
	ee, _ := NewSync(NewOptions())
	unsubscribe := ee.AddListener("test", ListenerFunc(func(_ context.Context, _ any) error { return nil }))

	unsubscribe()
	if ee.ListenerCount("test") != 0 {
		t.Errorf("expected 0 listeners, got %d", ee.ListenerCount("test"))
	}

	// Should not panic or cause issues
	unsubscribe()
}

func TestEventEmitter_RemoveAllListeners_NoExist(_ *testing.T) {
	ee, _ := NewSync(NewOptions())
	// Should not panic
	ee.RemoveAllListeners("non-existent")
}

func assertCallOrder(t *testing.T, got, want []string) {
	t.Helper()

	if !slices.Equal(got, want) {
		t.Fatalf("call order=%v, want %v", got, want)
	}
}

func waitForCallCount(t *testing.T, count *atomic.Int32, want int32) {
	t.Helper()

	deadline := time.NewTimer(time.Second)
	defer deadline.Stop()

	ticker := time.NewTicker(time.Millisecond)
	defer ticker.Stop()

	for {
		if count.Load() >= want {
			return
		}

		select {
		case <-deadline.C:
			t.Fatalf("call count=%d, want at least %d", count.Load(), want)
		case <-ticker.C:
		}
	}
}
