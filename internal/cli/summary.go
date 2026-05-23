package cli

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/your/cronlog/internal/logentry"
)

// JobSummary holds aggregated metrics for a single job.
type JobSummary struct {
	JobName     string
	TotalRuns   int
	ErrorCount  int
	LastRun     time.Time
	AvgDuration time.Duration
	SuccessRate float64
}

// ComputeSummary aggregates entries into per-job summaries.
func ComputeSummary(entries []logentry.Entry) []JobSummary {
	type accumulator struct {
		total    int
		errors   int
		lastRun  time.Time
		durTotal time.Duration
	}

	acc := make(map[string]*accumulator)

	for _, e := range entries {
		a, ok := acc[e.JobName]
		if !ok {
			a = &accumulator{}
			acc[e.JobName] = a
		}
		a.total++
		if e.IsError() {
			a.errors++
		}
		if e.RunAt.After(a.lastRun) {
			a.lastRun = e.RunAt
		}
		a.durTotal += e.Duration
	}

	summaries := make([]JobSummary, 0, len(acc))
	for job, a := range acc {
		avg := time.Duration(0)
		if a.total > 0 {
			avg = a.durTotal / time.Duration(a.total)
		}
		successRate := 0.0
		if a.total > 0 {
			successRate = float64(a.total-a.errors) / float64(a.total) * 100.0
		}
		summaries = append(summaries, JobSummary{
			JobName:     job,
			TotalRuns:   a.total,
			ErrorCount:  a.errors,
			LastRun:     a.lastRun,
			AvgDuration: avg,
			SuccessRate: successRate,
		})
	}

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].JobName < summaries[j].JobName
	})
	return summaries
}

// PrintSummary writes a human-readable summary table to w.
func PrintSummary(w io.Writer, summaries []JobSummary) {
	if len(summaries) == 0 {
		fmt.Fprintln(w, "No entries found.")
		return
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "JOB\tRUNS\tERRORS\tSUCCESS%\tAVG DURATION\tLAST RUN")
	for _, s := range summaries {
		fmt.Fprintf(tw, "%s\t%d\t%d\t%.1f%%\t%s\t%s\n",
			s.JobName,
			s.TotalRuns,
			s.ErrorCount,
			s.SuccessRate,
			s.AvgDuration.Round(time.Millisecond),
			s.LastRun.Format(time.RFC3339),
		)
	}
	tw.Flush()
}
