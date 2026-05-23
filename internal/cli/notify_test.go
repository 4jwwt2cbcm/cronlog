package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"cronlog/internal/logentry"
	"cronlog/internal/storage"
)

func tempNotifyDB(t *testing.T) *storage.Store {
	t.Helper()
	p := filepath.Join(t.TempDir(), "notify.json")
	s, err := storage.New(p)
	if err != nil {
		t.Fatalf("storage.New: %v", err)
	}
	return s
}

func makeNotifyEntry(job string, exitCode int, ts time.Time) logentry.Entry {
	e := logentry.New(job, "echo test", "output", exitCode, 100)
	e.Timestamp = ts
	return e
}

func TestRunNotifyNoAlerts(t *testing.T) {
	store := tempNotifyDB(t)
	now := time.Now()
	_ = store.Add(makeNotifyEntry("backup", 0, now))
	_ = store.Add(makeNotifyEntry("backup", 0, now.Add(-time.Minute)))

	var buf bytes.Buffer
	cfg := NotifyConfig{MinFails: 2, ErrorRate: 0.8}
	if err := RunNotify(&buf, store, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := buf.String(); got != "no alerts triggered\n" {
		t.Errorf("expected no-alert message, got %q", got)
	}
}

func TestRunNotifyTriggersAlert(t *testing.T) {
	store := tempNotifyDB(t)
	now := time.Now()
	for i := 0; i < 3; i++ {
		_ = store.Add(makeNotifyEntry("sync", 1, now.Add(time.Duration(-i)*time.Minute)))
	}

	var buf bytes.Buffer
	cfg := NotifyConfig{MinFails: 2, ErrorRate: 0.5}
	if err := RunNotify(&buf, store, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Error("expected alert output, got empty")
	}
}

func TestRunNotifyDryRun(t *testing.T) {
	store := tempNotifyDB(t)
	now := time.Now()
	_ = store.Add(makeNotifyEntry("report", 2, now))
	_ = store.Add(makeNotifyEntry("report", 2, now.Add(-time.Minute)))

	var buf bytes.Buffer
	cfg := NotifyConfig{Hook: "/usr/bin/true", MinFails: 1, ErrorRate: 0.1, DryRun: true}
	if err := RunNotify(&buf, store, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := buf.String(); len(got) == 0 {
		t.Error("expected dry-run output")
	}
	if contains := bytes.Contains(buf.Bytes(), []byte("dry-run")); !contains {
		t.Errorf("expected dry-run marker in output, got: %s", buf.String())
	}
}

func TestRunNotifyJobFilter(t *testing.T) {
	store := tempNotifyDB(t)
	now := time.Now()
	_ = store.Add(makeNotifyEntry("backup", 1, now))
	_ = store.Add(makeNotifyEntry("backup", 1, now.Add(-time.Minute)))
	_ = store.Add(makeNotifyEntry("sync", 0, now))

	var buf bytes.Buffer
	cfg := NotifyConfig{MinFails: 1, ErrorRate: 0.1, JobFilter: "sync"}
	if err := RunNotify(&buf, store, cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := buf.String(); got != "no alerts triggered\n" {
		t.Errorf("sync job should not trigger alert, got: %s", got)
	}
}

func TestBuildNotifyCmdRegistered(t *testing.T) {
	p := filepath.Join(t.TempDir(), "db.json")
	cmd := buildNotifyCmd(p)
	if cmd.Use != "notify" {
		t.Errorf("expected Use=notify, got %q", cmd.Use)
	}
	flags := []string{"hook", "min-failures", "error-rate", "job", "dry-run"}
	for _, f := range flags {
		if cmd.Flags().Lookup(f) == nil {
			t.Errorf("flag --%s not registered", f)
		}
	}
	_ = os.Remove(p)
}
