package gcache

import (
	"context"
	"database/sql"
	"errors"
	"github.com/alicebob/miniredis/v2"
	"github.com/allegro/bigcache/v3"
	"github.com/amerkurev/gcache/internal/hasher"
	"github.com/amerkurev/gcache/internal/marshaler"
	"github.com/amerkurev/gcache/store"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestMapCache_Int64(t *testing.T) {
	ctx := context.Background()
	c := New[int64, int64](store.MapStore(0))

	err := c.Set(1, 1000)
	assert.Nil(t, err)
	err = c.SetWithContext(ctx, 2, 2000)
	assert.Nil(t, err)

	v, err := c.GetWithContext(ctx, 1)
	assert.Nil(t, err)
	assert.Equal(t, v, int64(1000))
	v, err = c.Get(2)
	assert.Nil(t, err)
	assert.Equal(t, v, int64(2000))

	err = c.Delete(2)
	assert.Nil(t, err)
	err = c.DeleteWithContext(ctx, 2)
	assert.Nil(t, err)

	v, err = c.Get(2)
	assert.Equal(t, v, int64(0))
	assert.True(t, errors.Is(err, ErrNotFound))

	err = c.Clear()
	assert.Nil(t, err)
	err = c.ClearWithContext(ctx)
	assert.Nil(t, err)

	v, err = c.GetWithContext(ctx, 1)
	assert.Equal(t, v, int64(0))
	assert.True(t, errors.Is(err, ErrNotFound))
}

func TestMapCache_String(t *testing.T) {
	ctx := context.Background()
	c := New[string, string](store.MapStore(0))

	err := c.Set("a", "some value")
	assert.Nil(t, err)
	err = c.SetWithContext(ctx, "b", "another value")
	assert.Nil(t, err)

	v, err := c.GetWithContext(ctx, "a")
	assert.Nil(t, err)
	assert.Equal(t, v, "some value")
	v, err = c.Get("b")
	assert.Nil(t, err)
	assert.Equal(t, v, "another value")

	err = c.Delete("b")
	assert.Nil(t, err)
	err = c.DeleteWithContext(ctx, "b")
	assert.Nil(t, err)

	v, err = c.Get("b")
	assert.Equal(t, v, "")
	assert.True(t, errors.Is(err, ErrNotFound))

	err = c.Clear()
	assert.Nil(t, err)
	err = c.ClearWithContext(ctx)
	assert.Nil(t, err)

	v, err = c.GetWithContext(ctx, "a")
	assert.Equal(t, v, "")
	assert.True(t, errors.Is(err, ErrNotFound))
}

type user struct {
	Name string
}

func TestMapCache_Struct(t *testing.T) {
	ctx := context.Background()
	c := New[int, *user](store.MapStore(0))

	err := c.Set(100, &user{"John"})
	assert.Nil(t, err)
	err = c.SetWithContext(ctx, 200, &user{"Mary"})
	assert.Nil(t, err)

	v, err := c.GetWithContext(ctx, 100)
	assert.Nil(t, err)
	assert.Equal(t, v, &user{"John"})
	v, err = c.Get(200)
	assert.Nil(t, err)
	assert.Equal(t, v, &user{"Mary"})

	err = c.Delete(200)
	assert.Nil(t, err)
	err = c.DeleteWithContext(ctx, 200)
	assert.Nil(t, err)

	v, err = c.Get(200)
	assert.Nil(t, v)
	assert.True(t, errors.Is(err, ErrNotFound))

	err = c.Clear()
	assert.Nil(t, err)
	err = c.ClearWithContext(ctx)
	assert.Nil(t, err)

	v, err = c.GetWithContext(ctx, 100)
	assert.Nil(t, v)
	assert.True(t, errors.Is(err, ErrNotFound))
}

func TestMapCache_Unsupported(t *testing.T) {
	var hashError *hasher.Error
	c := New[complex128, int](store.MapStore(0))
	c.UseStats()

	_, err := c.Get(complex(10, 30))
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &hashError))

	err = c.Set(complex(10, 30), 100)
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &hashError))

	err = c.Delete(complex(10, 30))
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &hashError))

	s, ok := c.Stats()
	assert.True(t, ok)
	assert.Equal(t, s.ErrReadCount, 1)
	assert.Equal(t, s.ErrWriteCount, 1)
	assert.Equal(t, s.ErrDeleteCount, 1)
}

