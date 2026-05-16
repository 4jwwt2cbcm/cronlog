package storage

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"cronlog/internal/logentry"
)

// ErrStoreNotFound is returned when no store file exists at the given path.
var ErrStoreNotFound = errors.New("store file not found")

// Store persists log entries to a JSON file on disk.
type Store struct {
	mu      sync.RWMutex
	path    string
	entries []*logentry.Entry
}

// New creates or loads a Store at the given file path.
func New(path string) (*Store, error) {
	s := &Store{path: path}
	if err := s.load(); err != nil && !errors.Is(err, ErrStoreNotFound) {
		return nil, err
	}
	return s, nil
}

// Add appends a new entry and flushes to disk.
func (s *Store) Add(e *logentry.Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, e)
	return s.save()
}

// All returns a copy of all stored entries.
func (s *Store) All() []*logentry.Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*logentry.Entry, len(s.entries))
	copy(out, s.entries)
	return out
}

// Replace overwrites all entries and flushes to disk.
func (s *Store) Replace(entries []*logentry.Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = entries
	return s.save()
}

// load reads entries from the JSON file.
func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return ErrStoreNotFound
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.entries)
}

// save writes entries to the JSON file atomically via a temp file.
func (s *Store) save() error {
	data, err := json.Marshal(s.entries)
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
