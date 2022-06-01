package store

import (
	"context"
	"errors"
	"github.com/allegro/bigcache/v3"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBigcacheStore(t *testing.T) {
	ctx := context.Background()
	store, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	assert.Nil(t, err)
	s := BigcacheStore(store)

	key := "a"
	err = s.Set(ctx, key, nil)
	assert.Nil(t, err)

	key = "b"
	err = s.Set(ctx, key, []byte{1, 2, 3})
	assert.Nil(t, err)

	b, err := s.Get(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, b, []byte{1, 2, 3})

	err = s.Delete(ctx, key)
	assert.Nil(t, err)

	// repeated delete must be safe
	err = s.Delete(ctx, key)
	assert.Nil(t, err)

	b, err = s.Get(ctx, key)
	assert.Nil(t, b)
	assert.True(t, errors.Is(err, ErrNotFound))

	err = s.Clear(ctx)
	assert.Nil(t, err)

	b, err = s.Get(ctx, "a")
	assert.Nil(t, b)
	assert.True(t, errors.Is(err, ErrNotFound))
}
