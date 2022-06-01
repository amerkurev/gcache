package store

import (
	"context"
	"errors"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRedisStore(t *testing.T) {
	addr := "127.0.0.1:6379"
	m := miniredis.NewMiniRedis()
	if err := m.StartAddr(addr); err != nil {
		t.Fatalf("could not start miniredis: %s", err)
		// not reached
	}
	t.Cleanup(m.Close)

	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	err := rdb.Ping(ctx).Err()
	assert.Nil(t, err)
	s := RedisStore(rdb)

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
