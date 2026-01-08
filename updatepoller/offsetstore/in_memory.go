// Package offsetstore provides implementations for storing the update offset.
package offsetstore

import (
	"context"
	"sync/atomic"

	"github.com/tgbotkit/runtime/updatepoller"
)

// InMemoryOffsetStore is an in-memory implementation of OffsetStore.
type InMemoryOffsetStore struct {
	offset int64
}

var _ updatepoller.OffsetStore = (*InMemoryOffsetStore)(nil)

// NewInMemoryOffsetStore creates a new InMemoryOffsetStore.
func NewInMemoryOffsetStore(initial int) *InMemoryOffsetStore {
	return &InMemoryOffsetStore{offset: int64(initial)}
}

// Load retrieves the current offset.
func (s *InMemoryOffsetStore) Load(_ context.Context) (int, error) {
	return int(atomic.LoadInt64(&s.offset)), nil
}

// Save stores the new offset.
func (s *InMemoryOffsetStore) Save(_ context.Context, offset int) error {
	atomic.StoreInt64(&s.offset, int64(offset))

	return nil
}
