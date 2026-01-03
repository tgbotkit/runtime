package eventemitter

import (
	"context"
	"errors"
	"testing"
)

func TestEventEmitter_Emit_BreakOnError(t *testing.T) {
	ee, err := NewSync(NewOptions(WithStopOnError(true)))
	if err != nil {
		t.Fatalf("failed to create event emitter: %v", err)
	}
	ctx := context.Background()

	var callCount int
	errDummy := errors.New("dummy error")

	ee.AddListener("test", func(ctx context.Context, payload any) error {
		callCount++
		return errDummy
	})

	ee.AddListener("test", func(ctx context.Context, payload any) error {
		callCount++
		return nil
	})

	ee.Emit(ctx, "test", nil)

	if callCount != 1 {
		t.Errorf("expected 1 listener to be called, got %d", callCount)
	}
}

func TestEventEmitter_RemoveListener(t *testing.T) {
	ee, err := NewSync(NewOptions())
	if err != nil {
		t.Fatalf("failed to create event emitter: %v", err)
	}
	ctx := context.Background()

	var callCount int
	listener := func(ctx context.Context, payload any) error {
		callCount++
		return nil
	}

	ee.AddListener("test", listener)
	ee.Emit(ctx, "test", nil)

	if callCount != 1 {
		t.Errorf("expected 1 listener to be called, got %d", callCount)
	}

	ee.RemoveListener("test", listener)
	ee.Emit(ctx, "test", nil)

	if callCount != 1 {
		t.Errorf("expected 1 listener to be called, got %d", callCount)
	}
}
