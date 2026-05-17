package cli

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/cronlog/internal/logentry"
)

func makeStatsEntry(job string, exitCode int, dur time.Duration, runAt time.Time) logentry.Entry {
	e := logentry.Entry{
		Job:      job,
		ExitCode: exitCode,
		Duration: dur,
		RunAt:    runAt,
	}
	return e
}

func TestComputeStatsEmpty(t *testing.T) {
	stats := ComputeStats(nil)
	if len(stats) != 0 {
		t.Fatalf("expected 0 stats, got %d", len(stats))
	}
}

func TestComputeStatsCounts(t *testing.T) {
	now := time.Now()
	entries := []logentry.Entry{
		makeStatsEntry("backup", 0, 2*time.Second, now.Add(-2*time.Hour)),
		makeStatsEntry("backup", 1, 3*time.Second, now.Add(-1*time.Hour)),
		makeStatsEntry("backup", 0, 1*time.Second, now),
		makeStatsEntry("cleanup", 0, 500*time.Millisecond, now.Add(-30*time.Minute)),
	}

	stats := ComputeStats(entries)

	if len(stats) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(stats))
	}

	// sorted alphabetically: backup, cleanup
	backup := stats[0]
	if backup.JobName != "backup" {
		t.Errorf("expected backup, got %s", backup.JobName)
	}
	if backup.TotalRuns != 3 {
		t.Errorf("expected 3 runs, got %d", backup.TotalRuns)
	}
	if backup.ErrorCount != 1 {
		t.Errorf("expected 1 error, got %d", backup.ErrorCount)
	}
	expectedAvg := (2*time.Second + 3*time.Second + 1*time.Second) / 3
	if backup.AvgRuntime != expectedAvg {
		t.Errorf("expected avg %v, got %v", expectedAvg, backup.AvgRuntime)
	}
	if !backup.LastRun.Equal(now) {
		t.Errorf("expected last run to be now")
	}
}

func TestComputeStatsSorting(t *testing.T) {
	now := time.Now()
	entries := []logentry.Entry{
		makeStatsEntry("zebra", 0, time.Second, now),
		makeStatsEntry("alpha", 0, time.Second, now),
		makeStatsEntry("mango", 0, time.Second, now),
	}
	stats := ComputeStats(entries)
	names := []string{stats[0].JobName, stats[1].JobName, stats[2].JobName}
	if names[0] != "alpha" || names[1] != "mango" || names[2] != "zebra" {
		t.Errorf("unexpected sort order: %v", names)
	}
}

func TestPrintStatsEmpty(t *testing.T) {
	var buf bytes.Buffer
	PrintStats(&buf, nil)
	if !strings.Contains(buf.String(), "No log entries") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestPrintStatsOutput(t *testing.T) {
	now := time.Now()
	stats := []JobStats{
		{JobName: "backup", TotalRuns: 5, ErrorCount: 1, LastRun: now, AvgRuntime: 2 * time.Second},
	}
	var buf bytes.Buffer
	PrintStats(&buf, stats)
	out := buf.String()
	if !strings.Contains(out, "backup") {
		t.Errorf("expected job name in output")
	}
	if !strings.Contains(out, "5") {
		t.Errorf("expected run count in output")
	}
	if !strings.Contains(out, "1") {
		t.Errorf("expected error count in output")
	}
}
