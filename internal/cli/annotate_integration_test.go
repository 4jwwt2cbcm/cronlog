package cli

import (
	"bytes"
	"path/filepath"
	"testing"

	"cronlog/internal/storage"
)

// TestAnnotateRoundTrip verifies that annotate persists tags across a fresh
// store load, exercising the full storage round-trip.
func TestAnnotateRoundTrip(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "annotate.db")

	// First store instance: add entry and annotate it.
	st1, err := storage.New(dbPath)
	if err != nil {
		t.Fatalf("storage.New: %v", err)
	}
	e := makeAnnotateEntry("rt1", "nightly", nil)
	if err := st1.Add(e); err != nil {
		t.Fatalf("Add: %v", err)
	}

	cmd := buildAnnotateCmd(st1)
	cmd.SetArgs([]string{"rt1", "--tags", "env=prod,owner=alice"})
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("annotate Execute: %v", err)
	}

	// Second store instance: reload from disk and verify tags survived.
	st2, err := storage.New(dbPath)
	if err != nil {
		t.Fatalf("storage.New (reload): %v", err)
	}
	all, err := st2.All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("expected 1 entry after reload, got %d", len(all))
	}
	if all[0].Tags["env"] != "prod" {
		t.Errorf("env tag not persisted: got %q", all[0].Tags["env"])
	}
	if all[0].Tags["owner"] != "alice" {
		t.Errorf("owner tag not persisted: got %q", all[0].Tags["owner"])
	}
}

// TestAnnotateViaRootCmd ensures buildAnnotateCmd is wired into the root
// command and reachable through the CLI entry point.
func TestAnnotateViaRootCmd(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "root.db")

	st, err := storage.New(dbPath)
	if err != nil {
		t.Fatalf("storage.New: %v", err)
	}
	e := makeAnnotateEntry("cli1", "weekly", map[string]string{"tier": "low"})
	if err := st.Add(e); err != nil {
		t.Fatalf("Add: %v", err)
	}

	root := buildRoot(dbPath)
	root.SetArgs([]string{"annotate", "cli1", "--tags", "tier=high"})
	var out bytes.Buffer
	root.SetOut(&out)
	if err := root.Execute(); err != nil {
		t.Fatalf("root Execute: %v", err)
	}

	all, _ := st.All()
	if all[0].Tags["tier"] != "high" {
		t.Errorf("expected tier=high via root cmd, got %q", all[0].Tags["tier"])
	}
}
