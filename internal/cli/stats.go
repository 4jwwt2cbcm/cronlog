package cli

import (
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/cronlog/internal/logentry"
)

// JobStats holds aggregated statistics for a single job.
type JobStats struct {
	JobName    string
	TotalRuns  int
	ErrorCount int
	LastRun    time.Time
	AvgRuntime time.Duration
}

// ComputeStats aggregates log entries into per-job statistics.
func ComputeStats(entries []logentry.Entry) []JobStats {
	type accumulator struct {
		total    int
		errors   int
		lastRun  time.Time
		totalDur time.Duration
	}

	acc := make(map[string]*accumulator)

	for _, e := range entries {
		a, ok := acc[e.Job]
		if !ok {
			a = &accumulator{}
			acc[e.Job] = a
		}
		a.total++
		if e.IsError() {
			a.errors++
		}
		if e.RunAt.After(a.lastRun) {
			a.lastRun = e.RunAt
		}
		a.totalDur += e.Duration
	}

	stats := make([]JobStats, 0, len(acc))
	for job, a := range acc {
		avg := time.Duration(0)
		if a.total > 0 {
			avg = a.totalDur / time.Duration(a.total)
		}
		stats = append(stats, JobStats{
			JobName:    job,
			TotalRuns:  a.total,
			ErrorCount: a.errors,
			LastRun:    a.lastRun,
			AvgRuntime: avg,
		})
	}

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].JobName < stats[j].JobName
	})
	return stats
}

// PrintStats writes a human-readable stats summary to w.
func PrintStats(w io.Writer, stats []JobStats) {
	if len(stats) == 0 {
		fmt.Fprintln(w, "No log entries found.")
		return
	}
	fmt.Fprintf(w, "%-30s %8s %8s %12s %20s\n", "JOB", "RUNS", "ERRORS", "AVG RUNTIME", "LAST RUN")
	fmt.Fprintf(w, "%s\n", fmt.Sprintf("%s", repeat('-', 82)))
	for _, s := range stats {
		last := s.LastRun.Format(time.RFC3339)
		if s.LastRun.IsZero() {
			last = "never"
		}
		fmt.Fprintf(w, "%-30s %8d %8d %12s %20s\n",
			s.JobName, s.TotalRuns, s.ErrorCount,
			s.AvgRuntime.Round(time.Millisecond).String(),
			last,
		)
	}
}

func repeat(ch rune, n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = ch
	}
	return string(b)
}
