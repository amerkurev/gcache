package store

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeyNotFound(t *testing.T) {
	err := errors.New("some error")
	assert.Equal(t, ErrNotFound.Error(), "key not found")
	assert.False(t, errors.Is(err, ErrNotFound))
	assert.True(t, errors.As(err, &ErrNotFound))
}
