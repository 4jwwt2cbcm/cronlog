package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/cronlog/internal/logentry"
)

func makeEntry(job string, exitCode int) logentry.Entry {
	now := time.Now()
	return logentry.Entry{
		ID:         "test-id",
		JobName:    job,
		ExitCode:   exitCode,
		Output:     "some output",
		StartedAt:  now,
		FinishedAt: now.Add(2 * time.Second),
	}
}

func TestPrintText(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf, "text")

	entries := []logentry.Entry{
		makeEntry("backup", 0),
		makeEntry("cleanup", 1),
	}

	if err := p.Print(entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "backup") {
		t.Error("expected 'backup' in text output")
	}
	if !strings.Contains(out, "error") {
		t.Error("expected 'error' status in text output")
	}
	if !strings.Contains(out, "ok") {
		t.Error("expected 'ok' status in text output")
	}
}

func TestPrintJSON(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf, "json")

	entries := []logentry.Entry{makeEntry("sync", 0)}

	if err := p.Print(entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var decoded []logentry.Entry
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(decoded) != 1 {
		t.Errorf("expected 1 entry, got %d", len(decoded))
	}
	if decoded[0].JobName != "sync" {
		t.Errorf("expected job=sync, got %q", decoded[0].JobName)
	}
}

func TestPrintTextEmpty(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf, "text")

	if err := p.Print([]logentry.Entry{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "JOB") {
		t.Error("expected header in empty text output")
	}
}
