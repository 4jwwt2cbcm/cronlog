package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"cronlog/internal/logentry"
)

func makeExportEntry(job string, isErr bool) logentry.Entry {
	e := logentry.Entry{
		ID:        "export-id-1",
		Job:       job,
		Output:    "output line",
		ExitCode:  0,
		StartedAt: time.Now().Add(-5 * time.Second),
		EndedAt:   time.Now(),
	}
	if isErr {
		e.ExitCode = 1
		e.Error = "something failed"
	}
	return e
}

func TestParseExportFormat(t *testing.T) {
	cases := []struct {
		input   string
		want    ExportFormat
		wantErr bool
	}{
		{"json", ExportJSON, false},
		{"JSON", ExportJSON, false},
		{"csv", ExportCSV, false},
		{"text", ExportText, false},
		{"xml", "", true},
	}
	for _, tc := range cases {
		got, err := ParseExportFormat(tc.input)
		if tc.wantErr {
			if err == nil {
				t.Errorf("ParseExportFormat(%q): expected error, got nil", tc.input)
			}
			continue
		}
		if err != nil {
			t.Fatalf("ParseExportFormat(%q): unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParseExportFormat(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestExportEntriesJSON(t *testing.T) {
	entries := []logentry.Entry{
		makeExportEntry("backup", false),
		makeExportEntry("sync", true),
	}
	tmp := filepath.Join(t.TempDir(), "out.json")
	if err := ExportEntries(entries, tmp, ExportJSON); err != nil {
		t.Fatalf("ExportEntries JSON: %v", err)
	}
	data, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatalf("read exported file: %v", err)
	}
	var got []logentry.Entry
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal json: %v", err)
	}
	if len(got) != len(entries) {
		t.Errorf("got %d entries, want %d", len(got), len(entries))
	}
}

func TestExportEntriesCSV(t *testing.T) {
	entries := []logentry.Entry{makeExportEntry("daily", false)}
	tmp := filepath.Join(t.TempDir(), "out.csv")
	if err := ExportEntries(entries, tmp, ExportCSV); err != nil {
		t.Fatalf("ExportEntries CSV: %v", err)
	}
	data, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatalf("read exported file: %v", err)
	}
	if !strings.Contains(string(data), "daily") {
		t.Errorf("CSV output missing job name 'daily'")
	}
}

func TestExportEntriesText(t *testing.T) {
	entries := []logentry.Entry{makeExportEntry("weekly", false)}
	tmp := filepath.Join(t.TempDir(), "out.txt")
	if err := ExportEntries(entries, tmp, ExportText); err != nil {
		t.Fatalf("ExportEntries Text: %v", err)
	}
	data, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatalf("read exported file: %v", err)
	}
	if !strings.Contains(string(data), "weekly") {
		t.Errorf("text output missing job name 'weekly'")
	}
}

func TestExportEntriesInvalidFormat(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "out.xml")
	err := ExportEntries(nil, tmp, ExportFormat("xml"))
	if err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}
