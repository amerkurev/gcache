package gcache

import (
	"github.com/amerkurev/gcache/internal/hash"
)

// Cache represents the interface for all caches
type Cache[K comparable, V any] interface {
	Hash(any) (string, error)
	SetHasher(hash.Hasher)
}

type cache[K comparable, V any] struct {
	hash.Hasher
}

// New creates a new instance of gcache.
func New[K comparable, V any]() Cache[K, V] {
	return &cache[K, V]{&hash.MsgpackHasher{}}
}

// SetHasher sets a hasher that provides a hash function for cache key create.
func (c *cache[K, V]) SetHasher(hasher hash.Hasher) {
	c.Hasher = hasher
}
