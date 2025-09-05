package kv

import (
	"sync"
	"time"
)

type kvEntry struct {
	value  string
	perm   bool
	expiry time.Time
}

// TODO: Make this better by using a sync.Map maybe
type Store struct {
	mu       sync.Mutex
	kv       map[string]kvEntry
	interval time.Duration
}

func NewStore(interval time.Duration) *Store {
	return &Store{
		kv: make(map[string]kvEntry),
	}
}

func (s *Store) Set(key, value string, duration time.Duration) {
	s.mu.Lock()
	s.kv[key] = kvEntry{
		value:  value,
		expiry: time.Now().Add(duration),
		perm:   duration.Milliseconds() == 0,
	}
	s.mu.Unlock()
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.Lock()
	entry, ok := s.kv[key]
	s.mu.Unlock()
	return entry.value, ok
}

func (s *Store) ExpiryHandler() {
	for {
		time.Sleep(s.interval)
		s.mu.Lock()
		for k, v := range s.kv {
			if !v.perm && time.Now().After(v.expiry) {
				delete(s.kv, k)
			}
		}
		s.mu.Unlock()
	}
}
