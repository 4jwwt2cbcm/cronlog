package logentry

import (
	"time"
)

// Level represents the severity level of a log entry.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
	LevelDebug Level = "DEBUG"
)

// Entry represents a single structured log entry from a cron job.
type Entry struct {
	ID        string            `json:"id"`
	JobName   string            `json:"job_name"`
	Timestamp time.Time         `json:"timestamp"`
	Level     Level             `json:"level"`
	Message   string            `json:"message"`
	ExitCode  int               `json:"exit_code"`
	Duration  time.Duration     `json:"duration_ms"`
	Meta      map[string]string `json:"meta,omitempty"`
}

// IsError returns true if the entry represents a failed job execution.
func (e *Entry) IsError() bool {
	return e.Level == LevelError || e.ExitCode != 0
}

// Filter holds criteria for querying log entries.
type Filter struct {
	JobName  string
	Level    Level
	Since    time.Time
	Until    time.Time
	ExitCode *int
}

// Matches returns true if the entry satisfies all non-zero filter criteria.
func (f *Filter) Matches(e *Entry) bool {
	if f.JobName != "" && e.JobName != f.JobName {
		return false
	}
	if f.Level != "" && e.Level != f.Level {
		return false
	}
	if !f.Since.IsZero() && e.Timestamp.Before(f.Since) {
		return false
	}
	if !f.Until.IsZero() && e.Timestamp.After(f.Until) {
		return false
	}
	if f.ExitCode != nil && e.ExitCode != *f.ExitCode {
		return false
	}
	return true
}
