package gcache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/amerkurev/gcache/pkg/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

type StringHasher struct{}

func (*StringHasher) Hash(v any) (string, error) {
	return fmt.Sprintf("%+v", v), nil
}

func TestCache_SetHasher(t *testing.T) {
	c := New[string, int64](MapStore(0))
	c.SetHasher(&StringHasher{})
}

type JSONMarshaler struct{}

func (m *JSONMarshaler) Marshal(v any) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (m *JSONMarshaler) Unmarshal(b []byte, v any) error {
	err := json.Unmarshal(b, v)
	if err != nil {
		return err
	}
	return nil
}

func TestCache_SetMarshaler(t *testing.T) {
	c := New[string, int64](MapStore(0))
	c.SetMarshaler(&JSONMarshaler{})
}

func TestNewMapCache_Int(t *testing.T) {
	ctx := context.Background()
	c := NewMapCache[int64, int64](0)

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

func TestNewMapCache_String(t *testing.T) {
	ctx := context.Background()
	c := NewMapCache[string, string](0)

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

func TestNewMapCache_Struct(t *testing.T) {
	ctx := context.Background()
	c := NewMapCache[int, *user](0)

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
