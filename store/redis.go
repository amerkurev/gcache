package store

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type redisStore struct {
	rdb *redis.Client
}

// RedisStore creates a Redis data store.
// See Redis docs https://github.com/go-redis/redis
func RedisStore(rdb *redis.Client) Store {
	return &redisStore{rdb}
}

func (r *redisStore) Get(ctx context.Context, key string) ([]byte, error) {
	v, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, ErrNotFound
	}
	return []byte(v), err
}

func (r *redisStore) Set(ctx context.Context, key string, data []byte) error {
	return r.rdb.Set(ctx, key, data, 0).Err()
}

func (r *redisStore) Delete(ctx context.Context, key string) error {
	return r.rdb.Del(ctx, key).Err()
}

func (r *redisStore) Clear(ctx context.Context) error {
	return r.rdb.FlushDB(ctx).Err()
}
