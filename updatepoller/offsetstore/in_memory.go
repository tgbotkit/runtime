package offsetstore

import (
	"context"
	"sync/atomic"
)

// InMemoryOffsetStore is an in-memory implementation of OffsetStore.
type InMemoryOffsetStore struct {
	offset int64
}

// NewInMemoryOffsetStore creates a new InMemoryOffsetStore.
func NewInMemoryOffsetStore(initial int) *InMemoryOffsetStore {
	return &InMemoryOffsetStore{offset: int64(initial)}
}

func (s *InMemoryOffsetStore) Load(_ context.Context) (int, error) {
	return int(atomic.LoadInt64(&s.offset)), nil
}

func (s *InMemoryOffsetStore) Save(_ context.Context, offset int) error {
	atomic.StoreInt64(&s.offset, int64(offset))
	return nil
}
