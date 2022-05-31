package store

import (
	"context"
	"errors"
)

// Store is the interface implemented by types that can be data storage for cache.
type Store interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, data []byte) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
}

// ErrNotFound indicates that key not found in the store.
var ErrNotFound = errors.New("key not found")
