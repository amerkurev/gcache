package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestSQLiteStore_Context(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	db, err := sql.Open("sqlite3", "test.db")
	assert.Nil(t, err)
	s, err := SQLiteStore(ctx, db)
	assert.Nil(t, err)
	assert.NotNil(t, s)

	t.Cleanup(func() {
		assert.Nil(t, db.Close())
	})

	err = s.Clear(ctx)
	assert.Nil(t, err)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for k := 0; k < 100_000; k++ {
			key := "a"
			err := s.Set(ctx, key, nil)
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

func TestSQLiteStore(t *testing.T) {
	ctx := context.Background()
	db, err := sql.Open("sqlite3", "test.db")
	assert.Nil(t, err)
	s, err := SQLiteStore(ctx, db)
	assert.Nil(t, err)
	assert.NotNil(t, s)

	t.Cleanup(func() {
		assert.Nil(t, db.Close())
	})

	err = s.Clear(ctx)
	assert.Nil(t, err)

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
