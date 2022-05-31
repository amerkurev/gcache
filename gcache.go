package gcache

import (
	"context"
	"github.com/amerkurev/gcache/internal/hasher"
	"github.com/amerkurev/gcache/internal/marshaler"
	"github.com/amerkurev/gcache/internal/store"
	impl "github.com/amerkurev/gcache/store"
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

	SetHasher(hasher.Hasher)
	SetMarshaler(marshaler.Marshaler)
}

type cache[KeyType comparable, ValueType any] struct {
	hasher.Hasher
	marshaler.Marshaler
	store.Store
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
		return
	}
	b, err := c.Store.Get(ctx, k)
	if err != nil {
		return
	}
	err = c.Unmarshal(b, &value)
	return
}

func (c *cache[K, V]) SetWithContext(ctx context.Context, key K, value V) error {
	k, err := c.Hash(key)
	if err != nil {
		return err
	}
	v, err := c.Marshal(value)
	if err != nil {
		return err
	}
	return c.Store.Set(ctx, k, v)
}

func (c *cache[K, V]) DeleteWithContext(ctx context.Context, key K) error {
	k, err := c.Hash(key)
	if err != nil {
		return err
	}
	return c.Store.Delete(ctx, k)
}

func (c *cache[K, V]) ClearWithContext(ctx context.Context) error {
	return c.Store.Clear(ctx)
}

func (c *cache[K, V]) SetHasher(h hasher.Hasher) {
	c.Hasher = h
}

func (c *cache[K, V]) SetMarshaler(m marshaler.Marshaler) {
	c.Marshaler = m
}

// New creates a new instance of cache object.
func New[K comparable, V any](s store.Store) Cache[K, V] {
	return &cache[K, V]{
		Hasher:    &hasher.MsgpackHasher{},
		Marshaler: &marshaler.MsgpackMarshaler{},
		Store:     s,
	}
}

// NewMapCache creates a new instance of cache object with the MapStore as a data store.
// MapStore is like a Go map but is safe for concurrent use by multiple goroutines.
func NewMapCache[K comparable, V any](size int) Cache[K, V] {
	return New[K, V](impl.MapStore(size))
}

// ErrNotFound indicates that key not found in the cache.
var ErrNotFound = store.ErrNotFound
