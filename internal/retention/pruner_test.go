package retention

import (
	"testing"
	"time"

	"github.com/cronlog/cronlog/internal/logentry"
)

func makeEntry(t time.Time, level string, msg string) logentry.Entry {
	return logentry.Entry{
		Timestamp: t,
		Level:     level,
		Message:   msg,
	}
}

func TestPruneByAge(t *testing.T) {
	now := time.Now()
	policy := Policy{MaxAge: 24 * time.Hour, KeepErrors: false}
	pruner := NewPruner(policy)

	entries := []logentry.Entry{
		makeEntry(now.Add(-48*time.Hour), "info", "old entry"),
		makeEntry(now.Add(-1*time.Hour), "info", "recent entry"),
	}

	result := pruner.Prune(entries)
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0].Message != "recent entry" {
		t.Errorf("unexpected message: %s", result[0].Message)
	}
}

func TestPruneKeepErrors(t *testing.T) {
	now := time.Now()
	policy := Policy{MaxAge: 24 * time.Hour, KeepErrors: true}
	pruner := NewPruner(policy)

	entries := []logentry.Entry{
		makeEntry(now.Add(-48*time.Hour), "error", "old error"),
		makeEntry(now.Add(-48*time.Hour), "info", "old info"),
	}

	result := pruner.Prune(entries)
	if len(result) != 1 {
		t.Fatalf("expected 1 entry (error kept), got %d", len(result))
	}
	if result[0].Level != "error" {
		t.Errorf("expected error level, got %s", result[0].Level)
	}
}

func TestPruneMaxEntries(t *testing.T) {
	now := time.Now()
	policy := Policy{MaxEntries: 2, KeepErrors: false}
	pruner := NewPruner(policy)

	entries := []logentry.Entry{
		makeEntry(now.Add(-3*time.Hour), "info", "oldest"),
		makeEntry(now.Add(-2*time.Hour), "info", "middle"),
		makeEntry(now.Add(-1*time.Hour), "info", "newest"),
	}

	result := pruner.Prune(entries)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	// Newest entries should survive.
	if result[0].Message != "newest" {
		t.Errorf("expected newest first, got %s", result[0].Message)
	}
}

func TestValidatePolicy(t *testing.T) {
	if err := DefaultPolicy().Validate(); err != nil {
		t.Errorf("default policy should be valid: %v", err)
	}
	bad := Policy{MaxAge: -1}
	if err := bad.Validate(); err == nil {
		t.Error("expected error for negative MaxAge")
	}
}
