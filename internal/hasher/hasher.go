package hasher

import (
	"crypto/sha256"
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
	"reflect"
)

// Hasher is the interface implemented by types that can hash any value into string.
type Hasher interface {
	Hash(any) (string, error)
}

// Error represents an error from calling a Hash method.
type Error struct {
	Type reflect.Type
	Err  error
}

func (e *Error) Error() string {
	return "error when hashing object of type " + e.Type.String() +
		": " + e.Err.Error()
}

// Unwrap returns the underlying error.
func (e *Error) Unwrap() error { return e.Err }

// MsgpackHasher is a default hasher that uses msgpack marshaling and
// hashing algorithm for create a hash of value.
type MsgpackHasher struct{}

// Hash creates a hash string of any value.
func (*MsgpackHasher) Hash(v any) (string, error) {
	b, err := msgpack.Marshal(v)
	if err != nil {
		return "", &Error{Type: reflect.TypeOf(v), Err: err}
	}

	h := sha256.New()
	if _, err = h.Write(b); err != nil {
		return "", &Error{Type: reflect.TypeOf(v), Err: err}
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
