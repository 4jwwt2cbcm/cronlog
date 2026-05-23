package cli

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/your/cronlog/internal/logentry"
)

func makeSummaryEntry(job string, exitCode int, dur time.Duration, runAt time.Time) logentry.Entry {
	return logentry.Entry{
		ID:       "test-id",
		JobName:  job,
		ExitCode: exitCode,
		Duration: dur,
		RunAt:    runAt,
		Output:   "",
	}
}

func TestComputeSummaryEmpty(t *testing.T) {
	result := ComputeSummary(nil)
	if len(result) != 0 {
		t.Fatalf("expected empty summary, got %d entries", len(result))
	}
}

func TestComputeSummaryBasic(t *testing.T) {
	now := time.Now()
	entries := []logentry.Entry{
		makeSummaryEntry("backup", 0, 2*time.Second, now.Add(-10*time.Minute)),
		makeSummaryEntry("backup", 1, 3*time.Second, now.Add(-5*time.Minute)),
		makeSummaryEntry("backup", 0, 1*time.Second, now),
	}

	summaries := ComputeSummary(entries)
	if len(summaries) != 1 {
		t.Fatalf("expected 1 job summary, got %d", len(summaries))
	}

	s := summaries[0]
	if s.JobName != "backup" {
		t.Errorf("expected job 'backup', got %q", s.JobName)
	}
	if s.TotalRuns != 3 {
		t.Errorf("expected 3 runs, got %d", s.TotalRuns)
	}
	if s.ErrorCount != 1 {
		t.Errorf("expected 1 error, got %d", s.ErrorCount)
	}
	expectedAvg := 2 * time.Second // (2+3+1)/3
	if s.AvgDuration != expectedAvg {
		t.Errorf("expected avg %v, got %v", expectedAvg, s.AvgDuration)
	}
	if !s.LastRun.Equal(now) {
		t.Errorf("expected last run %v, got %v", now, s.LastRun)
	}
}

func TestComputeSummarySuccessRate(t *testing.T) {
	now := time.Now()
	entries := []logentry.Entry{
		makeSummaryEntry("sync", 0, time.Second, now),
		makeSummaryEntry("sync", 0, time.Second, now),
		makeSummaryEntry("sync", 0, time.Second, now),
		makeSummaryEntry("sync", 1, time.Second, now),
	}
	summaries := ComputeSummary(entries)
	if len(summaries) != 1 {
		t.Fatalf("expected 1 summary")
	}
	if summaries[0].SuccessRate != 75.0 {
		t.Errorf("expected 75%% success rate, got %.1f", summaries[0].SuccessRate)
	}
}

func TestComputeSummarySortedByJob(t *testing.T) {
	now := time.Now()
	entries := []logentry.Entry{
		makeSummaryEntry("zebra", 0, time.Second, now),
		makeSummaryEntry("alpha", 0, time.Second, now),
		makeSummaryEntry("mango", 0, time.Second, now),
	}
	summaries := ComputeSummary(entries)
	names := []string{summaries[0].JobName, summaries[1].JobName, summaries[2].JobName}
	if names[0] != "alpha" || names[1] != "mango" || names[2] != "zebra" {
		t.Errorf("expected sorted order, got %v", names)
	}
}

func TestPrintSummaryEmpty(t *testing.T) {
	var buf bytes.Buffer
	PrintSummary(&buf, nil)
	if !strings.Contains(buf.String(), "No entries") {
		t.Errorf("expected empty message, got %q", buf.String())
	}
}

func TestPrintSummaryHeader(t *testing.T) {
	now := time.Now()
	summaries := ComputeSummary([]logentry.Entry{
		makeSummaryEntry("cleanup", 0, time.Second, now),
	})
	var buf bytes.Buffer
	PrintSummary(&buf, summaries)
	out := buf.String()
	for _, col := range []string{"JOB", "RUNS", "ERRORS", "SUCCESS%", "AVG DURATION", "LAST RUN"} {
		if !strings.Contains(out, col) {
			t.Errorf("expected column %q in output", col)
		}
	}
	if !strings.Contains(out, "cleanup") {
		t.Errorf("expected job name in output")
	}
}
