package updatepoller

import "context"

// OffsetStore is an interface for storing and loading the update offset.
type OffsetStore interface {
	// Load loads the offset from the store.
	Load(ctx context.Context) (int, error)
	// Save saves the offset to the store.
	Save(ctx context.Context, offset int) error
}
