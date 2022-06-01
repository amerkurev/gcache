package marshaler

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestMarshalError(t *testing.T) {
	err := errors.New("some marshal error")
	e := &MarshalError{
		Type: reflect.TypeOf(int64(1)),
		Err:  err,
	}
	assert.Equal(t, e.Error(), "cannot serialize object of type int64: some marshal error")
	assert.True(t, errors.Is(e, err))
	assert.True(t, errors.As(e, &err))
}

func TestUnmarshalError(t *testing.T) {
	err := errors.New("some unmarshal error")
	e := &UnmarshalError{
		Type: reflect.TypeOf(int64(1)),
		Err:  err,
	}
	assert.Equal(t, e.Error(), "cannot deserialize object of type int64: some unmarshal error")
	assert.True(t, errors.Is(e, err))
	assert.True(t, errors.As(e, &err))
}

func TestMsgpackMarshaler_Marshal(t *testing.T) {
	m := &MsgpackMarshaler{}
	b, err := m.Marshal(map[string]struct{}{})
	assert.Nil(t, err)
	assert.Equal(t, b, []byte{0x80})

	b, err = m.Marshal(100)
	assert.Nil(t, err)
	assert.Equal(t, b, []byte{0x64})

	b, err = m.Marshal("some value")
	assert.Nil(t, err)
	assert.Equal(t, b, []byte{0xaa, 0x73, 0x6f, 0x6d, 0x65, 0x20, 0x76, 0x61, 0x6c, 0x75, 0x65})
}

func TestMsgpackMarshaler_Unmarshal(t *testing.T) {
	m := &MsgpackMarshaler{}
	var v map[string]struct{}
	err := m.Unmarshal([]byte{0x80}, &v)
	assert.Nil(t, err)
	assert.Equal(t, v, map[string]struct{}{})

	var i int
	err = m.Unmarshal([]byte{0x64}, &i)
	assert.Nil(t, err)
	assert.Equal(t, i, 100)

	var s string
	err = m.Unmarshal([]byte{0xaa, 0x73, 0x6f, 0x6d, 0x65, 0x20, 0x76, 0x61, 0x6c, 0x75, 0x65}, &s)
	assert.Nil(t, err)
	assert.Equal(t, s, "some value")

	err = m.Unmarshal([]byte{}, s)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "cannot deserialize object of type string: msgpack: Decode(non-pointer string)")
}
