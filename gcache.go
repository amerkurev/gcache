package gcache

import (
	"github.com/amerkurev/gcache/internal/hasher"
	"github.com/amerkurev/gcache/internal/marshaler"
)

// Cache represents the interface for all caches
type Cache[K comparable, V any] interface {
	SetHasher(hasher.Hasher)
	SetMarshaler(marshaler.Marshaler)
}

type cache[K comparable, V any] struct {
	hasher.Hasher
	marshaler.Marshaler
}

// New creates a new instance of cache object.
func New[K comparable, V any]() Cache[K, V] {
	return &cache[K, V]{
		Hasher:    &hasher.MsgpackHasher{},
		Marshaler: &marshaler.MsgpackMarshaler{},
	}
}

func (c *cache[K, V]) SetHasher(h hasher.Hasher) {
	c.Hasher = h
}

func (c *cache[K, V]) SetMarshaler(m marshaler.Marshaler) {
	c.Marshaler = m
}
