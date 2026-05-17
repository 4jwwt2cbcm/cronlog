package logentry_test

import (
	"testing"
	"time"

	"github.com/example/cronlog/internal/logentry"
)

func newEntry(job string, level logentry.Level, exit int, ts time.Time) *logentry.Entry {
	return &logentry.Entry{
		ID:        "test-id",
		JobName:   job,
		Timestamp: ts,
		Level:     level,
		Message:   "test message",
		ExitCode:  exit,
	}
}

func TestIsError(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		entry    *logentry.Entry
		wantErr  bool
	}{
		{"error level", newEntry("job1", logentry.LevelError, 0, now), true},
		{"non-zero exit", newEntry("job1", logentry.LevelInfo, 1, now), true},
		{"info ok", newEntry("job1", logentry.LevelInfo, 0, now), false},
		{"warn ok", newEntry("job1", logentry.LevelWarn, 0, now), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.entry.IsError(); got != tt.wantErr {
				t.Errorf("IsError() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestFilterMatches(t *testing.T) {
	now := time.Now()
	exit1 := 1

	entry := newEntry("backup", logentry.LevelError, 1, now)

	tests := []struct {
		name    string
		filter  logentry.Filter
		want    bool
	}{
		{"empty filter matches all", logentry.Filter{}, true},
		{"matching job", logentry.Filter{JobName: "backup"}, true},
		{"non-matching job", logentry.Filter{JobName: "cleanup"}, false},
		{"matching level", logentry.Filter{Level: logentry.LevelError}, true},
		{"non-matching level", logentry.Filter{Level: logentry.LevelInfo}, false},
		{"since before entry", logentry.Filter{Since: now.Add(-time.Hour)}, true},
		{"since after entry", logentry.Filter{Since: now.Add(time.Hour)}, false},
		{"until after entry", logentry.Filter{Until: now.Add(time.Hour)}, true},
		{"until before entry", logentry.Filter{Until: now.Add(-time.Hour)}, false},
		{"matching exit code", logentry.Filter{ExitCode: &exit1}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.filter.Matches(entry); got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterMatchesNonMatchingExitCode(t *testing.T) {
	now := time.Now()
	exit0 := 0
	exit2 := 2

	entry := newEntry("backup", logentry.LevelError, 1, now)

	tests := []struct {
		name   string
		filter logentry.Filter
		want   bool
	}{
		{"non-matching exit code zero", logentry.Filter{ExitCode: &exit0}, false},
		{"non-matching exit code two", logentry.Filter{ExitCode: &exit2}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.filter.Matches(entry); got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}
