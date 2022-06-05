package hasher

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestHasherError(t *testing.T) {
	err := errors.New("some error")
	e := &Error{
		Type: reflect.TypeOf(int64(1)),
		Err:  err,
	}
	assert.Equal(t, e.Error(), "error when hashing object of type int64: some error")
	assert.True(t, errors.Is(e, err))
	assert.True(t, errors.As(e, &err))
}

func TestMsgpackHasher(t *testing.T) {
	h := &MsgpackHasher{}
	k, err := h.Hash(map[string]struct{}{})
	assert.Nil(t, err)
	assert.Equal(t, k, "76be8b528d0075f7aae98d6fa57a6d3c83ae480a8469e668d7b0af968995ac71")

	k, err = h.Hash(100)
	assert.Nil(t, err)
	assert.Equal(t, k, "18ac3e7343f016890c510e93f935261169d9e3f565436429830faf0934f4f8e4")

	k, err = h.Hash("some key")
	assert.Nil(t, err)
	assert.Equal(t, k, "0dc44df765b1ef70e8b5069777b6cb177fdeef0cc977327b9e19e4a3dad24818")

	var m map[string]struct{}
	k, err = h.Hash(m)
	assert.Nil(t, err)
	assert.Equal(t, k, "e4ff5e7d7a7f08e9800a3e25cb774533cb20040df30b6ba10f956f9acd0eb3f7")

	_, err = h.Hash(nil)
	assert.Nil(t, err)
	assert.Equal(t, k, "e4ff5e7d7a7f08e9800a3e25cb774533cb20040df30b6ba10f956f9acd0eb3f7")

	unsupportedTypes := map[string]any{
		"msgpack: Encode(unsupported func())":     func() {},
		"msgpack: Encode(unsupported chan int)":   make(chan int),
		"msgpack: Encode(unsupported complex128)": complex(10, 11),
	}

	for s, anyType := range unsupportedTypes {
		_, err = h.Hash(anyType)
		assert.NotNil(t, err)
		e, ok := err.(*Error)
		assert.True(t, ok)
		assert.Equal(t, e.Error(), fmt.Sprintf("error when hashing object of type %T: %s", anyType, s))
		assert.Equal(t, e.Unwrap().Error(), s)
	}
}