func TestMapCache_Concurrency(t *testing.T) {
	c := New[int, int](store.MapStore(0))

	goroutines := 10
	items := 10_000

	var wg sync.WaitGroup

	// concurrency write
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for k := 0; k < items; k++ {
				// write
				err := c.Set(k, k*k)
				assert.Nil(t, err)

				// read
				v, err := c.Get(k + 1)
				if err != nil {
					assert.True(t, errors.Is(err, ErrNotFound))
				} else {
					assert.Equal(t, v, (k+1)*(k+1))
				}

				// delete
				err = c.Delete(k - 10)
				assert.Nil(t, err)

				// clear
				if k%1000 == 0 {
					err := c.Clear()
					assert.Nil(t, err)
				}
			}
		}()
	}
	wg.Wait()
}

func TestBigCache_Concurrency(t *testing.T) {
	s, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	assert.Nil(t, err)

	c := New[int, int](store.BigcacheStore(s))

	goroutines := 10
	items := 10_000

	var wg sync.WaitGroup

	// concurrency write
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for k := 0; k < items; k++ {
				// write
				err := c.Set(k, k*k)
				assert.Nil(t, err)

				// read
				v, err := c.Get(k + 1)
				if err != nil {
					assert.True(t, errors.Is(err, ErrNotFound))
				} else {
					assert.Equal(t, v, (k+1)*(k+1))
				}

				// delete
				err = c.Delete(k - 10)
				assert.Nil(t, err)
				err = c.Delete(k - 10)
				assert.Nil(t, err)

				// clear
				if k%1000 == 0 {
					err := c.Clear()
					assert.Nil(t, err)
				}
			}
		}()
	}
	wg.Wait()
}

func TestRedisCache_Concurrency(t *testing.T) {
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

	c := New[int, int](store.RedisStore(rdb))

	goroutines := 10
	items := 1000

	var wg sync.WaitGroup

	// concurrency write
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for k := 0; k < items; k++ {
				// write
				err := c.SetWithContext(ctx, k, k*k)
				assert.Nil(t, err)

				// read
				v, err := c.GetWithContext(ctx, k+1)
				if err != nil {
					assert.True(t, errors.Is(err, ErrNotFound))
				} else {
					assert.Equal(t, v, (k+1)*(k+1))
				}

				// delete
				err = c.DeleteWithContext(ctx, k-10)
				assert.Nil(t, err)
				err = c.Delete(k - 10)
				assert.Nil(t, err)

				// clear
				if k%1000 == 0 {
					err := c.ClearWithContext(ctx)
					assert.Nil(t, err)
				}
			}
		}()
	}
	wg.Wait()
}

func TestSQLiteCache_Concurrency(t *testing.T) {
	db, err := sql.Open("sqlite3", "test.db")
	assert.Nil(t, err)
	t.Cleanup(func() {
		assert.Nil(t, db.Close())
	})

	ctx := context.Background()
	s, err := store.SQLiteStore(ctx, db)
	assert.Nil(t, err)
	c := New[int, int](s)

	goroutines := 10
	items := 100

	var wg sync.WaitGroup

	// concurrency write
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for k := 0; k < items; k++ {
				// write
				err := c.SetWithContext(ctx, k, k*k)
				assert.Nil(t, err)

				// read
				v, err := c.GetWithContext(ctx, k+1)
				if err != nil {
					assert.True(t, errors.Is(err, ErrNotFound))
				} else {
					assert.Equal(t, v, (k+1)*(k+1))
				}

				// delete
				err = c.DeleteWithContext(ctx, k-10)
				assert.Nil(t, err)
				err = c.Delete(k - 10)
				assert.Nil(t, err)

				// clear
				if k%1000 == 0 {
					err := c.ClearWithContext(ctx)
					assert.Nil(t, err)
				}
			}
		}()
	}
	wg.Wait()
}

