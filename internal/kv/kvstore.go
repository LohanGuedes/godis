package kv

import (
	"time"
)

type kvEntry struct {
	value  string
	expiry time.Time
}

// TODO: Make this better by using a sync.Map maybe
type Store struct {
	kv       map[string]kvEntry
	interval time.Duration
}

func NewStore(interval time.Duration) *Store {
	return &Store{
		kv: make(map[string]kvEntry),
	}
}

func (s *Store) Set(key, value string, duration time.Duration) {
	s.kv[key] = kvEntry{
		value:  value,
		expiry: time.Now().Add(duration),
	}
}

func (s *Store) Get(key string) (string, bool) {
	entry, ok := s.kv[key]
	return entry.value, ok
}

func (s *Store) ExpiryHandler() {
	time.Sleep(s.interval)
	for k, v := range s.kv {
		if time.Now().After(v.expiry) {
			delete(s.kv, k)
		}
	}
}
