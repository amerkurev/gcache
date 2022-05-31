package gcache

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type StringHasher struct{}

func (*StringHasher) Hash(v any) (string, error) {
	return fmt.Sprintf("%+v", v), nil
}

func TestCache_SetHasher(t *testing.T) {
	c := New[string, int64]()
	c.SetHasher(&StringHasher{})
	k, err := c.Hash([]int{1, 2, 3})
	assert.Nil(t, err)
	assert.Equal(t, k, "[1 2 3]")
}
