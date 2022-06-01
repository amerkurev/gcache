package hasher

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestHasherError_Error(t *testing.T) {
	err := errors.New("some error")
	e := &Error{
		Type: reflect.TypeOf(int64(1)),
		Err:  err,
	}
	assert.Equal(t, e.Error(), "error when hashing object of type int64: some error")
	assert.True(t, errors.Is(e, err))
	assert.True(t, errors.As(e, &err))
}

func TestMsgpackHasher_Hash(t *testing.T) {
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
}
