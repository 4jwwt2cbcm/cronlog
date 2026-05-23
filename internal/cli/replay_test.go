package cli

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"cronlog/internal/logentry"
	"cronlog/internal/storage"
)

func tempReplayDB(t *testing.T) *storage.Store {
	t.Helper()
	path := filepath.Join(t.TempDir(), "replay.db")
	s, err := storage.New(path)
	if err != nil {
		t.Fatalf("storage.New: %v", err)
	}
	return s
}

func makeReplayEntry(job, command string, exitCode int) logentry.Entry {
	e := logentry.New(job, command, "output", exitCode, time.Now(), 100)
	return e
}

func TestReplayCmdMissingID(t *testing.T) {
	s := tempReplayDB(t)
	cmd := buildReplayCmd(s)
	cmd.SetArgs([]string{})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing ID argument")
	}
}

func TestReplayCmdUnknownID(t *testing.T) {
	s := tempReplayDB(t)
	cmd := buildReplayCmd(s)
	cmd.SetArgs([]string{"nonexistent-id"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for unknown entry ID")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' in error, got: %v", err)
	}
}

func TestReplayCmdNoCommand(t *testing.T) {
	s := tempReplayDB(t)
	e := makeReplayEntry("noop", "", 0)
	if err := s.Add(e); err != nil {
		t.Fatalf("Add: %v", err)
	}

	cmd := buildReplayCmd(s)
	cmd.SetArgs([]string{e.ID})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error when command is empty")
	}
	if !strings.Contains(err.Error(), "no recorded command") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestReplayCmdSuccess(t *testing.T) {
	s := tempReplayDB(t)
	e := makeReplayEntry("echo-job", "echo hello", 0)
	if err := s.Add(e); err != nil {
		t.Fatalf("Add: %v", err)
	}

	cmd := buildReplayCmd(s)
	cmd.SetArgs([]string{e.ID})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Replay complete") {
		t.Errorf("expected 'Replay complete' in output, got: %s", out)
	}
	if !strings.Contains(out, "exit code: 0") {
		t.Errorf("expected exit code 0 in output, got: %s", out)
	}

	all, err := s.All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("expected 2 entries after replay, got %d", len(all))
	}
}
