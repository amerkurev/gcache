package store

import (
	"context"
	"sync"
)

type mapStore struct {
	mx sync.RWMutex
	m  map[string][]byte
}

// MapStore creates a store that is like a Go map but is safe for concurrent use by multiple goroutines.
func MapStore(size int) Store {
	m := make(map[string][]byte, size)
	return &mapStore{m: m}
}

func (s *mapStore) Get(_ context.Context, key string) ([]byte, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
	v, ok := s.m[key]
	if !ok {
		return nil, ErrNotFound
	}
	return v, nil
}

func (s *mapStore) Set(_ context.Context, key string, data []byte) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.m[key] = data
	return nil
}

func (s *mapStore) Delete(_ context.Context, key string) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	delete(s.m, key)
	return nil
}

func (s *mapStore) Clear(_ context.Context) error {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.m = make(map[string][]byte)
	return nil
}
