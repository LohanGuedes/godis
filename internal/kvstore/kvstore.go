package kvstore

import (
	"sync"
	"time"
)

type kvEntry struct {
	value  string
	expiry time.Time
}

// TODO: Make this better by using a sync.Map maybe
type Store struct {
	mu sync.Mutex
	kv map[string]kvEntry
}

func NewStore() *Store {
	return &Store{
		kv: make(map[string]kvEntry),
	}
}

func (s *Store) Set(key, value string) {
	s.mu.Lock()
	s.kv[key] = kvEntry{value: value}
	s.mu.Unlock()
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.Lock()
	entry, ok := s.kv[key]
	s.mu.Unlock()
	return entry.value, ok
}
