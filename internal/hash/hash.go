package hash

import (
	"crypto/sha256"
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
	"reflect"
)

// A HasherError represents an error from calling a Hash method.
type HasherError struct {
	Type reflect.Type
	Err  error
}

func (e *HasherError) Error() string {
	return "hash error for type " + e.Type.String() +
		": " + e.Err.Error()
}

// Unwrap returns the underlying error.
func (e *HasherError) Unwrap() error { return e.Err }

// Hasher is the interface implemented by types that can hash any value into string.
type Hasher interface {
	Hash(any) (string, error)
}

// MsgpackHasher is a default hasher that uses msgpack marshaling and
// hashing algorithm for create a hash of any value.
type MsgpackHasher struct{}

// Hash creates a hash string of any value.
func (*MsgpackHasher) Hash(v any) (string, error) {
	b, err := msgpack.Marshal(v)
	if err != nil {
		return "", &HasherError{Type: reflect.TypeOf(v), Err: err}
	}

	h := sha256.New()
	if _, err = h.Write(b); err != nil {
		return "", &HasherError{Type: reflect.TypeOf(v), Err: err}
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