func TestRedisCache_Context(t *testing.T) {
	addr := "127.0.0.1:6379"
	m := miniredis.NewMiniRedis()
	if err := m.StartAddr(addr); err != nil {
		t.Fatalf("could not start miniredis: %s", err)
		// not reached
	}
	t.Cleanup(m.Close)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	c := New[int, int](store.RedisStore(rdb))

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for k := 0; k < 100_000; k++ {
			err := c.SetWithContext(ctx, k, k*k)
			if err != nil {
				assert.True(t, errors.Is(err, context.Canceled))
				break
			}
		}
		wg.Done()
	}()

	cancel()
	wg.Wait()

	assert.Equal(t, <-ctx.Done(), struct{}{})
	assert.True(t, errors.Is(ctx.Err(), context.Canceled))
}

func TestCacheStats(t *testing.T) {
	c := New[int, int](store.MapStore(0))
	c.UseStats()

	goroutines := 10
	items := 10_000
	bytesCount := 494600

	var wg sync.WaitGroup

	// concurrency write
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for k := 0; k < items; k++ {
				err := c.Set(k, k*k)
				assert.Nil(t, err)
			}
		}()
	}
	wg.Wait()

	s, ok := c.Stats()
	assert.True(t, ok)
	assert.Equal(t, s.WriteCount, goroutines*items)
	assert.Equal(t, s.WriteBytes, bytesCount)
	assert.Equal(t, s.ErrWriteCount, 0)

	// concurrency hit read
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for k := 0; k < items; k++ {
				v, err := c.Get(k)
				assert.Equal(t, k*k, v)
				assert.Nil(t, err)
			}
		}()
	}
	wg.Wait()

	s, ok = c.Stats()
	assert.True(t, ok)
	assert.Equal(t, s.Hits, goroutines*items)
	assert.Equal(t, s.ReadCount, goroutines*items)
	assert.Equal(t, s.ReadBytes, bytesCount)
	assert.Equal(t, s.ErrReadCount, 0)

	// concurrency delete
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for k := 0; k < items; k++ {
				err := c.Delete(k)
				assert.Nil(t, err)
			}
		}()
	}
	wg.Wait()

	s, ok = c.Stats()
	assert.True(t, ok)
	assert.Equal(t, s.DeleteCount, goroutines*items)
	assert.Equal(t, s.ErrDeleteCount, 0)

	// concurrency miss read
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for k := 0; k < items; k++ {
				v, err := c.Get(k)
				assert.Equal(t, v, 0)
				assert.True(t, errors.Is(err, ErrNotFound))
			}
		}()
	}
	wg.Wait()

	s, ok = c.Stats()
	assert.True(t, ok)
	assert.Equal(t, s.Hits, goroutines*items)
	assert.Equal(t, s.Miss, goroutines*items)
	assert.Equal(t, s.ReadCount, 2*goroutines*items)
	assert.Equal(t, s.ReadBytes, bytesCount)
	assert.Equal(t, s.ErrReadCount, 0)

	// concurrency clear
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := c.Clear()
			assert.Nil(t, err)
		}()
	}
	wg.Wait()

	s, ok = c.Stats()
	assert.True(t, ok)
	assert.Equal(t, s.ClearCount, goroutines)
	assert.Equal(t, s.ErrClearCount, 0)

	c.ResetStats()

	s, ok = c.Stats()
	assert.True(t, ok)
	assert.Equal(t, s.ReadBytes, 0)
	assert.Equal(t, s.WriteBytes, 0)
}

func TestCacheStats_Error(t *testing.T) {
	var marshalError *marshaler.MarshalError
	c := New[string, complex128](store.MapStore(0))
	c.UseStats()

	err := c.Set("a", complex(10, 30))
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &marshalError))

	s, ok := c.Stats()
	assert.True(t, ok)
	assert.Equal(t, s.WriteCount, 0)
	assert.Equal(t, s.ErrWriteCount, 1)
}
