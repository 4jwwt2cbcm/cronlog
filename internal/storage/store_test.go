package storage_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"cronlog/internal/logentry"
	"cronlog/internal/storage"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "store.json")
}

func makeEntry(job, msg string, isErr bool) *logentry.Entry {
	return &logentry.Entry{
		Job:       job,
		Message:   msg,
		Timestamp: time.Now(),
		IsError:   isErr,
	}
}

func TestStoreAddAndAll(t *testing.T) {
	s, err := storage.New(tempPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	e1 := makeEntry("backup", "started", false)
	e2 := makeEntry("backup", "failed", true)

	if err := s.Add(e1); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if err := s.Add(e2); err != nil {
		t.Fatalf("Add: %v", err)
	}

	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestStorePersistence(t *testing.T) {
	path := tempPath(t)
	s, _ := storage.New(path)
	_ = s.Add(makeEntry("sync", "ok", false))

	// Reload from disk.
	s2, err := storage.New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if len(s2.All()) != 1 {
		t.Fatalf("expected 1 entry after reload, got %d", len(s2.All()))
	}
}

func TestStoreReplace(t *testing.T) {
	s, _ := storage.New(tempPath(t))
	for i := 0; i < 5; i++ {
		_ = s.Add(makeEntry("job", "msg", false))
	}

	if err := s.Replace([]*logentry.Entry{makeEntry("job", "only", false)}); err != nil {
		t.Fatalf("Replace: %v", err)
	}
	if len(s.All()) != 1 {
		t.Fatalf("expected 1 entry after replace, got %d", len(s.All()))
	}
}

func TestStoreNotFoundIsOK(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	s, err := storage.New(path)
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(s.All()) != 0 {
		t.Fatal("expected empty store for missing file")
	}
}

func TestStoreAtomicWrite(t *testing.T) {
	path := tempPath(t)
	s, _ := storage.New(path)
	_ = s.Add(makeEntry("job", "msg", false))

	// Ensure no leftover .tmp file.
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Fatal("tmp file should not exist after successful save")
	}
}
