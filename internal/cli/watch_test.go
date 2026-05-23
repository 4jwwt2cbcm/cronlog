package cli

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"cronlog/internal/logentry"
	"cronlog/internal/storage"
)

func tempWatchDB(t *testing.T) *storage.Store {
	t.Helper()
	dir := t.TempDir()
	store, err := storage.New(filepath.Join(dir, "watch.db"))
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	return store
}

func makeWatchEntry(job string, exitCode int) logentry.Entry {
	e := logentry.New(job, "echo hi", "output", exitCode, 10*time.Millisecond)
	return e
}

func TestFilterNewReturnsUnseen(t *testing.T) {
	entries := []logentry.Entry{
		makeWatchEntry("job1", 0),
		makeWatchEntry("job2", 1),
	}
	seen := map[string]bool{}
	opts := &WatchOptions{}

	out := filterNew(entries, seen, opts)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestFilterNewSkipsSeen(t *testing.T) {
	e := makeWatchEntry("job1", 0)
	seen := map[string]bool{e.ID: true}
	opts := &WatchOptions{}

	out := filterNew([]logentry.Entry{e}, seen, opts)
	if len(out) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(out))
	}
}

func TestFilterNewJobFilter(t *testing.T) {
	entries := []logentry.Entry{
		makeWatchEntry("backup", 0),
		makeWatchEntry("cleanup", 0),
	}
	seen := map[string]bool{}
	opts := &WatchOptions{JobFilter: "backup"}

	out := filterNew(entries, seen, opts)
	if len(out) != 1 || out[0].Job != "backup" {
		t.Fatalf("expected 1 backup entry, got %v", out)
	}
}

func TestFilterNewOnlyErrors(t *testing.T) {
	entries := []logentry.Entry{
		makeWatchEntry("job1", 0),
		makeWatchEntry("job2", 1),
	}
	seen := map[string]bool{}
	opts := &WatchOptions{OnlyErrors: true}

	out := filterNew(entries, seen, opts)
	if len(out) != 1 || out[0].ExitCode != 1 {
		t.Fatalf("expected 1 error entry, got %v", out)
	}
}

func TestBuildWatchCmdExists(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")
	cmd := buildWatchCmd(dbPath)
	if cmd.Use != "watch" {
		t.Errorf("expected use=watch, got %s", cmd.Use)
	}
	if cmd.Flags().Lookup("interval") == nil {
		t.Error("expected --interval flag")
	}
	if cmd.Flags().Lookup("errors") == nil {
		t.Error("expected --errors flag")
	}
	_ = os.RemoveAll(dir)
}
