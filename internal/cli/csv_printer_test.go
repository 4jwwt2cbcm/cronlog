package cli

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/example/cronlog/internal/logentry"
)

func makeCSVEntry(job string, exitCode int, output string) logentry.Entry {
	e, _ := logentry.New(job, output, exitCode, 100, time.Now())
	return e
}

func TestCSVPrinterHeader(t *testing.T) {
	var buf bytes.Buffer
	_, err := NewCSVPrinter(&buf)
	if err != nil {
		t.Fatalf("NewCSVPrinter error: %v", err)
	}
	line := buf.String()
	if !strings.Contains(line, "id") || !strings.Contains(line, "job") {
		t.Errorf("header missing expected fields, got: %q", line)
	}
}

func TestCSVPrinterSingleEntry(t *testing.T) {
	var buf bytes.Buffer
	cp, err := NewCSVPrinter(&buf)
	if err != nil {
		t.Fatalf("NewCSVPrinter error: %v", err)
	}

	e := makeCSVEntry("backup", 0, "done")
	if err := cp.Print(e); err != nil {
		t.Fatalf("Print error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "backup") {
		t.Errorf("output missing job name, got: %q", out)
	}
	if !strings.Contains(out, "false") {
		t.Errorf("expected is_error=false in output, got: %q", out)
	}
}

func TestCSVPrinterErrorEntry(t *testing.T) {
	var buf bytes.Buffer
	cp, err := NewCSVPrinter(&buf)
	if err != nil {
		t.Fatalf("NewCSVPrinter error: %v", err)
	}

	e := makeCSVEntry("cleanup", 1, "error occurred")
	if err := cp.Print(e); err != nil {
		t.Fatalf("Print error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "true") {
		t.Errorf("expected is_error=true in output, got: %q", out)
	}
}

func TestCSVPrinterPrintAll(t *testing.T) {
	var buf bytes.Buffer
	cp, err := NewCSVPrinter(&buf)
	if err != nil {
		t.Fatalf("NewCSVPrinter error: %v", err)
	}

	entries := []logentry.Entry{
		makeCSVEntry("job1", 0, "ok"),
		makeCSVEntry("job2", 2, "fail"),
		makeCSVEntry("job3", 0, "done"),
	}
	if err := cp.PrintAll(entries); err != nil {
		t.Fatalf("PrintAll error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// header + 3 data rows
	if len(lines) != 4 {
		t.Errorf("expected 4 lines (header+3), got %d", len(lines))
	}
}
