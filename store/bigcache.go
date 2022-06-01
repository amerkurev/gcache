package store

import (
	"context"
	"errors"
	"github.com/allegro/bigcache/v3"
)

type bigcacheStore struct {
	bc *bigcache.BigCache
}

// BigcacheStore creates a Bigcache data store.
// See Bigcache docs https://github.com/allegro/bigcache.
func BigcacheStore(config bigcache.Config) (Store, error) {
	cache, err := bigcache.NewBigCache(config)
	if err != nil {
		return nil, err
	}

	return &bigcacheStore{
		bc: cache,
	}, nil
}

func (b *bigcacheStore) Get(_ context.Context, key string) ([]byte, error) {
	v, err := b.bc.Get(key)
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return nil, ErrNotFound
		}
	}
	return v, err
}

func (b *bigcacheStore) Set(_ context.Context, key string, data []byte) error {
	return b.bc.Set(key, data)
}

func (b *bigcacheStore) Delete(_ context.Context, key string) error {
	err := b.bc.Delete(key)
	if err != nil {
		// repeated delete must be safe
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return nil
		}
	}
	return err
}

func (b *bigcacheStore) Clear(_ context.Context) error {
	return b.bc.Reset()
}
