package gcache

import (
	"context"
	"errors"
	"github.com/allegro/bigcache/v3"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestMapCache_Int64(t *testing.T) {
	ctx := context.Background()
	c, err := NewMapCache[int64, int64](0)
	assert.Nil(t, err)

	err = c.Set(1, 1000)
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
	c, err := NewMapCache[string, string](0)
	assert.Nil(t, err)

	err = c.Set("a", "some value")
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
	c, err := NewMapCache[int, *user](0)
	assert.Nil(t, err)

	err = c.Set(100, &user{"John"})
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

func TestMapCache_Concurrency(t *testing.T) {
	c, err := NewMapCache[int, int](0)
	assert.Nil(t, err)

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
	c, err := NewBigCache[int, int](bigcache.DefaultConfig(10 * time.Minute))
	assert.Nil(t, err)

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

func TestCacheStats(t *testing.T) {
	c, err := NewMapCache[int, int](0)
	assert.Nil(t, err)
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
}
