package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func tempDB(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "test.json")
}

func TestBuildRootHasSubcommands(t *testing.T) {
	root := buildRoot()
	names := make(map[string]bool)
	for _, sub := range root.Commands() {
		names[sub.Name()] = true
	}
	for _, expected := range []string{"run", "logs", "prune"} {
		if !names[expected] {
			t.Errorf("expected subcommand %q to be registered", expected)
		}
	}
}

func TestRunCmdStoresEntry(t *testing.T) {
	db := tempDB(t)
	root := buildRoot()

	buf := &bytes.Buffer{}
	root.SetOut(buf)

	root.SetArgs([]string{"--db", db, "--format", "json", "run", "--job", "test-job", "echo", "hello"})
	if err := root.Execute(); err != nil {
		t.Fatalf("run command failed: %v", err)
	}
}

func TestLogsCmdOutputsJSON(t *testing.T) {
	db := tempDB(t)

	// first run a command to populate the store
	run := buildRoot()
	run.SetArgs([]string{"--db", db, "run", "--job", "myjob", "echo", "world"})
	run.SetOut(os.Discard)
	if err := run.Execute(); err != nil {
		t.Fatalf("run failed: %v", err)
	}

	// now query logs in JSON format
	buf := &bytes.Buffer{}
	logs := buildRoot()
	logs.SetOut(buf)
	logs.SetArgs([]string{"--db", db, "--format", "json", "logs", "--job", "myjob"})
	if err := logs.Execute(); err != nil {
		t.Fatalf("logs command failed: %v", err)
	}

	output := strings.TrimSpace(buf.String())
	if output == "" {
		t.Fatal("expected JSON output, got empty string")
	}

	var entries []map[string]interface{}
	if err := json.Unmarshal([]byte(output), &entries); err != nil {
		t.Fatalf("invalid JSON output: %v — output was: %s", err, output)
	}
	if len(entries) == 0 {
		t.Error("expected at least one log entry")
	}
}

func TestPruneCmdRuns(t *testing.T) {
	db := tempDB(t)

	// populate store first
	run := buildRoot()
	run.SetArgs([]string{"--db", db, "run", "echo", "hi"})
	run.SetOut(os.Discard)
	_ = run.Execute()

	buf := &bytes.Buffer{}
	prune := buildRoot()
	prune.SetOut(buf)
	prune.SetArgs([]string{"--db", db, "prune", "--max-age", "30", "--max-entries", "500"})
	if err := prune.Execute(); err != nil {
		t.Fatalf("prune command failed: %v", err)
	}
	if !strings.Contains(buf.String(), "pruned") {
		t.Errorf("expected 'pruned' in output, got: %s", buf.String())
	}
}
