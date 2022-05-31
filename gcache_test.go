package gcache

import (
	"encoding/json"
	"fmt"
	"testing"
)

type StringHasher struct{}

func (*StringHasher) Hash(v any) (string, error) {
	return fmt.Sprintf("%+v", v), nil
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

func TestCache_SetHasher(t *testing.T) {
	c := New[string, int64]()
	c.SetHasher(&StringHasher{})
	c.SetMarshaler(&JSONMarshaler{})
}
