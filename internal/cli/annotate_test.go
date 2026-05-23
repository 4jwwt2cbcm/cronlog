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

func tempAnnotateDB(t *testing.T) *storage.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := storage.New(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("storage.New: %v", err)
	}
	return st
}

func makeAnnotateEntry(id, job string, tags map[string]string) logentry.Entry {
	return logentry.Entry{
		ID:        id,
		Job:       job,
		StartedAt: time.Now(),
		Duration:  time.Second,
		ExitCode:  0,
		Output:    "ok",
		Tags:      tags,
	}
}

func TestAnnotateCmdMissingID(t *testing.T) {
	st := tempAnnotateDB(t)
	cmd := buildAnnotateCmd(st)
	cmd.SetArgs([]string{"--tags", "env=prod"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for missing id argument")
	}
}

func TestAnnotateCmdUnknownID(t *testing.T) {
	st := tempAnnotateDB(t)
	cmd := buildAnnotateCmd(st)
	cmd.SetArgs([]string{"nonexistent", "--tags", "env=prod"})
	var buf bytes.Buffer
	cmd.SetErr(&buf)
	if err := cmd.Execute(); err == nil {
		t.Fatal("expected error for unknown id")
	}
}

func TestAnnotateCmdAddsTag(t *testing.T) {
	st := tempAnnotateDB(t)
	e := makeAnnotateEntry("id1", "backup", map[string]string{"env": "dev"})
	if err := st.Add(e); err != nil {
		t.Fatalf("Add: %v", err)
	}

	cmd := buildAnnotateCmd(st)
	cmd.SetArgs([]string{"id1", "--tags", "env=prod,team=ops"})
	var out bytes.Buffer
	cmd.SetOut(&out)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}

	all, _ := st.All()
	if len(all) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(all))
	}
	if all[0].Tags["env"] != "prod" {
		t.Errorf("expected env=prod, got %s", all[0].Tags["env"])
	}
	if all[0].Tags["team"] != "ops" {
		t.Errorf("expected team=ops, got %s", all[0].Tags["team"])
	}
}

func TestAnnotateCmdMergesExistingTags(t *testing.T) {
	st := tempAnnotateDB(t)
	e := makeAnnotateEntry("id2", "sync", map[string]string{"env": "dev", "region": "us-east"})
	if err := st.Add(e); err != nil {
		t.Fatalf("Add: %v", err)
	}

	cmd := buildAnnotateCmd(st)
	cmd.SetArgs([]string{"id2", "--tags", "env=prod"})
	cmd.SetOut(os.Stdout)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}

	all, _ := st.All()
	if all[0].Tags["region"] != "us-east" {
		t.Errorf("existing tag region should be preserved, got %q", all[0].Tags["region"])
	}
	if all[0].Tags["env"] != "prod" {
		t.Errorf("env tag should be overwritten to prod, got %q", all[0].Tags["env"])
	}
}

func TestParseFlagTags(t *testing.T) {
	cases := []struct {
		input    string
		wantLen  int
		wantKey  string
		wantVal  string
	}{
		{"env=prod", 1, "env", "prod"},
		{"env=prod,team=ops", 2, "team", "ops"},
		{"", 0, "", ""},
	}
	for _, tc := range cases {
		tags := parseFlagTags(tc.input)
		if len(tags) != tc.wantLen {
			t.Errorf("input %q: want %d tags, got %d", tc.input, tc.wantLen, len(tags))
		}
		if tc.wantKey != "" && tags[tc.wantKey] != tc.wantVal {
			t.Errorf("input %q: want %s=%s, got %s", tc.input, tc.wantKey, tc.wantVal, tags[tc.wantKey])
		}
	}
}
