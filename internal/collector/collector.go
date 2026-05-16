// Package collector provides functionality for capturing and recording
// cron job output into structured log entries.
package collector

import (
	"bytes"
	"os/exec"
	"time"

	"github.com/example/cronlog/internal/logentry"
	"github.com/example/cronlog/internal/storage"
)

// Collector runs cron job commands and stores their output.
type Collector struct {
	store *storage.Store
}

// New creates a new Collector backed by the given store.
func New(store *storage.Store) *Collector {
	return &Collector{store: store}
}

// RunResult holds the outcome of a collected job run.
type RunResult struct {
	Entry    logentry.Entry
	Err      error
}

// Run executes the named command with the provided arguments, captures its
// combined output, and persists a log entry to the store.
func (c *Collector) Run(jobName string, command string, args ...string) RunResult {
	start := time.Now()

	var buf bytes.Buffer
	cmd := exec.Command(command, args...)
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	cmdErr := cmd.Run()
	duration := time.Since(start)

	exitCode := 0
	if cmdErr != nil {
		if exitErr, ok := cmdErr.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}

	entry := logentry.New(jobName, buf.String(), exitCode, duration)

	storeErr := c.store.Add(entry)
	return RunResult{
		Entry: entry,
		Err:   storeErr,
	}
}
