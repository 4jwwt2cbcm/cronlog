// Package storage provides persistent storage for cronlog log entries.
//
// Entries are serialised as a JSON array on disk. Writes are made atomic
// by writing to a temporary file and renaming it into place, so a crash
// during a write cannot corrupt the existing data.
//
// Basic usage:
//
//	s, err := storage.New("/var/lib/cronlog/entries.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Append a new entry.
//	_ = s.Add(entry)
//
//	// Query with filters.
//	results := s.Query(storage.Filter{
//		Job:     "backup",
//		OnlyErr: true,
//		Since:   time.Now().Add(-24 * time.Hour),
//	})
//
//	// Remove entries older than 30 days.
//	_ = s.Prune(time.Now().Add(-30 * 24 * time.Hour))
//
// The Store is safe for concurrent use by multiple goroutines.
package storage
