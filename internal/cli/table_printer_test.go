package cli

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/example/cronlog/internal/logentry"
)

func makeTableEntry(job string, exit int, ago time.Duration) logentry.Entry {
	now := time.Now().Add(-ago)
	return logentry.Entry{
		ID:        "abcdef1234567890",
		JobName:   job,
		ExitCode:  exit,
		StartedAt: now,
		Duration:  250 * time.Millisecond,
		Output:    "some output",
	}
}

func TestTablePrinterHeader(t *testing.T) {
	var buf bytes.Buffer
	p := NewTablePrinter(&buf)
	p.PrintHeader()
	p.w.Flush()
	out := buf.String()
	for _, col := range []string{"ID", "JOB", "STATUS", "EXIT", "STARTED", "DURATION"} {
		if !strings.Contains(out, col) {
			t.Errorf("expected header to contain %q, got:\n%s", col, out)
		}
	}
}

func TestTablePrinterSingleEntry(t *testing.T) {
	var buf bytes.Buffer
	p := NewTablePrinter(&buf)
	e := makeTableEntry("backup", 0, time.Hour)
	p.PrintAll([]logentry.Entry{e})
	out := buf.String()
	if !strings.Contains(out, "backup") {
		t.Errorf("expected output to contain job name, got:\n%s", out)
	}
	if !strings.Contains(out, "OK") {
		t.Errorf("expected status OK, got:\n%s", out)
	}
}

func TestTablePrinterErrorEntry(t *testing.T) {
	var buf bytes.Buffer
	p := NewTablePrinter(&buf)
	e := makeTableEntry("cleanup", 1, time.Minute)
	p.PrintAll([]logentry.Entry{e})
	out := buf.String()
	if !strings.Contains(out, "ERROR") {
		t.Errorf("expected status ERROR for non-zero exit, got:\n%s", out)
	}
}

func TestTablePrinterEmpty(t *testing.T) {
	var buf bytes.Buffer
	p := NewTablePrinter(&buf)
	p.PrintAll([]logentry.Entry{})
	out := buf.String()
	if !strings.Contains(out, "No log entries found.") {
		t.Errorf("expected empty message, got:\n%s", out)
	}
}

func TestTablePrinterMultipleEntries(t *testing.T) {
	var buf bytes.Buffer
	p := NewTablePrinter(&buf)
	entries := []logentry.Entry{
		makeTableEntry("jobA", 0, 2*time.Hour),
		makeTableEntry("jobB", 2, time.Hour),
		makeTableEntry("jobC", 0, 30*time.Minute),
	}
	p.PrintAll(entries)
	out := buf.String()
	for _, name := range []string{"jobA", "jobB", "jobC"} {
		if !strings.Contains(out, name) {
			t.Errorf("expected output to contain %q, got:\n%s", name, out)
		}
	}
	if strings.Count(out, "ERROR") != 1 {
		t.Errorf("expected exactly 1 ERROR row, got:\n%s", out)
	}
}
