package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"cronlog/internal/cli"
	"cronlog/internal/logentry"
	"cronlog/internal/storage"
)

func setupNotifyStore(t *testing.T) (*storage.Store, string) {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "notify_int.json")
	s, err := storage.New(p)
	if err != nil {
		t.Fatalf("storage.New: %v", err)
	}
	return s, dir
}

func addNotifyEntry(t *testing.T, s *storage.Store, job string, code int) {
	t.Helper()
	e := logentry.New(job, "cmd", "out", code, 50)
	e.Timestamp = time.Now()
	if err := s.Add(e); err != nil {
		t.Fatalf("store.Add: %v", err)
	}
}

func TestNotifyHookInvoked(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell hook test skipped on windows")
	}

	store, dir := setupNotifyStore(t)
	outFile := filepath.Join(dir, "hook_out.txt")

	for i := 0; i < 3; i++ {
		addNotifyEntry(t, store, "nightly", 1)
	}

	// Write a tiny shell script that records its argument.
	hookPath := filepath.Join(dir, "hook.sh")
	script := "#!/bin/sh\necho \"$1\" > " + outFile + "\n"
	if err := os.WriteFile(hookPath, []byte(script), 0o755); err != nil {
		t.Fatalf("write hook: %v", err)
	}

	var buf bytes.Buffer
	cfg := cli.NotifyConfig{
		Hook:      hookPath,
		MinFails:  2,
		ErrorRate: 0.5,
	}
	if err := cli.RunNotify(&buf, store, cfg); err != nil {
		t.Fatalf("RunNotify: %v", err)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("hook output file not created: %v", err)
	}
	if len(data) == 0 {
		t.Error("hook was invoked but wrote no output")
	}
}

func TestNotifyDryRunDoesNotInvokeHook(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell hook test skipped on windows")
	}

	store, dir := setupNotifyStore(t)
	sentinel := filepath.Join(dir, "should_not_exist.txt")
	hookPath := filepath.Join(dir, "hook2.sh")
	script := "#!/bin/sh\ntouch " + sentinel + "\n"
	if err := os.WriteFile(hookPath, []byte(script), 0o755); err != nil {
		t.Fatalf("write hook: %v", err)
	}

	addNotifyEntry(t, store, "daily", 1)
	addNotifyEntry(t, store, "daily", 1)

	var buf bytes.Buffer
	cfg := cli.NotifyConfig{
		Hook:      hookPath,
		MinFails:  1,
		ErrorRate: 0.1,
		DryRun:    true,
	}
	if err := cli.RunNotify(&buf, store, cfg); err != nil {
		t.Fatalf("RunNotify: %v", err)
	}

	if _, err := os.Stat(sentinel); err == nil {
		t.Error("hook was invoked despite dry-run flag")
	}
}
