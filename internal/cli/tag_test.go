package cli

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"cronlog/internal/logentry"
)

func makeTagEntry(job, tags string, exitCode int) logentry.Entry {
	return logentry.Entry{
		ID:        "tid-" + job,
		Job:       job,
		Tags:      tags,
		ExitCode:  exitCode,
		StartedAt: time.Now(),
	}
}

func TestComputeTagSummaryEmpty(t *testing.T) {
	summaries := ComputeTagSummary(nil)
	if len(summaries) != 0 {
		t.Fatalf("expected 0 summaries, got %d", len(summaries))
	}
}

func TestComputeTagSummaryBasic(t *testing.T) {
	entries := []logentry.Entry{
		makeTagEntry("backup", "infra,daily", 0),
		makeTagEntry("cleanup", "infra", 1),
		makeTagEntry("report", "daily", 0),
	}

	summaries := ComputeTagSummary(entries)

	if len(summaries) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(summaries))
	}

	byTag := map[string]TagSummary{}
	for _, s := range summaries {
		byTag[s.Tag] = s
	}

	infra := byTag["infra"]
	if infra.Total != 2 {
		t.Errorf("infra total: want 2, got %d", infra.Total)
	}
	if infra.Errors != 1 {
		t.Errorf("infra errors: want 1, got %d", infra.Errors)
	}
	if len(infra.Jobs) != 2 {
		t.Errorf("infra jobs: want 2, got %d", len(infra.Jobs))
	}

	daily := byTag["daily"]
	if daily.Total != 2 {
		t.Errorf("daily total: want 2, got %d", daily.Total)
	}
	if daily.Errors != 0 {
		t.Errorf("daily errors: want 0, got %d", daily.Errors)
	}
}

func TestComputeTagSummaryNoTags(t *testing.T) {
	entries := []logentry.Entry{
		makeTagEntry("backup", "", 0),
		makeTagEntry("cleanup", "", 0),
	}
	summaries := ComputeTagSummary(entries)
	if len(summaries) != 0 {
		t.Fatalf("expected 0 summaries for untagged entries, got %d", len(summaries))
	}
}

func TestComputeTagSummarySortedByTag(t *testing.T) {
	entries := []logentry.Entry{
		makeTagEntry("job1", "zebra", 0),
		makeTagEntry("job2", "alpha", 0),
		makeTagEntry("job3", "mango", 0),
	}
	summaries := ComputeTagSummary(entries)
	if summaries[0].Tag != "alpha" || summaries[1].Tag != "mango" || summaries[2].Tag != "zebra" {
		t.Errorf("expected alphabetical order, got %v", summaries)
	}
}

func TestBuildTagCmdOutput(t *testing.T) {
	db := tempDB(t)
	cmd := buildTagCmd(db)
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	cmd.SetArgs([]string{})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "no tagged entries found") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}
