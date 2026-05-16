package storage

import (
	"time"

	"cronlog/internal/logentry"
)

// Filter holds optional criteria for querying entries.
type Filter struct {
	Job      string    // exact job name match; empty means any
	OnlyErr  bool      // if true, return only error entries
	Since    time.Time // zero means no lower bound
	Until    time.Time // zero means no upper bound
	MaxCount int       // 0 means no limit
}

// Query returns entries from the store that match f.
func (s *Store) Query(f Filter) []*logentry.Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var out []*logentry.Entry
	for _, e := range s.entries {
		if f.Job != "" && e.Job != f.Job {
			continue
		}
		if f.OnlyErr && !e.IsError {
			continue
		}
		if !f.Since.IsZero() && e.Timestamp.Before(f.Since) {
			continue
		}
		if !f.Until.IsZero() && e.Timestamp.After(f.Until) {
			continue
		}
		out = append(out, e)
		if f.MaxCount > 0 && len(out) >= f.MaxCount {
			break
		}
	}
	return out
}
