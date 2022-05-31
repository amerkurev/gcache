package marshaler

import (
	"github.com/vmihailenco/msgpack/v5"
	"reflect"
)

// Marshaler is the interface implemented by types that can marshal and unmarshal any value into bytes.
type Marshaler interface {
	Marshal(any) ([]byte, error)
	Unmarshal([]byte, any) error
}

// MarshalError represents an error from calling a Marshal method.
type MarshalError struct {
	Type reflect.Type
	Err  error
}

func (e *MarshalError) Error() string {
	return "cannot serialize object of type " + e.Type.String() +
		": " + e.Err.Error()
}

// Unwrap returns the underlying error.
func (e *MarshalError) Unwrap() error { return e.Err }

// UnmarshalError represents an error from calling an Unmarshal method.
type UnmarshalError struct {
	Type reflect.Type
	Err  error
}

func (e *UnmarshalError) Error() string {
	return "cannot deserialize object of type " + e.Type.String() +
		": " + e.Err.Error()
}

// Unwrap returns the underlying error.
func (e *UnmarshalError) Unwrap() error { return e.Err }

// MsgpackMarshaler is a default marshaler that uses msgpack for marshaling and unmarshaling.
type MsgpackMarshaler struct{}

// Marshal returns the binary encoding of any value.
func (m *MsgpackMarshaler) Marshal(v any) ([]byte, error) {
	b, err := msgpack.Marshal(v)
	if err != nil {
		return nil, &MarshalError{Type: reflect.TypeOf(v), Err: err}
	}
	return b, nil
}

// Unmarshal decodes the binary data.
func (m *MsgpackMarshaler) Unmarshal(b []byte, v any) error {
	err := msgpack.Unmarshal(b, v)
	if err != nil {
		return &UnmarshalError{Type: reflect.TypeOf(v), Err: err}
	}
	return nil
}
