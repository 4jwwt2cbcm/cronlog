package collector_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/example/cronlog/internal/collector"
	"github.com/example/cronlog/internal/retention"
	"github.com/example/cronlog/internal/storage"
)

func tempStore(t *testing.T) *storage.Store {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.json")
	s, err := storage.New(path, retention.DefaultPolicy())
	if err != nil {
		t.Fatalf("storage.New: %v", err)
	}
	t.Cleanup(func() { os.Remove(path) })
	return s
}

func TestRunSuccess(t *testing.T) {
	s := tempStore(t)
	c := collector.New(s)

	res := c.Run("echo-job", "echo", "hello")
	if res.Err != nil {
		t.Fatalf("unexpected error: %v", res.Err)
	}
	if res.Entry.IsError() {
		t.Errorf("expected success, got exit code %d", res.Entry.ExitCode)
	}
	if res.Entry.Job != "echo-job" {
		t.Errorf("job name = %q, want %q", res.Entry.Job, "echo-job")
	}
}

func TestRunFailure(t *testing.T) {
	s := tempStore(t)
	c := collector.New(s)

	res := c.Run("fail-job", "false")
	if res.Err != nil {
		t.Fatalf("unexpected store error: %v", res.Err)
	}
	if !res.Entry.IsError() {
		t.Error("expected error entry for failing command")
	}
}

func TestRunStoredInStore(t *testing.T) {
	s := tempStore(t)
	c := collector.New(s)

	c.Run("job-a", "echo", "first")
	c.Run("job-b", "echo", "second")

	entries, err := s.All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("store has %d entries, want 2", len(entries))
	}
}
