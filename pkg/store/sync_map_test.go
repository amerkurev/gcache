package store

import (
	"context"
	"errors"
	"github.com/amerkurev/gcache/internal/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapStore(t *testing.T) {
	s := MapStore(0)
	ctx := context.Background()
	key := "a"
	err := s.Set(ctx, key, nil)
	assert.Nil(t, err)

	key = "b"
	err = s.Set(ctx, key, []byte{1, 2, 3})
	assert.Nil(t, err)

	b, err := s.Get(ctx, key)
	assert.Nil(t, err)
	assert.Equal(t, b, []byte{1, 2, 3})

	err = s.Delete(ctx, key)
	assert.Nil(t, err)

	b, err = s.Get(ctx, key)
	assert.Nil(t, b)
	assert.True(t, errors.Is(err, store.ErrNotFound))

	err = s.Clear(ctx)
	assert.Nil(t, err)

	b, err = s.Get(ctx, "a")
	assert.Nil(t, b)
	assert.True(t, errors.Is(err, store.ErrNotFound))
}
