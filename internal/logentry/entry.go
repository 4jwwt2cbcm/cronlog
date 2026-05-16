// Package logentry defines the core log entry type used throughout cronlog.
package logentry

import (
	"strings"
	"time"
)

// Entry represents a single recorded execution of a cron job.
type Entry struct {
	ID        string        `json:"id"`
	Job       string        `json:"job"`
	Output    string        `json:"output"`
	ExitCode  int           `json:"exit_code"`
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time     `json:"timestamp"`
}

// New constructs an Entry with the current timestamp and a generated ID.
func New(job, output string, exitCode int, duration time.Duration) Entry {
	return Entry{
		ID:        generateID(),
		Job:       job,
		Output:    output,
		ExitCode:  exitCode,
		Duration:  duration,
		Timestamp: time.Now().UTC(),
	}
}

// IsError reports whether the entry represents a failed job execution.
func (e Entry) IsError() bool {
	return e.ExitCode != 0
}

// FilterMatches reports whether the entry satisfies the given filter string.
// An empty filter matches all entries.
func (e Entry) FilterMatches(filter string) bool {
	if filter == "" {
		return true
	}
	return strings.Contains(e.Job, filter) ||
		strings.Contains(e.Output, filter)
}

// generateID returns a simple unique identifier based on current time.
func generateID() string {
	return time.Now().UTC().Format("20060102150405.000000000")
}
