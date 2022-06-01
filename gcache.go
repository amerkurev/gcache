package gcache

import (
	"context"
	"errors"
	"github.com/amerkurev/gcache/internal/hasher"
	"github.com/amerkurev/gcache/internal/marshaler"
	"github.com/amerkurev/gcache/internal/stats"
	"github.com/amerkurev/gcache/store"
)

// Cache represents the interface for all caches
type Cache[KeyType comparable, ValueType any] interface {
	Get(KeyType) (ValueType, error)
	Set(KeyType, ValueType) error
	Delete(KeyType) error
	Clear() error

	GetWithContext(context.Context, KeyType) (ValueType, error)
	SetWithContext(context.Context, KeyType, ValueType) error
	DeleteWithContext(context.Context, KeyType) error
	ClearWithContext(context.Context) error

	UseStats()
	ResetStats()
	Stats() (stats.Stats, bool)
}

type cache[KeyType comparable, ValueType any] struct {
	hasher.Hasher
	marshaler.Marshaler
	store.Store

	*stats.SyncStats
	useStats bool
}

func (c *cache[K, V]) Get(key K) (V, error) {
	return c.GetWithContext(context.Background(), key)
}

func (c *cache[K, V]) Set(key K, value V) error {
	return c.SetWithContext(context.Background(), key, value)
}

func (c *cache[K, V]) Delete(key K) error {
	return c.DeleteWithContext(context.Background(), key)
}

func (c *cache[K, V]) Clear() error {
	return c.ClearWithContext(context.Background())
}

func (c *cache[K, V]) GetWithContext(ctx context.Context, key K) (value V, err error) {
	k, err := c.Hash(key)
	if err != nil {
		if c.useStats {
			c.ErrRead()
		}
		return
	}

	b, err := c.Store.Get(ctx, k)
	if err != nil {
		if c.useStats {
			if errors.Is(err, ErrNotFound) {
				c.IncRead(false, 0)
			} else {
				c.ErrRead()
			}
		}
		return
	}

	err = c.Unmarshal(b, &value)
	if c.useStats {
		if err != nil {
			c.ErrRead()
		} else {
			c.IncRead(true, len(b))
		}
	}
	return
}

func (c *cache[K, V]) SetWithContext(ctx context.Context, key K, value V) error {
	k, err := c.Hash(key)
	if err != nil {
		if c.useStats {
			c.ErrWrite()
		}
		return err
	}

	v, err := c.Marshal(value)
	if err != nil {
		if c.useStats {
			c.ErrWrite()
		}
		return err
	}

	err = c.Store.Set(ctx, k, v)
	if c.useStats {
		if err != nil {
			c.ErrWrite()
		} else {
			c.IncWrite(len(v))
		}
	}
	return err
}

func (c *cache[K, V]) DeleteWithContext(ctx context.Context, key K) error {
	k, err := c.Hash(key)
	if err != nil {
		if c.useStats {
			c.ErrDelete()
		}
		return err
	}

	err = c.Store.Delete(ctx, k)
	if c.useStats {
		if err != nil {
			c.ErrDelete()
		} else {
			c.IncDelete()
		}
	}
	return err
}

func (c *cache[K, V]) ClearWithContext(ctx context.Context) error {
	err := c.Store.Clear(ctx)
	if c.useStats {
		if err != nil {
			c.ErrClear()
		} else {
			c.IncClear()
		}
	}
	return err
}

func (c *cache[K, V]) UseStats() {
	c.useStats = true
}

func (c *cache[K, V]) ResetStats() {
	c.SyncStats.Reset()
}

func (c *cache[K, V]) Stats() (stats.Stats, bool) {
	return c.SyncStats.Snapshot(), c.useStats
}

// New creates a new instance of cache object.
func New[K comparable, V any](s store.Store) Cache[K, V] {
	return &cache[K, V]{
		Hasher:    &hasher.MsgpackHasher{},
		Marshaler: &marshaler.MsgpackMarshaler{},
		Store:     s,
		SyncStats: &stats.SyncStats{},
	}
}

// NewMapCache creates a new instance of cache object with the MapStore as a data store.
// MapStore is like a Go map but is safe for concurrent use by multiple goroutines.
func NewMapCache[K comparable, V any](size int) Cache[K, V] {
	return New[K, V](store.MapStore(size))
}

// ErrNotFound indicates that key not found in the cache.
var ErrNotFound = store.ErrNotFound
