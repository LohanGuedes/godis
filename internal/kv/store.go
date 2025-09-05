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
	defer s.mu.Unlock()
	entry, ok := s.kv[key]
	if !entry.perm && time.Now().After(entry.expiry) {
		delete(s.kv, key)
		return entry.value, ok
	}
	return entry.value, ok
}

// TODO: Improve this from O(n) to O(Log(n)) (using a binary-tree)
func (s *Store) ExpiryHandler() {
	ticker := time.NewTicker(s.interval)
	for range ticker.C {
		s.mu.Lock()
		for k, v := range s.kv {
			if !v.perm && time.Now().After(v.expiry) {
				delete(s.kv, k)
			}
		}
		s.mu.Unlock()
	}
}
