package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSQLiteStore(t *testing.T) {
	ctx := context.Background()
	db, err := sql.Open("sqlite3", "test.db")
	assert.Nil(t, err)
	s, err := SQLiteStore(ctx, db)
	assert.Nil(t, err)
	assert.NotNil(t, s)

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
