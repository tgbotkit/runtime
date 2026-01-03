package offsetstore

import (
	"context"
	"testing"
)

func TestInMemoryOffsetStore(t *testing.T) {
	initial := 100
	store := NewInMemoryOffsetStore(initial)

	ctx := context.Background()

	// Test Load initial
	got, err := store.Load(ctx)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got != initial {
		t.Errorf("Load() got = %v, want %v", got, initial)
	}

	// Test Save and Load
	newOffset := 200
	if err := store.Save(ctx, newOffset); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	got, err = store.Load(ctx)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got != newOffset {
		t.Errorf("Load() got = %v, want %v", got, newOffset)
	}
}
