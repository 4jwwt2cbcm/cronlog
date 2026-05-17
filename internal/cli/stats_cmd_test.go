package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cronlog/internal/logentry"
	"github.com/cronlog/internal/storage"
)

func tempStatsDB(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "stats_test.db")
}

func seedStore(t *testing.T, path string, entries []logentry.Entry) {
	t.Helper()
	s, err := storage.New(path)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	for _, e := range entries {
		if err := s.Add(e); err != nil {
			t.Fatalf("add entry: %v", err)
		}
	}
}

func TestStatsCmdNoEntries(t *testing.T) {
	db := tempStatsDB(t)
	cmd := buildStatsCmd(db)

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	// redirect stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := cmd.Execute()
	w.Close()
	os.Stdout = old

	var out bytes.Buffer
	out.ReadFrom(r)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "No log entries") {
		t.Errorf("expected empty message, got: %s", out.String())
	}
}

func TestStatsCmdWithEntries(t *testing.T) {
	db := tempStatsDB(t)
	now := time.Now()
	entries := []logentry.Entry{
		{Job: "backup", ExitCode: 0, Duration: 2 * time.Second, RunAt: now},
		{Job: "backup", ExitCode: 1, Duration: 3 * time.Second, RunAt: now.Add(-time.Hour)},
		{Job: "cleanup", ExitCode: 0, Duration: 500 * time.Millisecond, RunAt: now},
	}
	seedStore(t, db, entries)

	cmd := buildStatsCmd(db)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := cmd.Execute()
	w.Close()
	os.Stdout = old

	var out bytes.Buffer
	out.ReadFrom(r)
	outStr := out.String()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(outStr, "backup") {
		t.Errorf("expected backup in output")
	}
	if !strings.Contains(outStr, "cleanup") {
		t.Errorf("expected cleanup in output")
	}
}

func TestStatsCmdJobFilter(t *testing.T) {
	db := tempStatsDB(t)
	now := time.Now()
	entries := []logentry.Entry{
		{Job: "backup", ExitCode: 0, Duration: time.Second, RunAt: now},
		{Job: "cleanup", ExitCode: 0, Duration: time.Second, RunAt: now},
	}
	seedStore(t, db, entries)

	cmd := buildStatsCmd(db)
	cmd.SetArgs([]string{"--job", "backup"})

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := cmd.Execute()
	w.Close()
	os.Stdout = old

	var out bytes.Buffer
	out.ReadFrom(r)
	outStr := out.String()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(outStr, "backup") {
		t.Errorf("expected backup in output")
	}
	if strings.Contains(outStr, "cleanup") {
		t.Errorf("cleanup should be filtered out")
	}
}
